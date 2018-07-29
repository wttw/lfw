package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)

// EnvMap returns the docker environment variables
func EnvMap(c *cli.Context) (map[string]string, error) {
	parseRe := regexp.MustCompile(`^export ([^=]+)\s*=\s*(.*)`)
	dmBinary := filepath.Join(c.GlobalString("dockerdir"), "docker-machine")
	cmd := exec.Command(dmBinary, "env", c.GlobalString("dockername"))
	output, err := cmd.Output()
	if err != nil {
		return map[string]string{}, err
	}
	ret := map[string]string{}
	reader := bufio.NewReader(bytes.NewReader(output))
	for {
		line, rerr := reader.ReadString('\n')
		if rerr != nil && rerr != io.EOF {
			return map[string]string{}, rerr
		}
		matches := parseRe.FindStringSubmatch(line)
		if matches != nil {
			ret[matches[1]] = strings.Trim(strings.TrimSpace(matches[2]), `"`)
		}
		if rerr != nil {
			return ret, nil
		}
	}
}

// Env returns the docker environment variables
func Env(c *cli.Context) ([]string, error) {
	e, err := EnvMap(c)
	if err != nil {
		return []string{}, err
	}
	ret := []string{}
	for k, v := range e {
		ret = append(ret, fmt.Sprintf("%s=%s", k, v))
	}
	return ret, nil
}

// CmdEnv implements the command "env" - show docker environment
func CmdEnv(c *cli.Context) error {
	env, err := EnvMap(c)
	if err != nil {
		return err
	}
	for k, v := range env {
		fmt.Printf("export %s=\"%s\"\n", k, v)
	}
	return nil
}
