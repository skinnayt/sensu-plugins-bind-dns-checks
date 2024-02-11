package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	StatisticsFormat   string
	StatisticsFilePath string
	StatisticsIP       string
	StatisticsPort     int
	OutputFormat       string
	returnMetrics      []*Metric
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-plugins-bind-dns-checks",
			Short:    "Sensu check to pull metrics from bind DNS server",
			Keyspace: "sensu.io/plugins/sensu-plugins-bind-dns-checks/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "statistics-format",
			Env:       "STATISTICS_FORMAT",
			Argument:  "statistics-format",
			Shorthand: "f",
			Default:   "file",
			Usage:     "The format of the statistics file (file, xml, json)",
			Value:     &plugin.StatisticsFormat,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "statistics-filepath",
			Env:       "STATISTICS_FILEPATH",
			Argument:  "statistics-filepath",
			Shorthand: "p",
			Default:   "",
			Usage:     "The file path to the statistics file",
			Value:     &plugin.StatisticsFilePath,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "statistics-ip",
			Env:       "STATISTICS_IP",
			Argument:  "statistics-ip",
			Shorthand: "a",
			Default:   "",
			Usage:     "The IP address to listen on for the statistics channel",
			Value:     &plugin.StatisticsIP,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "statistics-port",
			Env:       "STATISTICS_PORT",
			Argument:  "statistics-port",
			Shorthand: "P",
			Default:   0,
			Usage:     "The port to listen on for the statistics channel",
			Value:     &plugin.StatisticsPort,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "output-format",
			Env:       "OUTPUT_FORMAT",
			Argument:  "output-format",
			Shorthand: "o",
			Default:   "",
			Usage:     "The format to output the metrics in (graphite, prometheus)",
			Value:     &plugin.OutputFormat,
		},
	}
)

type MetricTag [2]string

func (mt *MetricTag) String() string {
	return fmt.Sprintf("%s_%s", mt[0], mt[1])
}

type Metric struct {
	Name      string
	Value     int64
	Timestamp time.Time
	Tags      []*MetricTag
}

type namedStats struct {
	statsTags []*MetricTag
	curLevel  string
}

func (m *Metric) Graphite(tag_prefix string) string {
	var tags []string
	if tag_prefix != "" {
		tags = append(tags, tag_prefix)
	}
	for _, tag := range m.Tags {
		tags = append(tags, tag.String())
	}

	return fmt.Sprintf(
		"%s.%s %d %d",
		strings.Join(tags, "."),
		strings.Replace(m.Name, " ", "_", -1),
		m.Value,
		m.Timestamp.Unix(),
	)
}

func (m *Metric) Prometheus() string {
	var tags []string
	for _, tag := range m.Tags {
		tags = append(tags, fmt.Sprintf("%s=\"%s\"", tag[0], tag[1]))
	}

	return fmt.Sprintf(
		"%s{%s} %d %d",
		strings.Replace(m.Name, " ", "_", -1),
		strings.Join(tags, ","),
		m.Value,
		m.Timestamp.UnixMilli(),
	)
}

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	// Check that we got an appropriate statistics format
	if plugin.StatisticsFormat == "file" {
		if plugin.StatisticsFilePath == "" {
			return sensu.CheckStateUnknown, fmt.Errorf("no statistics file path specified when using file format")
		}
		// Check that the file exists
		if _, err := os.Stat(plugin.StatisticsFilePath); os.IsNotExist(err) {
			return sensu.CheckStateUnknown, fmt.Errorf("statistics file does not exist: %s", plugin.StatisticsFilePath)
		}
		// Check that the file is readable
		if _, err := os.Open(plugin.StatisticsFilePath); err != nil {
			return sensu.CheckStateUnknown, fmt.Errorf("unable to read statistics file: %s", err)
		}
	} else if (plugin.StatisticsFormat == "xml") || (plugin.StatisticsFormat == "json") {
		if plugin.StatisticsIP == "" {
			return sensu.CheckStateUnknown, fmt.Errorf("no statistics IP specified when using %s format", plugin.StatisticsFormat)
		}
		// Check that the IP is valid
		// TODO: Maybe allow a hostname here?
		if net.ParseIP(plugin.StatisticsIP) == nil {
			return sensu.CheckStateUnknown, fmt.Errorf("invalid statistics IP specified: %s", plugin.StatisticsIP)
		}

		// Check that the port is valid
		if plugin.StatisticsPort < 1 || plugin.StatisticsPort > 65535 {
			return sensu.CheckStateUnknown, fmt.Errorf("invalid statistics port specified: %d", plugin.StatisticsPort)
		}
	} else {
		return sensu.CheckStateUnknown, fmt.Errorf("invalid statistics format: %s", plugin.StatisticsFormat)
	}

	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	if plugin.StatisticsFormat == "file" {
		if err := readStatisticsFile(); err != nil {
			return sensu.CheckStateCritical, fmt.Errorf("error reading statistics file: %s", err)
		}
	} else if plugin.StatisticsFormat == "xml" || plugin.StatisticsFormat == "json" {
		if err := readStatisticsChannel(); err != nil {
			return sensu.CheckStateCritical, fmt.Errorf("error reading statistics channel: %s", err)
		}
	}

	// Dump out the metrics loaded from the statistics file or channel
	if plugin.OutputFormat == "graphite" {
		OutputMetricsGraphite()
	} else if plugin.OutputFormat == "prometheus" {
		OutputMetricsPrometheus()
	}

	return sensu.CheckStateOK, nil
}

