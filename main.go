package main

import (
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
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	return sensu.CheckStateOK, nil
}
