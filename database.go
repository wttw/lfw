package main

import (
	"fmt"

	"github.com/urfave/cli"
)

// CmdDburi prints the external database URI
func CmdDburi(c *cli.Context) error {
	conf, err := CurrentConfig(c)
	if err != nil {
		return err
	}
	ip, err := Ip(c)
	if err != nil {
		return err
	}
	fmt.Printf("mysql://%s:%s@%s:%d/%s\n", conf.Mysql.User, conf.Mysql.Password, ip, conf.Ports.MYSQL, conf.Mysql.Database)
	return nil
}
