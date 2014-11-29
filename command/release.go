package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/semaphore"
)

type ReleaseCommand struct {
	Ui   cli.Ui
	Name string
}

func (c *ReleaseCommand) Run(args []string) int {
	parser, err := newParser(c.Name, args, nil)
	if err != nil {
		return 1
	}

	sem, err := semaphore.New(parser.Path, parser.Holder)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing semaphore: %s", err))
		return 1
	}

	err = sem.Release()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error releasing semaphore: %s", err))
		return 1
	}

	return 0
}

func (c *ReleaseCommand) Synopsis() string {
	return "Releases a previously acquired semaphore in consul"
}

func (c *ReleaseCommand) Help() string {
	helpText := `
Usage consul-semaphore release [options]

  Releases a previously acquired semaphore in consul.

Options:

%s
	`

	return strings.TrimSpace(fmt.Sprintf(helpText, commonHelp()))
}
