package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/semaphore"
)

type InitCommand struct {
	Ui   cli.Ui
	Name string
}

func (c *InitCommand) Run(args []string) int {
	var max int
	parser, err := newParser(c.Name, args, func(f *flag.FlagSet) {
		f.IntVar(&max, "max", -0xdefa, "maximum concurrent")
	})
	if err != nil {
		return 1
	}

	if max < 0 && max != -0xdefa {
		c.Ui.Error(fmt.Sprintf("Max must be a positive integer: %v", max))
		return 1
	}

	sem, err := semaphore.New(parser.Path, parser.Holder)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing semaphore: %s", err))
		return 1
	}

	if max > 0 {
		_, err = sem.SetMax(uint(max))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error setting maximum for semaphore: %s", err))
			return 1
		}
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
