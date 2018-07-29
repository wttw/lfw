package main

import "syscall"

// ConfigDir is where we can find the sites.json config file
const ConfigDir = "~/Library/Application Support/Local by Flywheel"

// DockerDir is where
const DockerDir = "/Applications/Local by Flywheel.app/Contents/Resources/extraResources/virtual-machine/vendor/docker/osx"

// DockerName is the name of the docker instance
const DockerName = "local-by-flywheel"

// FIXME(steve): break out into Windows + Mac variants

// Exec runs a command, replacing this process
func Exec(argv []string, env []string) error {
	return syscall.Exec(argv[0], argv, env)
}
