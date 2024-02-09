package main

import (
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-plugins-bind-dns-checks",
			Short:    "Sensu check to pull metrics from bind DNS server",
			Keyspace: "sensu.io/plugins/sensu-plugins-bind-dns-checks/config",
		},
	}

	options = []sensu.ConfigOption{}
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
