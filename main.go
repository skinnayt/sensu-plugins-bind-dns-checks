package main

import (
	"encoding/json"
	"encoding/xml"
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
			Usage:     "The format to output the metrics in (graphite)",
			Value:     &plugin.OutputFormat,
		},
	}
)

type MetricTag [2]string

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

func (m *Metric) String() string {
	return fmt.Sprintf("%s %s: %d", m.Timestamp, m.Name, m.Value)
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
		if err := readXmlStats(statsData); err != nil {
			return err
		}
	} else if plugin.StatisticsFormat == "json" {
		// Read the JSON statistics
		if err := readJsonStats(statsData); err != nil {
			return err
		}
	}

	return nil
}

type bindXmlStats struct {
	Stats []struct {
		Name  string `xml:"name,attr"`
		Value int    `xml:"value,attr"`
	} `xml:"statistics"`
}

type bindJsonStats struct {
	JsonStatsVersion string    `json:"json-stats-version"`
	BootTime         time.Time `json:"boot-time"`
	ConfigTime       time.Time `json:"config-time"`
	CurrentTime      time.Time `json:"current-time"`
	Version          string    `json:"version"`
	OpCodes          struct {
		Query      int `json:"QUERY"`
		IQuery     int `json:"IQUERY"`
		Status     int `json:"STATUS"`
		Reserved3  int `json:"RESERVED3"`
		Notify     int `json:"NOTIFY"`
		Update     int `json:"UPDATE"`
		Reserved6  int `json:"RESERVED6"`
		Reserved7  int `json:"RESERVED7"`
		Reserved8  int `json:"RESERVED8"`
		Reserved9  int `json:"RESERVED9"`
		Reserved10 int `json:"RESERVED10"`
		Reserved11 int `json:"RESERVED11"`
		Reserved12 int `json:"RESERVED12"`
		Reserved13 int `json:"RESERVED13"`
		Reserved14 int `json:"RESERVED14"`
		Reserved15 int `json:"RESERVED15"`
	} `json:"opcodes"`
	RCodes struct {
		Noerror    int `json:"NOERROR"`
		Formerr    int `json:"FORMERR"`
		Servfail   int `json:"SERVFAIL"`
		Nxdomain   int `json:"NXDOMAIN"`
		Notimp     int `json:"NOTIMP"`
		Refused    int `json:"REFUSED"`
		Yxdomain   int `json:"YXDOMAIN"`
		Yxrrset    int `json:"YXRRSET"`
		Nxrrset    int `json:"NXRRSET"`
		Notauth    int `json:"NOTAUTH"`
		Notzone    int `json:"NOTZONE"`
		Reserved11 int `json:"RESERVED11"`
		Reserved12 int `json:"RESERVED12"`
		Reserved13 int `json:"RESERVED13"`
		Reserved14 int `json:"RESERVED14"`
		Reserved15 int `json:"RESERVED15"`
		Badvers    int `json:"BADVERS"`
		R17        int `json:"17"`
		R18        int `json:"18"`
		R19        int `json:"19"`
		R20        int `json:"20"`
		R21        int `json:"21"`
		R22        int `json:"22"`
		Badcookie  int `json:"BADCOOKIE"`
	} `json:"rcodes"`
	QTypes struct {
		Others     int `json:"Others"`
		A          int `json:"A"`
		Ns         int `json:"NS"`
		Cname      int `json:"CNAME"`
		Soa        int `json:"SOA"`
		Ptr        int `json:"PTR"`
		Mx         int `json:"MX"`
		Txt        int `json:"TXT"`
		Afsdb      int `json:"AFSDB"`
		Aaaa       int `json:"AAAA"`
		Srv        int `json:"SRV"`
		Naptr      int `json:"NAPTR"`
		Dname      int `json:"DNAME"`
		Ds         int `json:"DS"`
		Rrsig      int `json:"RRSIG"`
		Dnskey     int `json:"DNSKEY"`
		Nsec3param int `json:"NSEC3PARAM"`
		Tlsa       int `json:"TLSA"`
		Cds        int `json:"CDS"`
		Cdnskey    int `json:"CDNSKEY"`
		Zonemd     int `json:"ZONEMD"`
		Svcb       int `json:"SVCB"`
		Https      int `json:"HTTPS"`
		Spf        int `json:"SPF"`
		Any        int `json:"ANY"`
	} `json:"qtypes"`
	NSStats struct {
		Requestv4        int `json:"Requestv4"`
		Requestv6        int `json:"Requestv6"`
		ReqEdns0         int `json:"ReqEdns0"`
		ReqTCP           int `json:"ReqTCP"`
		TCPConnHighWater int `json:"TCPConnHighWater"`
		AuthQryRej       int `json:"AuthQryRej"`
		RecQryRej        int `json:"RecQryRej"`
		Response         int `json:"Response"`
		TruncatedResp    int `json:"TruncatedResp"`
		RespEDNS0        int `json:"RespEDNS0"`
		QrySuccess       int `json:"QrySuccess"`
		QryAuthAns       int `json:"QryAuthAns"`
		QryNoauthAns     int `json:"QryNoauthAns"`
		QryReferral      int `json:"QryReferral"`
		QryNxrrset       int `json:"QryNxrrset"`
		QryNXDOMAIN      int `json:"QryNXDOMAIN"`
		QryFailure       int `json:"QryFailure"`
		QryUDP           int `json:"QryUDP"`
		QryTCP           int `json:"QryTCP"`
		CookieIn         int `json:"CookieIn"`
		CookieNew        int `json:"CookieNew"`
		CookieMatch      int `json:"CookieMatch"`
		ECSOpt           int `json:"ECSOpt"`
	} `json:"nsstats"`
	ZoneStats struct {
		NotifyInv4 int `json:"NotifyInv4"`
		SOAOutv4   int `json:"SOAOutv4"`
		AXFRReqv4  int `json:"AXFRReqv4"`
		IXFRReqv4  int `json:"IXFRReqv4"`
		XfrSuccess int `json:"XfrSuccess"`
	} `json:"zonestats"`
	Views struct {
		Default struct {
			Zones    []*ZoneView `json:"zones"`
			Resolver struct {
				Stats struct {
					Queryv6         int `json:"Queryv6"`
					Responsev6      int `json:"Responsev6"`
					NXDOMAIN        int `json:"NXDOMAIN"`
					Truncated       int `json:"Truncated"`
					Retry           int `json:"Retry"`
					ValAttempt      int `json:"ValAttempt"`
					ValOk           int `json:"ValOk"`
					ValNegOk        int `json:"ValNegOk"`
					QryRTT100       int `json:"QryRTT100"`
					QryRTT500       int `json:"QryRTT500"`
					BucketSize      int `json:"BucketSize"`
					ClientCookieOut int `json:"ClientCookieOut"`
					ServerCookieOut int `json:"ServerCookieOut"`
					CookieIn        int `json:"CookieIn"`
					CookieClientOk  int `json:"CookieClientOk"`
					Priming         int `json:"Priming"`
				} `json:"stats"`
				QTypes struct {
					Others     int `json:"Others"`
					A          int `json:"A"`
					Ns         int `json:"NS"`
					Cname      int `json:"CNAME"`
					Soa        int `json:"SOA"`
					Ptr        int `json:"PTR"`
					Mx         int `json:"MX"`
					Txt        int `json:"TXT"`
					Afsdb      int `json:"AFSDB"`
					Aaaa       int `json:"AAAA"`
					Srv        int `json:"SRV"`
					Naptr      int `json:"NAPTR"`
					Dname      int `json:"DNAME"`
					Ds         int `json:"DS"`
					Rrsig      int `json:"RRSIG"`
					Dnskey     int `json:"DNSKEY"`
					Nsec3param int `json:"NSEC3PARAM"`
					Tlsa       int `json:"TLSA"`
					Cds        int `json:"CDS"`
					Cdnskey    int `json:"CDNSKEY"`
					Zonemd     int `json:"ZONEMD"`
					Svcb       int `json:"SVCB"`
					Https      int `json:"HTTPS"`
					Spf        int `json:"SPF"`
					Any        int `json:"ANY"`
				} `json:"qtypes"`
				Cache struct {
					Others     int `json:"Others"`
					A          int `json:"A"`
					Ns         int `json:"NS"`
					Cname      int `json:"CNAME"`
					Soa        int `json:"SOA"`
					Ptr        int `json:"PTR"`
					Mx         int `json:"MX"`
					Txt        int `json:"TXT"`
					Afsdb      int `json:"AFSDB"`
					Aaaa       int `json:"AAAA"`
					Srv        int `json:"SRV"`
					Naptr      int `json:"NAPTR"`
					Dname      int `json:"DNAME"`
					Ds         int `json:"DS"`
					Rrsig      int `json:"RRSIG"`
					Dnskey     int `json:"DNSKEY"`
					Nsec3param int `json:"NSEC3PARAM"`
					Tlsa       int `json:"TLSA"`
					Cds        int `json:"CDS"`
					Cdnskey    int `json:"CDNSKEY"`
					Zonemd     int `json:"ZONEMD"`
					Svcb       int `json:"SVCB"`
					Https      int `json:"HTTPS"`
					Spf        int `json:"SPF"`
					Any        int `json:"ANY"`
				} `json:"cache"`
				CacheStats struct {
					CacheHits    int `json:"CacheHits"`
					CacheMisses  int `json:"CacheMisses"`
					QueryHits    int `json:"QueryHits"`
					QueryMisses  int `json:"QueryMisses"`
					DeleteLRU    int `json:"DeleteLRU"`
					DeleteTTL    int `json:"DeleteTTL"`
					CacheNodes   int `json:"CacheNodes"`
					CacheBuckets int `json:"CacheBuckets"`
					TreeMemTotal int `json:"TreeMemTotal"`
					TreeMemInUse int `json:"TreeMemInUse"`
					TreeMemMax   int `json:"TreeMemMax"`
					HeapMemTotal int `json:"HeapMemTotal"`
					HeapMemInUse int `json:"HeapMemInUse"`
					HeapMemMax   int `json:"HeapMemMax"`
				} `json:"cache-stats"`
				Adb struct {
					Nentries   int `json:"nentries"`
					Entriescnt int `json:"entriescnt"`
					Nnames     int `json:"nnames"`
					Namescnt   int `json:"namescnt"`
				} `json:"adb"`
			} `json:"resolver"`
		} `json:"default"`
		Bind struct {
			Zones    []*ZoneView `json:"zones"`
			Resolver struct {
				Stats struct {
					Queryv6         int `json:"Queryv6"`
					Responsev6      int `json:"Responsev6"`
					NXDOMAIN        int `json:"NXDOMAIN"`
					Truncated       int `json:"Truncated"`
					Retry           int `json:"Retry"`
					ValAttempt      int `json:"ValAttempt"`
					ValOk           int `json:"ValOk"`
					ValNegOk        int `json:"ValNegOk"`
					QryRTT100       int `json:"QryRTT100"`
					QryRTT500       int `json:"QryRTT500"`
					BucketSize      int `json:"BucketSize"`
					ClientCookieOut int `json:"ClientCookieOut"`
					ServerCookieOut int `json:"ServerCookieOut"`
					CookieIn        int `json:"CookieIn"`
					CookieClientOk  int `json:"CookieClientOk"`
					Priming         int `json:"Priming"`
				} `json:"stats"`
				QTypes struct {
					Others     int `json:"Others"`
					A          int `json:"A"`
					Ns         int `json:"NS"`
					Cname      int `json:"CNAME"`
					Soa        int `json:"SOA"`
					Ptr        int `json:"PTR"`
					Mx         int `json:"MX"`
					Txt        int `json:"TXT"`
					Afsdb      int `json:"AFSDB"`
					Aaaa       int `json:"AAAA"`
					Srv        int `json:"SRV"`
					Naptr      int `json:"NAPTR"`
					Dname      int `json:"DNAME"`
					Ds         int `json:"DS"`
					Rrsig      int `json:"RRSIG"`
					Dnskey     int `json:"DNSKEY"`
					Nsec3param int `json:"NSEC3PARAM"`
					Tlsa       int `json:"TLSA"`
					Cds        int `json:"CDS"`
					Cdnskey    int `json:"CDNSKEY"`
					Zonemd     int `json:"ZONEMD"`
					Svcb       int `json:"SVCB"`
					Https      int `json:"HTTPS"`
					Spf        int `json:"SPF"`
					Any        int `json:"ANY"`
				} `json:"qtypes"`
				Cache struct {
					Others     int `json:"Others"`
					A          int `json:"A"`
					Ns         int `json:"NS"`
					Cname      int `json:"CNAME"`
					Soa        int `json:"SOA"`
					Ptr        int `json:"PTR"`
					Mx         int `json:"MX"`
					Txt        int `json:"TXT"`
					Afsdb      int `json:"AFSDB"`
					Aaaa       int `json:"AAAA"`
					Srv        int `json:"SRV"`
					Naptr      int `json:"NAPTR"`
					Dname      int `json:"DNAME"`
					Ds         int `json:"DS"`
					Rrsig      int `json:"RRSIG"`
					Dnskey     int `json:"DNSKEY"`
					Nsec3param int `json:"NSEC3PARAM"`
					Tlsa       int `json:"TLSA"`
					Cds        int `json:"CDS"`
					Cdnskey    int `json:"CDNSKEY"`
					Zonemd     int `json:"ZONEMD"`
					Svcb       int `json:"SVCB"`
					Https      int `json:"HTTPS"`
					Spf        int `json:"SPF"`
					Any        int `json:"ANY"`
				} `json:"cache"`
				CacheStats struct {
					CacheHits    int `json:"CacheHits"`
					CacheMisses  int `json:"CacheMisses"`
					QueryHits    int `json:"QueryHits"`
					QueryMisses  int `json:"QueryMisses"`
					DeleteLRU    int `json:"DeleteLRU"`
					DeleteTTL    int `json:"DeleteTTL"`
					CacheNodes   int `json:"CacheNodes"`
					CacheBuckets int `json:"CacheBuckets"`
					TreeMemTotal int `json:"TreeMemTotal"`
					TreeMemInUse int `json:"TreeMemInUse"`
					TreeMemMax   int `json:"TreeMemMax"`
					HeapMemTotal int `json:"HeapMemTotal"`
					HeapMemInUse int `json:"HeapMemInUse"`
					HeapMemMax   int `json:"HeapMemMax"`
				} `json:"cache-stats"`
				Adb struct {
					Nentries   int `json:"nentries"`
					Entriescnt int `json:"entriescnt"`
					Nnames     int `json:"nnames"`
					Namescnt   int `json:"namescnt"`
				} `json:"adb"`
			} `json:"resolver"`
		} `json:"bind"`
	} `json:"views"`
	SocketStats struct {
		UDP4Open    int `json:"UDP4Open"`
		UDP6Open    int `json:"UDP6Open"`
		TCP4Open    int `json:"TCP4Open"`
		TCP6Open    int `json:"TCP6Open"`
		RawOpen     int `json:"RawOpen"`
		UDP4Close   int `json:"UDP4Close"`
		UDP6Close   int `json:"UDP6Close"`
		TCP4Close   int `json:"TCP4Close"`
		TCP6Close   int `json:"TCP6Close"`
		UDP6Conn    int `json:"UDP6Conn"`
		TCP4Conn    int `json:"TCP4Conn"`
		TCP6Conn    int `json:"TCP6Conn"`
		TCP4Accept  int `json:"TCP4Accept"`
		TCP6Accept  int `json:"TCP6Accept"`
		TCP4RecvErr int `json:"TCP4RecvErr"`
		UDP4Active  int `json:"UDP4Active"`
		UDP6Active  int `json:"UDP6Active"`
		TCP4Active  int `json:"TCP4Active"`
		TCP6Active  int `json:"TCP6Active"`
		RawActive   int `json:"RawActive"`
	} `json:"socketstats"`
	SocketMgr struct {
		Sockets []SocketMgrSocket `json:"sockets"`
	} `json:"socketmgr"`
	TaskMgr struct {
		ThreadModel    string         `json:"thread-model"`
		DefaultQuantum int            `json:"default-quantum"`
		Tasks          []*TaskMgrTask `json:"tasks"`
	} `json:"taskmgr"`
	Memory struct {
		TotalUse    int        `json:"TotalUse"`
		InUse       int        `json:"InUse"`
		Malloced    int        `json:"Malloced"`
		BlockSize   int        `json:"BlockSize"`
		ContextSize int        `json:"ContextSize"`
		Lost        int        `json:"Lost"`
		Contexts    []*Context `json:"Contexts"`
	} `json:"memory"`
	Traffic struct {
		DnsUDPRequestsSizesReceivedIPv4 struct {
			U0_15    int `json:"0-15,omitempty"`
			U16_31   int `json:"16-31,omitempty"`
			U32_47   int `json:"32-47,omitempty"`
			U48_63   int `json:"48-63,omitempty"`
			U64_79   int `json:"64-79,omitempty"`
			U80_95   int `json:"80-95,omitempty"`
			U96_111  int `json:"96-111,omitempty"`
			U112_127 int `json:"112-127,omitempty"`
			U128_143 int `json:"128-143,omitempty"`
			U144_159 int `json:"144-159,omitempty"`
		} `json:"dns-udp-requests-sizes-received-ipv4"`
		DnsUDPResponsesSizesSentIPv4 struct {
			U0_15      int `json:"0-15,omitempty"`
			U16_31     int `json:"16-31,omitempty"`
			U32_47     int `json:"32-47,omitempty"`
			U48_63     int `json:"48-63,omitempty"`
			U64_79     int `json:"64-79,omitempty"`
			U80_95     int `json:"80-95,omitempty"`
			U96_111    int `json:"96-111,omitempty"`
			U112_127   int `json:"112-127,omitempty"`
			U128_143   int `json:"128-143,omitempty"`
			U144_159   int `json:"144-159,omitempty"`
			U160_175   int `json:"160-175,omitempty"`
			U176_191   int `json:"176-191,omitempty"`
			U192_207   int `json:"192-207,omitempty"`
			U208_223   int `json:"208-223,omitempty"`
			U224_239   int `json:"224-239,omitempty"`
			U240_255   int `json:"240-255,omitempty"`
			U256_271   int `json:"256-271,omitempty"`
			U272_287   int `json:"272-287,omitempty"`
			U288_303   int `json:"288-303,omitempty"`
			U304_319   int `json:"304-319,omitempty"`
			U320_335   int `json:"320-335,omitempty"`
			U336_351   int `json:"336-351,omitempty"`
			U352_367   int `json:"352-367,omitempty"`
			U368_383   int `json:"368-383,omitempty"`
			U384_399   int `json:"384-399,omitempty"`
			U400_415   int `json:"400-415,omitempty"`
			U416_431   int `json:"416-431,omitempty"`
			U432_447   int `json:"432-447,omitempty"`
			U448_463   int `json:"448-463,omitempty"`
			U464_479   int `json:"464-479,omitempty"`
			U480_495   int `json:"480-495,omitempty"`
			U496_511   int `json:"496-511,omitempty"`
			U512_527   int `json:"512-527,omitempty"`
			U528_543   int `json:"528-543,omitempty"`
			U544_559   int `json:"544-559,omitempty"`
			U560_575   int `json:"560-575,omitempty"`
			U576_591   int `json:"576-591,omitempty"`
			U592_607   int `json:"592-607,omitempty"`
			U608_623   int `json:"608-623,omitempty"`
			U624_639   int `json:"624-639,omitempty"`
			U640_655   int `json:"640-655,omitempty"`
			U656_671   int `json:"656-671,omitempty"`
			U672_687   int `json:"672-687,omitempty"`
			U688_703   int `json:"688-703,omitempty"`
			U704_719   int `json:"704-719,omitempty"`
			U720_735   int `json:"720-735,omitempty"`
			U736_751   int `json:"736-751,omitempty"`
			U752_767   int `json:"752-767,omitempty"`
			U768_783   int `json:"768-783,omitempty"`
			U784_799   int `json:"784-799,omitempty"`
			U800_815   int `json:"800-815,omitempty"`
			U816_831   int `json:"816-831,omitempty"`
			U832_847   int `json:"832-847,omitempty"`
			U848_863   int `json:"848-863,omitempty"`
			U864_879   int `json:"864-879,omitempty"`
			U880_895   int `json:"880-895,omitempty"`
			U896_911   int `json:"896-911,omitempty"`
			U912_927   int `json:"912-927,omitempty"`
			U928_943   int `json:"928-943,omitempty"`
			U944_959   int `json:"944-959,omitempty"`
			U960_975   int `json:"960-975,omitempty"`
			U976_991   int `json:"976-991,omitempty"`
			U992_1007  int `json:"992-1007,omitempty"`
			U1008_1023 int `json:"1008-1023,omitempty"`
			U1024_1039 int `json:"1024-1039,omitempty"`
			U1040_1055 int `json:"1040-1055,omitempty"`
			U1056_1071 int `json:"1056-1071,omitempty"`
			U1072_1087 int `json:"1072-1087,omitempty"`
			U1088_1103 int `json:"1088-1103,omitempty"`
			U1104_1119 int `json:"1104-1119,omitempty"`
			U1120_1135 int `json:"1120-1135,omitempty"`
			U1136_1151 int `json:"1136-1151,omitempty"`
			U1152_1167 int `json:"1152-1167,omitempty"`
			U1168_1183 int `json:"1168-1183,omitempty"`
			U1184_1199 int `json:"1184-1199,omitempty"`
			U1200_1215 int `json:"1200-1215,omitempty"`
			U1216_1231 int `json:"1216-1231,omitempty"`
		} `json:"dns-udp-responses-sizes-sent-ipv4"`
		DnsTCPRequestsSizesReceivedIPv4 struct {
			T0_15  int `json:"0-15,omitempty"`
			T16_31 int `json:"16-31,omitempty"`
			T32_47 int `json:"32-47,omitempty"`
			T48_63 int `json:"48-63,omitempty"`
			T64_79 int `json:"64-79,omitempty"`
			T80_95 int `json:"80-95,omitempty"`
		} `json:"dns-tcp-requests-sizes-received-ipv4"`
		DnsTCPResponsesSizesSentIPv4 struct {
			T0_15      int `json:"0-15,omitempty"`
			T16_31     int `json:"16-31,omitempty"`
			T32_47     int `json:"32-47,omitempty"`
			T48_63     int `json:"48-63,omitempty"`
			T64_79     int `json:"64-79,omitempty"`
			T80_95     int `json:"80-95,omitempty"`
			T96_111    int `json:"96-111,omitempty"`
			T112_127   int `json:"112-127,omitempty"`
			T128_143   int `json:"128-143,omitempty"`
			T144_159   int `json:"144-159,omitempty"`
			T160_175   int `json:"160-175,omitempty"`
			T176_191   int `json:"176-191,omitempty"`
			T192_207   int `json:"192-207,omitempty"`
			T208_223   int `json:"208-223,omitempty"`
			T224_239   int `json:"224-239,omitempty"`
			T240_255   int `json:"240-255,omitempty"`
			T256_271   int `json:"256-271,omitempty"`
			T272_287   int `json:"272-287,omitempty"`
			T288_303   int `json:"288-303,omitempty"`
			T304_319   int `json:"304-319,omitempty"`
			T320_335   int `json:"320-335,omitempty"`
			T336_351   int `json:"336-351,omitempty"`
			T352_367   int `json:"352-367,omitempty"`
			T368_383   int `json:"368-383,omitempty"`
			T384_399   int `json:"384-399,omitempty"`
			T400_415   int `json:"400-415,omitempty"`
			T416_431   int `json:"416-431,omitempty"`
			T432_447   int `json:"432-447,omitempty"`
			T448_463   int `json:"448-463,omitempty"`
			T464_479   int `json:"464-479,omitempty"`
			T480_495   int `json:"480-495,omitempty"`
			T496_511   int `json:"496-511,omitempty"`
			T512_527   int `json:"512-527,omitempty"`
			T528_543   int `json:"528-543,omitempty"`
			T544_559   int `json:"544-559,omitempty"`
			T560_575   int `json:"560-575,omitempty"`
			T576_591   int `json:"576-591,omitempty"`
			T592_607   int `json:"592-607,omitempty"`
			T608_623   int `json:"608-623,omitempty"`
			T624_639   int `json:"624-639,omitempty"`
			T640_655   int `json:"640-655,omitempty"`
			T656_671   int `json:"656-671,omitempty"`
			T672_687   int `json:"672-687,omitempty"`
			T688_703   int `json:"688-703,omitempty"`
			T704_719   int `json:"704-719,omitempty"`
			T720_735   int `json:"720-735,omitempty"`
			T736_751   int `json:"736-751,omitempty"`
			T752_767   int `json:"752-767,omitempty"`
			T768_783   int `json:"768-783,omitempty"`
			T784_799   int `json:"784-799,omitempty"`
			T800_815   int `json:"800-815,omitempty"`
			T816_831   int `json:"816-831,omitempty"`
			T832_847   int `json:"832-847,omitempty"`
			T848_863   int `json:"848-863,omitempty"`
			T864_879   int `json:"864-879,omitempty"`
			T880_895   int `json:"880-895,omitempty"`
			T896_911   int `json:"896-911,omitempty"`
			T912_927   int `json:"912-927,omitempty"`
			T928_943   int `json:"928-943,omitempty"`
			T944_959   int `json:"944-959,omitempty"`
			T960_975   int `json:"960-975,omitempty"`
			T976_991   int `json:"976-991,omitempty"`
			T992_1007  int `json:"992-1007,omitempty"`
			T1008_1023 int `json:"1008-1023,omitempty"`
			T1024_1039 int `json:"1024-1039,omitempty"`
			T1040_1055 int `json:"1040-1055,omitempty"`
			T1056_1071 int `json:"1056-1071,omitempty"`
			T1072_1087 int `json:"1072-1087,omitempty"`
			T1088_1103 int `json:"1088-1103,omitempty"`
			T1104_1119 int `json:"1104-1119,omitempty"`
			T1120_1135 int `json:"1120-1135,omitempty"`
			T1136_1151 int `json:"1136-1151,omitempty"`
			T1152_1167 int `json:"1152-1167,omitempty"`
			T1168_1183 int `json:"1168-1183,omitempty"`
			T1184_1199 int `json:"1184-1199,omitempty"`
			T1200_1215 int `json:"1200-1215,omitempty"`
			T1216_1231 int `json:"1216-1231,omitempty"`
			T1232_1247 int `json:"1232-1247,omitempty"`
			T1248_1263 int `json:"1248-1263,omitempty"`
			T1264_1279 int `json:"1264-1279,omitempty"`
			T1280_1295 int `json:"1280-1295,omitempty"`
			T1296_1311 int `json:"1296-1311,omitempty"`
			T1312_1327 int `json:"1312-1327,omitempty"`
			T1328_1343 int `json:"1328-1343,omitempty"`
			T1344_1359 int `json:"1344-1359,omitempty"`
			T1360_1375 int `json:"1360-1375,omitempty"`
			T1376_1391 int `json:"1376-1391,omitempty"`
			T1392_1407 int `json:"1392-1407,omitempty"`
			T1408_1423 int `json:"1408-1423,omitempty"`
			T1424_1439 int `json:"1424-1439,omitempty"`
			T1440_1455 int `json:"1440-1455,omitempty"`
			T1456_1471 int `json:"1456-1471,omitempty"`
			T1472_1487 int `json:"1472-1487,omitempty"`
			T1488_1503 int `json:"1488-1503,omitempty"`
			T1504_1519 int `json:"1504-1519,omitempty"`
			T1520_1535 int `json:"1520-1535,omitempty"`
			T1536_1551 int `json:"1536-1551,omitempty"`
			T1552_1567 int `json:"1552-1567,omitempty"`
			T1568_1583 int `json:"1568-1583,omitempty"`
			T1584_1599 int `json:"1584-1599,omitempty"`
			T1600_1615 int `json:"1600-1615,omitempty"`
			T1616_1631 int `json:"1616-1631,omitempty"`
			T1632_1647 int `json:"1632-1647,omitempty"`
			T1648_1663 int `json:"1648-1663,omitempty"`
			T1664_1679 int `json:"1664-1679,omitempty"`
			T1680_1695 int `json:"1680-1695,omitempty"`
			T1696_1711 int `json:"1696-1711,omitempty"`
			T1712_1727 int `json:"1712-1727,omitempty"`
			T1728_1743 int `json:"1728-1743,omitempty"`
			T1744_1759 int `json:"1744-1759,omitempty"`
			T1760_1775 int `json:"1760-1775,omitempty"`
			T1776_1791 int `json:"1776-1791,omitempty"`
			T1792_1807 int `json:"1792-1807,omitempty"`
			T1808_1823 int `json:"1808-1823,omitempty"`
			T1824_1839 int `json:"1824-1839,omitempty"`
			T1840_1855 int `json:"1840-1855,omitempty"`
			T1856_1871 int `json:"1856-1871,omitempty"`
			T1872_1887 int `json:"1872-1887,omitempty"`
			T1888_1903 int `json:"1888-1903,omitempty"`
			T1904_1919 int `json:"1904-1919,omitempty"`
			T1920_1935 int `json:"1920-1935,omitempty"`
			T1936_1951 int `json:"1936-1951,omitempty"`
			T1952_1967 int `json:"1952-1967,omitempty"`
			T1968_1983 int `json:"1968-1983,omitempty"`
			T1984_1999 int `json:"1984-1999,omitempty"`
			T2000_2015 int `json:"2000-2015,omitempty"`
			T2016_2031 int `json:"2016-2031,omitempty"`
			T2032_2047 int `json:"2032-2047,omitempty"`
			T2048_2063 int `json:"2048-2063,omitempty"`
			T2064_2079 int `json:"2064-2079,omitempty"`
			T2080_2095 int `json:"2080-2095,omitempty"`
			T2096_2111 int `json:"2096-2111,omitempty"`
			T2112_2127 int `json:"2112-2127,omitempty"`
			T2128_2143 int `json:"2128-2143,omitempty"`
			T2144_2159 int `json:"2144-2159,omitempty"`
			T2160_2175 int `json:"2160-2175,omitempty"`
			T2176_2191 int `json:"2176-2191,omitempty"`
			T2192_2207 int `json:"2192-2207,omitempty"`
			T2208_2223 int `json:"2208-2223,omitempty"`
			T2224_2239 int `json:"2224-2239,omitempty"`
			T2240_2255 int `json:"2240-2255,omitempty"`
			T2256_2271 int `json:"2256-2271,omitempty"`
			T2272_2287 int `json:"2272-2287,omitempty"`
			T2288_2303 int `json:"2288-2303,omitempty"`
			T2304_2319 int `json:"2304-2319,omitempty"`
			T2320_2335 int `json:"2320-2335,omitempty"`
		} `json:"dns-tcp-responses-sizes-sent-ipv4"`
		DnsUDPRequestsSizesReceivedIPv6 struct {
			U0_15    int `json:"0-15,omitempty"`
			U16_31   int `json:"16-31,omitempty"`
			U32_47   int `json:"32-47,omitempty"`
			U48_63   int `json:"48-63,omitempty"`
			U64_79   int `json:"64-79,omitempty"`
			U80_95   int `json:"80-95,omitempty"`
			U96_111  int `json:"96-111,omitempty"`
			U112_127 int `json:"112-127,omitempty"`
			U128_143 int `json:"128-143,omitempty"`
			U144_159 int `json:"144-159,omitempty"`
		} `json:"dns-udp-requests-sizes-received-ipv6"`
		DnsUDPResponsesSizesSentIPv6 struct {
			U0_15      int `json:"0-15,omitempty"`
			U16_31     int `json:"16-31,omitempty"`
			U32_47     int `json:"32-47,omitempty"`
			U48_63     int `json:"48-63,omitempty"`
			U64_79     int `json:"64-79,omitempty"`
			U80_95     int `json:"80-95,omitempty"`
			U96_111    int `json:"96-111,omitempty"`
			U112_127   int `json:"112-127,omitempty"`
			U128_143   int `json:"128-143,omitempty"`
			U144_159   int `json:"144-159,omitempty"`
			U160_175   int `json:"160-175,omitempty"`
			U176_191   int `json:"176-191,omitempty"`
			U192_207   int `json:"192-207,omitempty"`
			U208_223   int `json:"208-223,omitempty"`
			U224_239   int `json:"224-239,omitempty"`
			U240_255   int `json:"240-255,omitempty"`
			U256_271   int `json:"256-271,omitempty"`
			U272_287   int `json:"272-287,omitempty"`
			U288_303   int `json:"288-303,omitempty"`
			U304_319   int `json:"304-319,omitempty"`
			U320_335   int `json:"320-335,omitempty"`
			U336_351   int `json:"336-351,omitempty"`
			U352_367   int `json:"352-367,omitempty"`
			U368_383   int `json:"368-383,omitempty"`
			U384_399   int `json:"384-399,omitempty"`
			U400_415   int `json:"400-415,omitempty"`
			U416_431   int `json:"416-431,omitempty"`
			U432_447   int `json:"432-447,omitempty"`
			U448_463   int `json:"448-463,omitempty"`
			U464_479   int `json:"464-479,omitempty"`
			U480_495   int `json:"480-495,omitempty"`
			U496_511   int `json:"496-511,omitempty"`
			U512_527   int `json:"512-527,omitempty"`
			U528_543   int `json:"528-543,omitempty"`
			U544_559   int `json:"544-559,omitempty"`
			U560_575   int `json:"560-575,omitempty"`
			U576_591   int `json:"576-591,omitempty"`
			U592_607   int `json:"592-607,omitempty"`
			U608_623   int `json:"608-623,omitempty"`
			U624_639   int `json:"624-639,omitempty"`
			U640_655   int `json:"640-655,omitempty"`
			U656_671   int `json:"656-671,omitempty"`
			U672_687   int `json:"672-687,omitempty"`
			U688_703   int `json:"688-703,omitempty"`
			U704_719   int `json:"704-719,omitempty"`
			U720_735   int `json:"720-735,omitempty"`
			U736_751   int `json:"736-751,omitempty"`
			U752_767   int `json:"752-767,omitempty"`
			U768_783   int `json:"768-783,omitempty"`
			U784_799   int `json:"784-799,omitempty"`
			U800_815   int `json:"800-815,omitempty"`
			U816_831   int `json:"816-831,omitempty"`
			U832_847   int `json:"832-847,omitempty"`
			U848_863   int `json:"848-863,omitempty"`
			U864_879   int `json:"864-879,omitempty"`
			U880_895   int `json:"880-895,omitempty"`
			U896_911   int `json:"896-911,omitempty"`
			U912_927   int `json:"912-927,omitempty"`
			U928_943   int `json:"928-943,omitempty"`
			U944_959   int `json:"944-959,omitempty"`
			U960_975   int `json:"960-975,omitempty"`
			U976_991   int `json:"976-991,omitempty"`
			U992_1007  int `json:"992-1007,omitempty"`
			U1008_1023 int `json:"1008-1023,omitempty"`
			U1024_1039 int `json:"1024-1039,omitempty"`
			U1040_1055 int `json:"1040-1055,omitempty"`
			U1056_1071 int `json:"1056-1071,omitempty"`
			U1072_1087 int `json:"1072-1087,omitempty"`
			U1088_1103 int `json:"1088-1103,omitempty"`
			U1104_1119 int `json:"1104-1119,omitempty"`
			U1120_1135 int `json:"1120-1135,omitempty"`
			U1136_1151 int `json:"1136-1151,omitempty"`
			U1152_1167 int `json:"1152-1167,omitempty"`
			U1168_1183 int `json:"1168-1183,omitempty"`
			U1184_1199 int `json:"1184-1199,omitempty"`
			U1200_1215 int `json:"1200-1215,omitempty"`
			U1216_1231 int `json:"1216-1231,omitempty"`
		} `json:"dns-udp-responses-sizes-sent-ipv6"`
		DnsTCPRequestsSizesReceivedIPv6 struct {
			T0_15  int `json:"0-15,omitempty"`
			T16_31 int `json:"16-31,omitempty"`
			T32_47 int `json:"32-47,omitempty"`
			T48_63 int `json:"48-63,omitempty"`
			T64_79 int `json:"64-79,omitempty"`
			T80_95 int `json:"80-95,omitempty"`
		} `json:"dns-tcp-requests-sizes-received-ipv6"`
		DnsTCPResponsesSizesSentIPv6 struct {
			T0_15      int `json:"0-15,omitempty"`
			T16_31     int `json:"16-31,omitempty"`
			T32_47     int `json:"32-47,omitempty"`
			T48_63     int `json:"48-63,omitempty"`
			T64_79     int `json:"64-79,omitempty"`
			T80_95     int `json:"80-95,omitempty"`
			T96_111    int `json:"96-111,omitempty"`
			T112_127   int `json:"112-127,omitempty"`
			T128_143   int `json:"128-143,omitempty"`
			T144_159   int `json:"144-159,omitempty"`
			T160_175   int `json:"160-175,omitempty"`
			T176_191   int `json:"176-191,omitempty"`
			T192_207   int `json:"192-207,omitempty"`
			T208_223   int `json:"208-223,omitempty"`
			T224_239   int `json:"224-239,omitempty"`
			T240_255   int `json:"240-255,omitempty"`
			T256_271   int `json:"256-271,omitempty"`
			T272_287   int `json:"272-287,omitempty"`
			T288_303   int `json:"288-303,omitempty"`
			T304_319   int `json:"304-319,omitempty"`
			T320_335   int `json:"320-335,omitempty"`
			T336_351   int `json:"336-351,omitempty"`
			T352_367   int `json:"352-367,omitempty"`
			T368_383   int `json:"368-383,omitempty"`
			T384_399   int `json:"384-399,omitempty"`
			T400_415   int `json:"400-415,omitempty"`
			T416_431   int `json:"416-431,omitempty"`
			T432_447   int `json:"432-447,omitempty"`
			T448_463   int `json:"448-463,omitempty"`
			T464_479   int `json:"464-479,omitempty"`
			T480_495   int `json:"480-495,omitempty"`
			T496_511   int `json:"496-511,omitempty"`
			T512_527   int `json:"512-527,omitempty"`
			T528_543   int `json:"528-543,omitempty"`
			T544_559   int `json:"544-559,omitempty"`
			T560_575   int `json:"560-575,omitempty"`
			T576_591   int `json:"576-591,omitempty"`
			T592_607   int `json:"592-607,omitempty"`
			T608_623   int `json:"608-623,omitempty"`
			T624_639   int `json:"624-639,omitempty"`
			T640_655   int `json:"640-655,omitempty"`
			T656_671   int `json:"656-671,omitempty"`
			T672_687   int `json:"672-687,omitempty"`
			T688_703   int `json:"688-703,omitempty"`
			T704_719   int `json:"704-719,omitempty"`
			T720_735   int `json:"720-735,omitempty"`
			T736_751   int `json:"736-751,omitempty"`
			T752_767   int `json:"752-767,omitempty"`
			T768_783   int `json:"768-783,omitempty"`
			T784_799   int `json:"784-799,omitempty"`
			T800_815   int `json:"800-815,omitempty"`
			T816_831   int `json:"816-831,omitempty"`
			T832_847   int `json:"832-847,omitempty"`
			T848_863   int `json:"848-863,omitempty"`
			T864_879   int `json:"864-879,omitempty"`
			T880_895   int `json:"880-895,omitempty"`
			T896_911   int `json:"896-911,omitempty"`
			T912_927   int `json:"912-927,omitempty"`
			T928_943   int `json:"928-943,omitempty"`
			T944_959   int `json:"944-959,omitempty"`
			T960_975   int `json:"960-975,omitempty"`
			T976_991   int `json:"976-991,omitempty"`
			T992_1007  int `json:"992-1007,omitempty"`
			T1008_1023 int `json:"1008-1023,omitempty"`
			T1024_1039 int `json:"1024-1039,omitempty"`
			T1040_1055 int `json:"1040-1055,omitempty"`
			T1056_1071 int `json:"1056-1071,omitempty"`
			T1072_1087 int `json:"1072-1087,omitempty"`
			T1088_1103 int `json:"1088-1103,omitempty"`
			T1104_1119 int `json:"1104-1119,omitempty"`
			T1120_1135 int `json:"1120-1135,omitempty"`
			T1136_1151 int `json:"1136-1151,omitempty"`
			T1152_1167 int `json:"1152-1167,omitempty"`
			T1168_1183 int `json:"1168-1183,omitempty"`
			T1184_1199 int `json:"1184-1199,omitempty"`
			T1200_1215 int `json:"1200-1215,omitempty"`
			T1216_1231 int `json:"1216-1231,omitempty"`
			T1232_1247 int `json:"1232-1247,omitempty"`
			T1248_1263 int `json:"1248-1263,omitempty"`
			T1264_1279 int `json:"1264-1279,omitempty"`
			T1280_1295 int `json:"1280-1295,omitempty"`
			T1296_1311 int `json:"1296-1311,omitempty"`
			T1312_1327 int `json:"1312-1327,omitempty"`
			T1328_1343 int `json:"1328-1343,omitempty"`
			T1344_1359 int `json:"1344-1359,omitempty"`
			T1360_1375 int `json:"1360-1375,omitempty"`
			T1376_1391 int `json:"1376-1391,omitempty"`
			T1392_1407 int `json:"1392-1407,omitempty"`
			T1408_1423 int `json:"1408-1423,omitempty"`
			T1424_1439 int `json:"1424-1439,omitempty"`
			T1440_1455 int `json:"1440-1455,omitempty"`
			T1456_1471 int `json:"1456-1471,omitempty"`
			T1472_1487 int `json:"1472-1487,omitempty"`
			T1488_1503 int `json:"1488-1503,omitempty"`
			T1504_1519 int `json:"1504-1519,omitempty"`
			T1520_1535 int `json:"1520-1535,omitempty"`
			T1536_1551 int `json:"1536-1551,omitempty"`
			T1552_1567 int `json:"1552-1567,omitempty"`
			T1568_1583 int `json:"1568-1583,omitempty"`
			T1584_1599 int `json:"1584-1599,omitempty"`
			T1600_1615 int `json:"1600-1615,omitempty"`
			T1616_1631 int `json:"1616-1631,omitempty"`
			T1632_1647 int `json:"1632-1647,omitempty"`
			T1648_1663 int `json:"1648-1663,omitempty"`
			T1664_1679 int `json:"1664-1679,omitempty"`
			T1680_1695 int `json:"1680-1695,omitempty"`
			T1696_1711 int `json:"1696-1711,omitempty"`
			T1712_1727 int `json:"1712-1727,omitempty"`
			T1728_1743 int `json:"1728-1743,omitempty"`
			T1744_1759 int `json:"1744-1759,omitempty"`
			T1760_1775 int `json:"1760-1775,omitempty"`
			T1776_1791 int `json:"1776-1791,omitempty"`
			T1792_1807 int `json:"1792-1807,omitempty"`
			T1808_1823 int `json:"1808-1823,omitempty"`
			T1824_1839 int `json:"1824-1839,omitempty"`
			T1840_1855 int `json:"1840-1855,omitempty"`
			T1856_1871 int `json:"1856-1871,omitempty"`
			T1872_1887 int `json:"1872-1887,omitempty"`
			T1888_1903 int `json:"1888-1903,omitempty"`
			T1904_1919 int `json:"1904-1919,omitempty"`
			T1920_1935 int `json:"1920-1935,omitempty"`
			T1936_1951 int `json:"1936-1951,omitempty"`
			T1952_1967 int `json:"1952-1967,omitempty"`
			T1968_1983 int `json:"1968-1983,omitempty"`
			T1984_1999 int `json:"1984-1999,omitempty"`
			T2000_2015 int `json:"2000-2015,omitempty"`
			T2016_2031 int `json:"2016-2031,omitempty"`
			T2032_2047 int `json:"2032-2047,omitempty"`
			T2048_2063 int `json:"2048-2063,omitempty"`
			T2064_2079 int `json:"2064-2079,omitempty"`
			T2080_2095 int `json:"2080-2095,omitempty"`
			T2096_2111 int `json:"2096-2111,omitempty"`
			T2112_2127 int `json:"2112-2127,omitempty"`
			T2128_2143 int `json:"2128-2143,omitempty"`
			T2144_2159 int `json:"2144-2159,omitempty"`
			T2160_2175 int `json:"2160-2175,omitempty"`
			T2176_2191 int `json:"2176-2191,omitempty"`
			T2192_2207 int `json:"2192-2207,omitempty"`
			T2208_2223 int `json:"2208-2223,omitempty"`
			T2224_2239 int `json:"2224-2239,omitempty"`
			T2240_2255 int `json:"2240-2255,omitempty"`
			T2256_2271 int `json:"2256-2271,omitempty"`
			T2272_2287 int `json:"2272-2287,omitempty"`
			T2288_2303 int `json:"2288-2303,omitempty"`
			T2304_2319 int `json:"2304-2319,omitempty"`
			T2320_2335 int `json:"2320-2335,omitempty"`
		} `json:"dns-tcp-responses-sizes-sent-ipv6"`
	} `json:"traffic"`
}

