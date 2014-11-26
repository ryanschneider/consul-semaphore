package command

import (
	"flag"
	"fmt"
	"os"

	api "github.com/armon/consul-api"
	"github.com/mitchellh/cli"
	"github.com/ryanschneider/consul-semaphore/lock"
)

type CommonFlags struct {
	Consul  string
	Path    string
	Holder  string
	Verbose bool

	flags  *flag.FlagSet
	client *lock.ConsulLockClient
}

func parseFlags(name string, ui cli.Ui, args []string, setupFunc func(*flag.FlagSet)) (common *CommonFlags, err error) {
	common = new(CommonFlags)

	// Parse the flags and options
	common.flags = flag.NewFlagSet(name, flag.ContinueOnError)

	common.flags.StringVar(&common.Consul, "consul", "127.0.0.1:8500",
		"address of the Consul instance")
	common.flags.StringVar(&common.Path, "path", "global/semaphore",
		"KV path to the semaphore to use")
	common.flags.StringVar(&common.Holder, "holder", "",
		"the holder of the semaphore (default hostname)")
	common.flags.BoolVar(&common.Verbose, "verbose", false, "enables verbose output")

	//call setupFunc if supplied
	if setupFunc != nil {
		setupFunc(common.flags)
	}

	// If there was a parser error, stop
	if len(args) > 0 {
		if err := common.flags.Parse(args); err != nil {
			ui.Error(fmt.Sprintf("Error parsing flags: %s", err))
			return nil, err
		}
	}

	if common.Holder == "" {
		hostname, err := os.Hostname()
		if err != nil {
			ui.Error(fmt.Sprintf("Error determining hostname: %s", err))
			return nil, err
		}

		common.Holder = hostname
	}

	common.client, err = getClient(common)
	if err != nil {
		return nil, err
	}

	return
}

func getClient(common *CommonFlags) (client *lock.ConsulLockClient, err error) {
	apiClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	client, err = lock.NewConsulLockClient(apiClient)
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
