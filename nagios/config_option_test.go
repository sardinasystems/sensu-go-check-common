package nagios

import (
	"testing"

	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/stretchr/testify/assert"
)

var commonOpt = ThresholdConfigOption{
	Option: sensu.PluginConfigOption[string]{
		Argument:  "string",
		Default:   "Default1",
		Env:       "ENV_1",
		Path:      "path1",
		Shorthand: "d",
		Usage:     "First argument",
		Secret:    true,
	},
}

func TestSetupFlag(t *testing.T) {

	finalValue := Threshold{}
	expectedValue := Threshold{
		Start:    10.0,
		End:      20.0,
		Inverted: true,
	}

	strValue := ""

	option := commonOpt
	option.Option.Value = &strValue
	option.Value = &finalValue

	err := option.SetValue("@10:20")
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, finalValue)
}

func TestVerifyType(t *testing.T) {
	assert.Implements(t, (*sensu.ConfigOption)(nil), &commonOpt)
}
