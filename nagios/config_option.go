package nagios

import (
	corev2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
}

func (p *ThresholdConfigOption) SetupFlag(cmd *cobra.Command) error {
	// 1. Provide receiver if not set
	if p.Option.Value == nil {
		var strValue string
		p.Option.Value = &strValue
	}

	// 2. Process default -- part of p.Option.SetupFlag(cmd)
	err := viper.BindEnv(p.Option.Argument, p.Option.Env)
	if err != nil {
		return err
	}
	viper.SetDefault(p.Option.Argument, p.Option.Default)

	*p.Option.Value = viper.GetString(p.Option.Argument)
	err = p.Value.Set(*p.Option.Value)
	if err != nil {
		return err
	}

	wrap := &Wrapper{
		Threshold: p.Value,
		StrValue:  p.Option.Value,
	}

	cmd.Flags().VarP(wrap, p.Option.Argument, p.Option.Shorthand, p.Option.Usage)

	return err
}

func (p *ThresholdConfigOption) SetValue(valueStr string) (err error) {
	// 1. Set inner string value
	err = p.Option.SetValue(valueStr)
	if err != nil {
		return err
	}

	// 2. Parse threshold
	return p.Value.Set(*p.Option.Value)
}

func (p *ThresholdConfigOption) SetAnnotationValue(keySpace string, event *corev2.Event) (sensu.SetAnnotationResult, error) {
	return p.Option.SetAnnotationValue(keySpace, event)
}

type Wrapper struct {
	*Threshold
	StrValue *string
}

func (w *Wrapper) Set(s string) error {
	*w.StrValue = s
	return w.Threshold.Set(s)
}
