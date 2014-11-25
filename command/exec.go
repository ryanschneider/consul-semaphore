package command

import (
	"strings"
)

type ExecCommand struct {
}

func (c *ExecCommand) Help() string {
	helpText := `
Usage consul-semaphore exec [options]

  Executes a command, wrapped in a consul semaphore.

Options:

  -verbose                   Enables verbose output
	`

	return strings.TrimSpace(helpText)
}
