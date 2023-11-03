package convert

import (
	"github.com/khulnasoft-lab/go-toolkit/test/cmd/app"
	"github.com/khulnasoft-lab/go-toolkit/test/cmd/cli/options"
)

func CatalogerConfig(c *options.CommandConfig) app.CatalogerConfig {
	cfg := app.DefaultCatalogerConfig()
	return cfg
}
