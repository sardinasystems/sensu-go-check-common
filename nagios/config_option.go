package nagios

import (
	corev2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/spf13/cobra"
)

type ThresholdConfigOption struct {
	// Inner string option
	Option sensu.PluginConfigOption[string]

	// Value is the value to read the configured flag or environment variable into.
	// Pass a pointer to any value in your plugin in order to fill it in with the
	// data from a flag or environment variable. The parsing will be done with
	// a function supplied by viper. See the viper documentation for details on
	// how various data types are parsed.
	Value *Threshold

	strValue string
}

func (p *ThresholdConfigOption) SetupFlag(cmd *cobra.Command) error {
	p.Option.Value = &p.strValue
	return p.Option.SetupFlag(cmd)
}

func (p *ThresholdConfigOption) SetValue(valueStr string) (err error) {
	// 1. Set inner string value
	err = p.Option.SetValue(valueStr)
	if err != nil {
		return err
	}

	// 2. Parse threshold
	th, err := ParseThreshold(*p.Option.Value)
	if err != nil {
		return err
	}

	*p.Value = th
	return nil
}

func (p *ThresholdConfigOption) SetAnnotationValue(keySpace string, event *corev2.Event) (sensu.SetAnnotationResult, error) {
	return p.Option.SetAnnotationValue(keySpace, event)
}
