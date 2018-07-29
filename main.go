package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/gommon/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

// Mysql holds the mysql connection information for a site
type Mysql struct {
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Ports holds the ports each service runs on
type Ports struct {
	HTTP        int `json:"HTTP"`
	HTTPS       int `json:"HTTPS"`
	MYSQL       int `json:"MYSQL"`
	MAILCATCHER int `json:"MAILCATCHER"`
}

// Config holds the configuration for a single site
type Config struct {
	Blueprint          string `json:"blueprint"`
	Environment        string `json:"environment"`
	PhpVersion         string `json:"phpVersion"`
	MysqlVersion       string `json:"mysqlVersion"`
	WebServer          string `json:"webServer"`
	DevMode            bool   `json:"devMode"`
	AdminUsername      string `json:"adminUsername"`
	AdminPassword      string `json:"adminPassword"`
	AdminEmail         string `json:"adminEmail"`
	MultiSite          string `json:"multiSite"`
	Name               string `json:"name"`
	Path               string `json:"path"`
	Domain             string `json:"domain"`
	ID                 string `json:"id"`
	LocalVersion       string `json:"localVersion"`
	Container          string `json:"container"`
	Mysql              Mysql  `json:"mysql"`
	Ports              Ports  `json:"ports"`
	EnvironmentVersion string `json:"environmentVersion"`
	SslSHA1            string `json:"sslSHA1"`
}

func main() {
	app := cli.NewApp()

	app.Name = "lfw"
	app.Version = "0.1.0"
	app.Description = "Client for Local by Flywheel"
	app.Usage = "Client for Local by Flywheel"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "site",
			Usage:  "Site to access",
			EnvVar: "LFW_SITE",
		},
		cli.StringFlag{
			Name:   "config",
			Value:  ConfigDir,
			Usage:  "Local by Flywheel configuration directory",
			EnvVar: "LFW_CONFIG",
		},
		cli.StringFlag{
			Name:   "dockerdir",
			Value:  DockerDir,
			Usage:  "Directory containing docker binaries",
			EnvVar: "LFW_DOCKER_DIR",
		},
		cli.StringFlag{
			Name:   "dockername",
			Value:  DockerName,
			Usage:  "Docker instance name",
			EnvVar: "LFW_DOCKER_NAME",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "shell",
			Aliases: []string{"sh", "ssh"},
			Usage:   "open a shell on a site",
			Action:  CmdShell,
		},
		{
			Name:   "wp",
			Usage:  "run 'wp' on a site",
			Action: CmdWp,
		},
		{
			Name:   "mysql",
			Usage:  "run 'mysql' on a site",
			Action: CmdMysql,
		},
		{
			Name:            "command",
			Aliases:         []string{"cmd", "c"},
			Usage:           "run a command on a site",
			SkipFlagParsing: true,
			Action:          CmdCommand,
		},
		{
			Name:   "list",
			Usage:  "list all sites",
			Action: CmdList,
		},
		{
			Name:   "env",
			Usage:  "show docker environment",
			Action: CmdEnv,
		},
		{
			Name:   "id",
			Usage:  "show container id for a site",
			Action: CmdId,
		},
		{
			Name:   "info",
			Usage:  "show configuration for a site",
			Action: CmdInfo,
		},
		{
			Name:   "ip",
			Usage:  "show IP address for a site",
			Action: CmdIp,
		},
		{
			Name:   "dburi",
			Usage:  "get MySQL connection string for a site",
			Action: CmdDburi,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		// log.Fatal(err)
		os.Exit(1)
	}
}

var cachedConfig map[string]Config

// Configs loads the configuration for all sites
func Configs(c *cli.Context) (map[string]Config, error) {
	if cachedConfig != nil {
		return cachedConfig, nil
	}
	var err error
	ret := map[string]Config{}
	configFile := filepath.Join(c.GlobalString("config"), "sites.json")
	if strings.HasPrefix(configFile, "~") {

		configFile, err = homedir.Expand(configFile)
		if err != nil {
			return map[string]Config{}, err
		}
	}
	cf, err := os.Open(configFile)
	if err != nil {
		return map[string]Config{}, err
	}
	decoder := json.NewDecoder(cf)
	err = decoder.Decode(&ret)
	if err != nil {
		return map[string]Config{}, err
	}
	cachedConfig = ret
	return cachedConfig, nil
}

// ConfigByName returns the configuration for a server
func ConfigByName(c *cli.Context, name string) (Config, error) {
	configs, err := Configs(c)
	if err != nil {
		return Config{}, err
	}

	for _, conf := range configs {
		if conf.Domain == name {
			return conf, nil
		}
	}

	for _, conf := range configs {
		if conf.Name == name {
			return conf, nil
		}
	}

	for _, conf := range configs {
		if conf.Domain == name+".local" {
			return conf, nil
		}
	}

	conf, ok := configs[name]
	if ok {
		return conf, nil
	}
	return Config{}, fmt.Errorf("no such server as '%s'", name)
}

// CurrentSite returns the current site, either explicitly or heuristically
func CurrentSite(c *cli.Context) (string, error) {
	// Explicitly set, as flag or environment variable
	site := c.GlobalString("site")
	if site != "" {
		return site, nil
	}

	confs, err := Configs(c)
	if err != nil {
		return "", err
	}

	// Heuristics time - check to see if we're in an app directory
	wd, err := os.Getwd()
	if err != nil {
		log.Errorf("Couldn't get working directory: %s", err.Error())
	} else {

		for _, conf := range confs {
			installdir, err := homedir.Expand(conf.Path)
			if err == nil {
				if strings.HasPrefix(wd, installdir) {
					return conf.Domain, nil
				}
			}
		}
	}

	// Maybe there's just a single site running
	var statuses map[string]string
	statusesFile := filepath.Join(c.GlobalString("config"), "site-statuses.json")
	if strings.HasPrefix(statusesFile, "~") {

		statusesFile, err = homedir.Expand(statusesFile)
		if err != nil {
			return "", err
		}
	}
	cf, err := os.Open(statusesFile)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(cf)
	err = decoder.Decode(&statuses)
	if err != nil {
		return "", err
	}

	running := []string{}
	for k, v := range statuses {
		if v == "running" {
			running = append(running, k)
		}
	}

	if len(running) != 1 {
		list := ""
		if len(running) <= 3 {
			names := []string{}
			for _, tag := range running {
				cf, ok := confs[tag]
				if ok {
					names = append(names, strings.TrimSuffix(cf.Domain, ".local"))
				}
			}
			list = fmt.Sprintf(" (%s)", strings.Join(names, ", "))
		}
		return "", fmt.Errorf("I can't guess which site to access - try with --site or $LFW_SITE%s", list)
	}
	conf, ok := confs[running[0]]
	if !ok {
		return "", fmt.Errorf("Can't find supposedly running site '%s'", running[0])
	}

	// TODO(steve): dotfiles in a project directory?

	return conf.Domain, nil
}

// CurrentConfig returns the config of the current site
func CurrentConfig(c *cli.Context) (Config, error) {
	name, err := CurrentSite(c)
	if err != nil {
		return Config{}, err
	}
	return ConfigByName(c, name)
}