type SocketMgrSocket struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	References   int      `json:"references"`
	Type         string   `json:"type"`
	PeerAddress  string   `json:"peer_address,omitempty"`
	LocalAddress string   `json:"local_address,omitempty"`
	States       []string `json:"states"`
}

type TaskMgrTask struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	References int    `json:"references"`
	State      string `json:"state"`
	Quantun    int    `json:"quantum"`
	Events     int    `json:"events"`
}

type Context struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	References  int    `json:"references"`
	Total       int    `json:"total"`
	Inuse       int    `json:"inuse"`
	Maxinuse    int    `json:"maxinuse"`
	Malloced    int    `json:"malloced"`
	Maxmalloced int    `json:"maxmalloced"`
	Blocksize   int    `json:"blocksize"`
	Pools       int    `json:"pools"`
	Hiwater     int    `json:"hiwater"`
	Lowater     int    `json:"lowater"`
}

type ZoneView struct {
	Name    string    `json:"name"`
	Class   string    `json:"class"`
	Serial  int       `json:"serial"`
	Type    string    `json:"type"`
	Loaded  time.Time `json:"loaded"`
	Expires time.Time `json:"expires,omitempty"`
	Refresh time.Time `json:"refresh,omitempty"`
	RCodes  struct {
		Noerror    int `json:"NOERROR,omitempty"`
		Formerr    int `json:"FORMERR,omitempty"`
		Servfail   int `json:"SERVFAIL,omitempty"`
		Nxdomain   int `json:"NXDOMAIN,omitempty"`
		Notimp     int `json:"NOTIMP,omitempty"`
		Refused    int `json:"REFUSED,omitempty"`
		Yxdomain   int `json:"YXDOMAIN,omitempty"`
		Yxrrset    int `json:"YXRRSET,omitempty"`
		Nxrrset    int `json:"NXRRSET,omitempty"`
		Notauth    int `json:"NOTAUTH,omitempty"`
		Notzone    int `json:"NOTZONE,omitempty"`
		Reserved11 int `json:"RESERVED11,omitempty"`
		Reserved12 int `json:"RESERVED12,omitempty"`
		Reserved13 int `json:"RESERVED13,omitempty"`
		Reserved14 int `json:"RESERVED14,omitempty"`
		Reserved15 int `json:"RESERVED15,omitempty"`
		Badvers    int `json:"BADVERS,omitempty"`
		R17        int `json:"17,omitempty"`
		R18        int `json:"18,omitempty"`
		R19        int `json:"19,omitempty"`
		R20        int `json:"20,omitempty"`
		R21        int `json:"21,omitempty"`
		R22        int `json:"22,omitempty"`
		Badcookie  int `json:"BADCOOKIE,omitempty"`
	} `json:"rcodes,omitempty"`
	QTypes struct {
		Others     int `json:"Others,omitempty"`
		A          int `json:"A,omitempty"`
		Ns         int `json:"NS,omitempty"`
		Cname      int `json:"CNAME,omitempty"`
		Soa        int `json:"SOA,omitempty"`
		Ptr        int `json:"PTR,omitempty"`
		Mx         int `json:"MX,omitempty"`
		Txt        int `json:"TXT,omitempty"`
		Afsdb      int `json:"AFSDB,omitempty"`
		Aaaa       int `json:"AAAA,omitempty"`
		Srv        int `json:"SRV,omitempty"`
		Naptr      int `json:"NAPTR,omitempty"`
		Dname      int `json:"DNAME,omitempty"`
		Ds         int `json:"DS,omitempty"`
		Rrsig      int `json:"RRSIG,omitempty"`
		Dnskey     int `json:"DNSKEY,omitempty"`
		Nsec3param int `json:"NSEC3PARAM,omitempty"`
		Tlsa       int `json:"TLSA,omitempty"`
		Cds        int `json:"CDS,omitempty"`
		Cdnskey    int `json:"CDNSKEY,omitempty"`
		Zonemd     int `json:"ZONEMD,omitempty"`
		Svcb       int `json:"SVCB,omitempty"`
		Https      int `json:"HTTPS,omitempty"`
		Spf        int `json:"SPF,omitempty"`
		Any        int `json:"ANY,omitempty"`
	} `json:"qtypes"`
}

func readXmlStats(statsData []byte) error {
	fmt.Printf("Read %d bytes of XML\n", len(statsData))

	var xmlStats bindXmlStats

	// Parse the XML statistics
	err := xml.Unmarshal(statsData, &xmlStats)
	if err != nil {
		fmt.Printf("Error parsing XML: %s\n", err)
		return err
	}

	return nil
}

func readJsonStats(statsData []byte) error {
	// Read the JSON statistics
	var jsonStats bindJsonStats

	err := json.Unmarshal(statsData, &jsonStats)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return err
	}
	return nil
}

func OutputMetricsGraphite() {
	// Output metrics in Graphite format
	for _, metric := range plugin.returnMetrics {
		outTags := "bind.dns"
		for i := len(metric.Tags) - 1; i >= 0; i-- {
			if len(metric.Tags[i]) == 2 {
				outTags = fmt.Sprintf("%s.%s_%s", outTags, metric.Tags[i][0], metric.Tags[i][1])
			} else {
				outTags = fmt.Sprintf("%s.%s", outTags, metric.Tags[i][0])
			}
		}

		fmt.Printf("%s.%s %d %d\n", outTags, metric.Name, metric.Value, metric.Timestamp.Unix())
	}
}
