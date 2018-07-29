package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

// Ip returns the IP address of the site
func Ip(c *cli.Context) (string, error) {
	ipfile, err := homedir.Expand(filepath.Join(ConfigDir, "machine-ip.json"))
	if err != nil {
		return "", nil
	}
	content, err := ioutil.ReadFile(ipfile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

// CmdIp prints the IP address of the site
func CmdIp(c *cli.Context) error {
	ip, err := Ip(c)
	if err != nil {
		return err
	}
	fmt.Println(ip)
	return nil
}
