package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/command"
)

// Commands is the mapping of all the available Consul commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	Commands = map[string]cli.CommandFactory{

		"init": func() (cli.Command, error) {
			return &command.InitCommand{
				Ui:   ui,
				Name: "init",
			}, nil
		},

		"acquire": func() (cli.Command, error) {
			return &command.AcquireCommand{
				Ui:   ui,
				Name: "acquire",
			}, nil
		},

		"release": func() (cli.Command, error) {
			return &command.ReleaseCommand{
				Ui:   ui,
				Name: "release",
			}, nil
		},

		"exec": func() (cli.Command, error) {
			return &command.ExecCommand{
				Ui:   ui,
				Name: "exec",
			}, nil
		},

		"version": func() (cli.Command, error) {
			ver := Version
			rel := VersionPrerelease
			if GitDescribe != "" {
				ver = GitDescribe
			}
			if GitDescribe == "" && rel == "" {
				rel = "dev"
			}

			return &command.VersionCommand{
				Revision:          GitCommit,
				Version:           ver,
				VersionPrerelease: rel,
				Ui:                ui,
			}, nil
		},
	}
}
