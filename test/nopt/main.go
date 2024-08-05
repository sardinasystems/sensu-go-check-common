package main

import (
	"fmt"
	"os"

	"github.com/sardinasystems/sensu-go-check-common/nagios"
	corev2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	StrThreshold string
	Threshold    nagios.Threshold
	CheckValue   float64
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-go-check-common",
			Short:    "test option parsing",
			Keyspace: "sensu.io/plugins/sensu-go-check-common/config",
		},
	}

	options = []sensu.ConfigOption{
		&nagios.ThresholdConfigOption{
			Option: sensu.PluginConfigOption[string]{
				Path:      "threshold",
				Env:       "TEST_THRESHOLD",
				Argument:  "threshold",
				Shorthand: "t",
				Default:   "10:20",
				Usage:     "Threshold",
				Value:     &plugin.StrThreshold,
			},
			Value: &plugin.Threshold,
		},
		&sensu.PluginConfigOption[float64]{
			Path:      "check_value",
			Env:       "TEST_CHECK_VALUE",
			Argument:  "check-value",
			Shorthand: "v",
			Default:   15.0,
			Usage:     "Value to check with threshold",
			Value:     &plugin.CheckValue,
		},
	}
)

func main() {
	useStdin := false
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error check stdin: %v\n", err)
		panic(err)
	}
	//Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		useStdin = true
	}

	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	if plugin.Threshold.Check(plugin.CheckValue) {
		return sensu.CheckStateCritical, fmt.Errorf("alert generated, '%s' aka '%s' %v ~=~ %v", plugin.StrThreshold, plugin.Threshold.String(), plugin.Threshold, plugin.CheckValue)
	}

	fmt.Printf("no alert, '%s' aka '%s' %v ~=~ %v\n", plugin.StrThreshold, plugin.Threshold.String(), plugin.Threshold, plugin.CheckValue)
	return sensu.CheckStateOK, nil
}
