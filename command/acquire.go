package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/lock"
)

type AcquireCommand struct {
	Ui   cli.Ui
	Name string
}

func (c *AcquireCommand) Run(args []string) int {
	common, client, err := parseFlags(c.Name, c.Ui, args, func(f *flag.FlagSet) {
	})
	if err != nil {
		return 1
	}

	l, err := lock.New(common.Semaphore, common.Holder, client)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating semaphore: %s", err))
		return 1
	}

	err = l.Lock()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error locking semaphore: %s", err))
		return 1
	}

	return 0
}

func (c *AcquireCommand) Synopsis() string {
	return "Acquires a semaphore in Consul"
}

func (c *AcquireCommand) Help() string {
	helpText := `
Usage consul-semaphore acquire [options]

  Acquires a semaphore in Consul.

Options:

  -verbose                   Enables verbose output
	`

	return strings.TrimSpace(helpText)
}
