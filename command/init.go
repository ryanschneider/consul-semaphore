package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/lock"
)

type InitCommand struct {
	Ui cli.Ui
}

func (c *InitCommand) Run(args []string) int {
	client, err := getClient()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating client: %s", err))
		return 1
	}

	_, err = lock.New("TODO-semaphore", "TODO-holder", client)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating semaphore: %s", err))
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

  -verbose                   Enables verbose output
	`

	return strings.TrimSpace(helpText)
}
