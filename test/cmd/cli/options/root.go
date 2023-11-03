package options

import (
	"github.com/spf13/pflag"

	"github.com/khulnasoft-lab/go-toolkit/log"
)

type CommandConfig struct {
	Option1 bool `json:"option1" yaml:"option1" mapstructure:"option1"`
	Option2 uint `json:"option2" yaml:"option2" mapstructure:"option2"`
}

func NewCommandConfig(log *log.Config) *CommandConfig {
	return &CommandConfig{}
}

func (c *CommandConfig) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&c.Option1, "option-1", "", c.Option1, "option 1")
	flags.UintVarP(&c.Option2, "option-2", "", c.Option2, "option 2")
}