// Read from statistics file
func readStatisticsFile() error {
	dnsStats, err := os.ReadFile(plugin.StatisticsFilePath)
	if err != nil {
		return err
	}

	// Parse the statistics file
	var statsReadTime time.Time

	namedStats := &namedStats{}
	namedStats.statsTags = []*MetricTag{}

	// Regular expressions for parsing the statistics file
	var statsFile = make(map[string]*regexp.Regexp)
	statsFile["start_end"], _ = regexp.Compile(`^(?:[-+]{3}) Statistics Dump (?:[-+]{3}) \((?P<unixtime>[0-9]*)\)$`)
	statsFile["sections"], _ = regexp.Compile(`^(?:[+]{2}) (?P<section>[a-zA-Z0-9_/ ]+) (?:[+]{2})$`)
	statsFile["metric"], _ = regexp.Compile(`^\s*(?P<value>[0-9]+) (?P<name>[-a-zA-Z0-9_/!#()<> ]+)\s*$`)
	statsFile["view"], _ = regexp.Compile(`^\[View: (?P<view>[a-zA-Z0-9_/ ]+)\]$`)
	statsFile["view_cache"], _ = regexp.Compile(`^\[View: (?P<view>[a-zA-Z0-9_/ ]+) \(Cache: (?P<cache>[a-zA-Z0-9_/ ]+)\)\]$`)
	statsFile["subsection"], _ = regexp.Compile(`^\[(?P<subsection>[-a-zA-Z0-9_/!#()<>]+)\]$`)
	statsFile["zone"], _ = regexp.Compile(`^\[(?P<zone>\.|(?:[a-z]+)(?:\.[a-z]+){1,}|(?:[0-9A-F]+\.)*(?:IN-ADDR|IP6|HOME|EMPTY\.AS112)\.ARPA)\]$`)
	statsFile["bind_var"], _ = regexp.Compile(`^\[(?P<bind_var>[a-z.]+) \(view: _bind\)\]$`)

	for _, line := range strings.Split(string(dnsStats), "\n") {
		// Parse the line
		if matched := statsFile["start_end"].FindStringSubmatch(line); matched != nil {
			// Start of a new statistics file or end of the file
			if statsReadTime.Equal(time.Time{}) {
				unixTime, _ := strconv.ParseInt(matched[1], 10, 64)
				statsReadTime = time.Unix(unixTime, 0)
			}
		} else if section := statsFile["sections"].FindStringSubmatch(line); section != nil {
			// Start of a new section
			namedStats.curLevel = section[1]
			namedStats.statsTags = []*MetricTag{}
		} else if metric := statsFile["metric"].FindStringSubmatch(line); metric != nil {
			// Metric
			value, _ := strconv.ParseInt(metric[1], 10, 64)
			plugin.returnMetrics = append(plugin.returnMetrics, &Metric{
				Name:      metric[2],
				Value:     value,
				Timestamp: statsReadTime,
				Tags:      namedStats.statsTags,
			})
		} else if view := statsFile["view"].FindStringSubmatch(line); view != nil {
			namedStats.statsTags = append(namedStats.statsTags, &MetricTag{"view", view[1]})
		} else if viewCache := statsFile["view_cache"].FindStringSubmatch(line); viewCache != nil {
			namedStats.statsTags = append(namedStats.statsTags, &MetricTag{"view", viewCache[1]})
			namedStats.statsTags = append(namedStats.statsTags, &MetricTag{"cache", viewCache[2]})
		} else if subsection := statsFile["subsection"].FindStringSubmatch(line); subsection != nil {
			namedStats.statsTags = append(namedStats.statsTags, &MetricTag{"subsection", subsection[1]})
		} else if zone := statsFile["zone"].FindStringSubmatch(line); zone != nil {
			namedStats.statsTags = append(namedStats.statsTags, &MetricTag{"zone", zone[1]})
		} else if bindVar := statsFile["bind_var"].FindStringSubmatch(line); bindVar != nil {
			namedStats.statsTags = append(namedStats.statsTags, &MetricTag{"bind_var", bindVar[1]})
			namedStats.statsTags = append(namedStats.statsTags, &MetricTag{"view", "_bind"})
		} else if strings.Trim(line, " ") == "" {
			// Skip blank lines
			continue
		} else {
			// Unrecognized line
			// XXX: Handle this later
			continue
		}
	}

	return nil
}

