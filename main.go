package main

import (
	"fmt"
	"net"
	"os"

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
	}
)

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
	return sensu.CheckStateOK, nil
}
