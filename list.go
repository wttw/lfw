package main

import (
	"fmt"

	"github.com/urfave/cli"
)

// List prints all the sites
func CmdList(c *cli.Context) error {
	configs, err := Configs(c)
	if err != nil {
		return err
	}
	for _, conf := range configs {
		fmt.Printf("%s\n", conf.Domain)
	}
	return nil
}
