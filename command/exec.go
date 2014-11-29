package command

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/semaphore"
)

type ExecCommand struct {
	Ui   cli.Ui
	Name string
}

func (c *ExecCommand) Run(args []string) (ret int) {
	parser, err := newParser(c.Name, args, func(f *flag.FlagSet) {
	})
	if err != nil {
		return 1
	}

	remainingArgs := parser.flags.Args()
	if len(remainingArgs) == 0 {
		c.Ui.Error("Error: No command to execute given")
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	sem, err := semaphore.New(parser.Path, parser.Holder)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing semaphore: %s", err))
		return 1
	}

	err = sem.Acquire(true)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error acquiring semaphore: %s", err))
		return 1
	}

	//defer Release until after the command is executed
	defer func() {
		err = sem.Release()
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error releasing semaphore: %s", err))
			ret = 1
		}
	}()

	//execute the command
	err = c.execute(remainingArgs...)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error executing command: %s", err))
		return 1
	}

	return 0
}

// execute accepts a command string and runs that command string on the current
// system.
func (c *ExecCommand) execute(args ...string) error {
	// Create and invoke the command
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = WriterFunc(func(p []byte) (n int, err error) {
		c.Ui.Output(fmt.Sprintf("%s", p))
		return len(p), nil
	})
	cmd.Stderr = WriterFunc(func(p []byte) (n int, err error) {
		c.Ui.Error(fmt.Sprintf("%s", p))
		return len(p), nil
	})
	return cmd.Run()
}

func (c *ExecCommand) Synopsis() string {
	return "Executes a command, wrapped in a consul semaphore"
}

func (c *ExecCommand) Help() string {
	helpText := `
Usage consul-semaphore exec [options] <command> [args...]

  Executes command, wrapped in a consul semaphore.

Options:

%s
	--                         Stop parsing args, next arg is command
	`

	return strings.TrimSpace(fmt.Sprintf(helpText, commonHelp()))
}

/// Used to coerce my coroutines into Writers
type WriterFunc func(p []byte) (n int, err error)

func (wf WriterFunc) Write(p []byte) (n int, err error) {
	return wf(p)
}
