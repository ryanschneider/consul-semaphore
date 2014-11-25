package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/command"
)

// Commands is the mapping of all the available Consul commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{

		"init": func() (cli.Command, error) {
			return &command.InitCommand{
				Ui: ui,
			}, nil
		},

		"acquire": func() (cli.Command, error) {
			return &command.AcquireCommand{
				Ui: ui,
			}, nil
		},

		/*
			"release": func() (cli.Command, error) {
				return &command.ReleaseCommand{}, nil
			},

			"exec": func() (cli.Command, error) {
				return &command.ExecCommand{}, nil
			},
		*/

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
