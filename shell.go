package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

// CmdShell spawns an interactive shell
func CmdShell(c *cli.Context) error {
	conf, err := CurrentConfig(c)
	if err != nil {
		return err
	}
	env, err := Env(c)
	if err != nil {
		return err
	}
	dBinary := filepath.Join(c.GlobalString("dockerdir"), "docker")
	cmd := []string{dBinary, "exec", "-it", conf.Container, "/bin/bash"}
	return Exec(cmd, env)
}

// RunCmd runs a command in a bash shell in the docker container, replacing this process
func RunCmd(c *cli.Context, args []string) error {
	conf, err := CurrentConfig(c)
	if err != nil {
		return err
	}
	env, err := Env(c)
	if err != nil {
		return err
	}

	dBinary := filepath.Join(c.GlobalString("dockerdir"), "docker")
	cmd := []string{dBinary, "exec", "-it", conf.Container, "/bin/bash", "-c", strings.Join(args, " ")}
	return Exec(cmd, env)
}

// CmdCommand runs a command
func CmdCommand(c *cli.Context) error {
	args := c.Args()
	if len(args) < 1 {
		return errors.New("command requires a command to run")
	}
	return RunCmd(c, args)
}

// CmdWp runs wp
func CmdWp(c *cli.Context) error {
	args := append([]string{"/usr/local/bin/wp"}, c.Args()...)
	return RunCmd(c, args)
}

// CmdMysql runs mysql
func CmdMysql(c *cli.Context) error {
	conf, err := CurrentConfig(c)
	if err != nil {
		return err
	}
	args := []string{"/usr/bin/mysql"}
	args = append(args, fmt.Sprintf("--user=%s", conf.Mysql.User))
	if conf.Mysql.Password != "" {
		args = append(args, fmt.Sprintf("--password=%s", conf.Mysql.Password))
	}
	args = append(args, fmt.Sprintf("--database=%s", conf.Mysql.Database))
	args = append(args, c.Args()...)
	return RunCmd(c, args)
}
