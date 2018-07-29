package main

import (
	"fmt"

	"github.com/urfave/cli"
)

// CmdId implements the "id" command - show container id for a server
func CmdId(c *cli.Context) error {
	conf, err := CurrentConfig(c)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", conf.Container)
	return nil
}
