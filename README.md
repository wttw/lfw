# lfw
Access a [Local by Flywheel](https://local.getflywheel.com) installation from the commandline

## Installation

Download the file from the [github releases page](https://github.com/wttw/lfw/releases/latest), for your operating system, unzip it and put it somewhere on your path.

To install from source, with a recent go installation, run `go install -i github.com/wttw/lfw`.

Currently it is macOS specific, but a build for Windows would likely be fairly easy.

## Usage

`lfw` takes one of a number of subcommands, listed in the help text. A few (`lfw list`, `lfw env`) display information about the Local by Flywheel installation, but most perform some action on a particular site.

`lfw shell` will open a bash shell in the docker container for a site. You can specify the site using the --site flag, e.g. `lfw --site=example shell` will open a shell in the site called `example.local`.

If you are in a directory of the site (e.g. `~/Local Sites/example`) then lfw will use that site, or if only a single site is running it will use the running site.

```
NAME:
   lfw - Client for Local by Flywheel

USAGE:
   lfw [global options] command [command options] [arguments...]

VERSION:
   0.1.0

DESCRIPTION:
   Client for Local by Flywheel

COMMANDS:
     shell, sh, ssh   open a shell on a site
     wp               run 'wp' on a site
     mysql            run 'mysql' on a site
     command, cmd, c  run a command on a site
     list             list all sites
     env              show docker environment
     id               show container id for a site
     info             show configuration for a site
     ip               show IP address for a site
     dburi            get MySQL connection string for a site
     help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --site value        Site to access [$LFW_SITE]
   --config value      Local by Flywheel configuration directory (default: "~/Library/Application Support/Local by Flywheel") [$LFW_CONFIG]
   --dockerdir value   Directory containing docker binaries (default: "/Applications/Local by Flywheel.app/Contents/Resources/extraResources/virtual-machine/vendor/docker/osx") [$LFW_DOCKER_DIR]
   --dockername value  Docker instance name (default: "local-by-flywheel") [$LFW_DOCKER_NAME]
   --help, -h          show help
   --version, -v       print the version
```