package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/lock"
)

type ReleaseCommand struct {
	Ui cli.Ui
}

func (c *ReleaseCommand) Run(args []string) int {
	client, err := getClient()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating client: %s", err))
		return 1
	}

	l, err := lock.New("TODO-semaphore", "TODO-holder", client)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating semaphore: %s", err))
		return 1
	}

	err = l.Unlock()
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

  -verbose                   Enables verbose output
	`

	return strings.TrimSpace(helpText)
}
