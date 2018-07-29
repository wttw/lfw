package main

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli"
)

// CmdInfo prints the configuration for a server
func CmdInfo(c *cli.Context) error {
	conf, err := CurrentConfig(c)
	if err != nil {
		return err
	}
	json, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(json))
	return nil
}
