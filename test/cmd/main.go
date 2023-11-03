package main

import (
	"log"

	"github.com/khulnasoft-lab/go-toolkit/test/cmd/cli"
)

func main() {
	c := cli.New()

	if err := c.Execute(); err != nil {
		log.Fatalf("error during command execution: %v", err)
	}
}