// Read from statistics channel
func readStatisticsChannel() error {
	// Make the URL for connecting to the statistics channel
	tcpAddr := net.TCPAddr{
		IP:   net.ParseIP(plugin.StatisticsIP),
		Port: plugin.StatisticsPort,
	}

	statsUrl := url.URL{
		Scheme: "http",
		Host:   tcpAddr.String(),
		Path:   "/",
	}
	if plugin.StatisticsFormat == "xml" {
		statsUrl.Path = "/xml/v3"
	}
	if plugin.StatisticsFormat == "json" {
		statsUrl.Path = "/json/v1"
	}

	statsReq, _ := http.NewRequest("GET", statsUrl.String(), nil)
	statsReq.Header.Add("Accept", "application/"+plugin.StatisticsFormat)

	// Connect to the statistics channel
	statsClient := &http.Client{}
	statsResp, err := statsClient.Do(statsReq)
	if err != nil {
		return err
	}

	defer statsResp.Body.Close()

	if statsResp.StatusCode != 200 {
		return fmt.Errorf("error reading statistics channel: %s", statsResp.Status)
	}

	// Read the XML statistics
	var statsData []byte
	readData := make([]byte, 1024)
	for {
		n, err := statsResp.Body.Read(readData)

		if err != nil {
			if err == io.EOF {
				if n > 0 {
					statsData = append(statsData, readData[:n]...)
				}
				break
			}
			return err
		}
		if n > 0 {
			statsData = append(statsData, readData[:n]...)
		} else {
			break
		}
	}

	statsClient.CloseIdleConnections()

	statsResp.Body.Close()

	// Read the statistics from the channel
	if plugin.StatisticsFormat == "xml" {
		// Read the XML statistics
		if err := ReadXmlStats(statsData); err != nil {
			return err
		}
	} else if plugin.StatisticsFormat == "json" {
		// Read the JSON statistics
		if err := ReadJsonStats(statsData); err != nil {
			return err
		}
	}

	return nil
}

func OutputMetricsGraphite() {
	// Output metrics in Graphite format
	for _, metric := range plugin.returnMetrics {
		fmt.Println(metric.Graphite("bind.dns"))
	}
}

func OutputMetricsPrometheus() {
	// Output metrics in Prometheus format
	for _, metric := range plugin.returnMetrics {
		fmt.Println(metric.Prometheus())
	}
}
