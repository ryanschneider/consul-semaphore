package command

import (
	"strings"
)

type ReleaseCommand struct {
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
