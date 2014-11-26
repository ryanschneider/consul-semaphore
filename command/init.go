package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/lock"
)

type InitCommand struct {
	Ui   cli.Ui
	Name string
}

func (c *InitCommand) Run(args []string) int {
	var max int
	common, client, err := parseFlags(c.Name, c.Ui, args, func(f *flag.FlagSet) {
		f.IntVar(&max, "max", 1, "maximum concurrent")
	})
	if err != nil {
		return 1
	}

	l, err := lock.New(common.Semaphore, common.Holder, client)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating semaphore: %s", err))
		return 1
	}

	_, _, err = l.SetMax(max)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error setting maximum for semaphore: %s", err))
		return 1
	}
	return 0
}

func (c *InitCommand) Synopsis() string {
	return "Initializes a unowned semaphore in Consul"
}

func (c *InitCommand) Help() string {
	helpText := `
Usage consul-semaphore init [options]

	Initializes a unowned semaphore in Consul.

Options:

	-max                       Maximum concurrent, default 1
%s
	`

	return strings.TrimSpace(fmt.Sprintf(helpText, commonHelp()))
}
