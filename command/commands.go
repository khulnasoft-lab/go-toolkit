package command

import (
	"github.com/spf13/cobra"

	"github.com/khulnasoft-lab/go-toolkit/inject"
)

func MakeOne(c inject.Container, creatorFunc any) *cobra.Command {
	return Make(c, creatorFunc)[0]
}

func Make(c inject.Container, creatorFuncs ...any) (out []*cobra.Command) {
	for _, fn := range creatorFuncs {
		cmd := inject.MustInvoke[*cobra.Command](c, fn)
		out = append(out, cmd)
	}
	return
}

func Converters(c inject.Container, converterFuncs ...any) {
	for _, fn := range converterFuncs {
		c.Register(inject.Singleton(fn))
	}
}
