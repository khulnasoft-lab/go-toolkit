package options

import (
	"github.com/spf13/pflag"
)

type Format struct {
	Output string `json:"output" yaml:"output" mapstructure:"output"`
}

func NewFormat() *Format {
	return &Format{
		Output: "",
	}
}

func (r *Format) AddFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&r.Output, "output", "o", r.Output, "output options")
}
