package command

import (
	"flag"
	"fmt"
	"os"

	api "github.com/armon/consul-api"
	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/lock"
)

type CommandHelper struct {
	Consul  string
	Path    string
	Holder  string
	Verbose bool

	flags  *flag.FlagSet
	client *lock.ConsulLockClient
}

func newCommandHelper(name string, ui cli.Ui, args []string, addFlags func(*flag.FlagSet)) (helper *CommandHelper, err error) {
	helper = new(CommandHelper)

	// Parse the flags and options
	helper.flags = flag.NewFlagSet(name, flag.ContinueOnError)

	helper.flags.StringVar(&helper.Consul, "consul", "127.0.0.1:8500",
		"address of the Consul instance")
	helper.flags.StringVar(&helper.Path, "path", "global/semaphore",
		"KV path to the semaphore to use")
	helper.flags.StringVar(&helper.Holder, "holder", "",
		"the holder of the semaphore (default hostname)")
	helper.flags.BoolVar(&helper.Verbose, "verbose", false, "enables verbose output")

	//call setupFunc if supplied
	if addFlags != nil {
		addFlags(helper.flags)
	}

	// If there was a parser error, stop
	if len(args) > 0 {
		if err := helper.flags.Parse(args); err != nil {
			ui.Error(fmt.Sprintf("Error parsing flags: %s", err))
			return nil, err
		}
	}

	if helper.Holder == "" {
		hostname, err := os.Hostname()
		if err != nil {
			ui.Error(fmt.Sprintf("Error determining hostname: %s", err))
			return nil, err
		}

		helper.Holder = hostname
	}

	err = helper.getClient()
	if err != nil {
		return nil, err
	}

	return
}

func (c *CommandHelper) getClient() (err error) {
	apiClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	c.client, err = lock.NewConsulLockClient(apiClient)
	return
}

func commonHelp() string {
	helpText := `
	-path                      KV path to the semaphore to use
	-holder                    The name of the holder (defaults to hostname)
	-consul                    Consul server to use, defaults to localhost
	-verbose                   Enables verbose output
	`

	return helpText[1 : len(helpText)-1]
}
