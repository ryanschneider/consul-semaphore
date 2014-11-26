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
	Consul    string
	Semaphore string
	Holder    string
	Verbose   bool
}

func parseFlags(name string, ui cli.Ui, args []string, setupFunc func(*flag.FlagSet)) (common *CommonFlags, client *lock.ConsulLockClient, err error) {
	// Parse the flags and options
	common = new(CommonFlags)

	// Parse the flags and options
	flags := flag.NewFlagSet(name, flag.ContinueOnError)

	flags.StringVar(&common.Consul, "consul", "127.0.0.1:8500",
		"address of the Consul instance")
	flags.StringVar(&common.Semaphore, "semaphore", "global/semaphore",
		"KV path to the semaphore to use")
	flags.StringVar(&common.Holder, "holder", "",
		"the holder of the semaphore (default hostname)")
	flags.BoolVar(&common.Verbose, "verbose", false, "enables verbose output")

	//call setupFunc if supplied
	if setupFunc != nil {
		setupFunc(flags)
	}

	// If there was a parser error, stop
	if len(args) > 0 {
		if err := flags.Parse(args); err != nil {
			ui.Error(fmt.Sprintf("Error parsing flags: %s", err))
			return nil, nil, err
		}
	}

	if common.Holder == "" {
		hostname, err := os.Hostname()
		if err != nil {
			ui.Error(fmt.Sprintf("Error determining hostname: %s", err))
			return nil, nil, err
		}

		common.Holder = hostname
	}

	client, err = getClient(common)
	if err != nil {
		return nil, nil, err
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
	-semaphore                 KV path to the semaphore to use
	-holder                    The name of the holder (defaults to hostname)
	-consul                    Consul server to use, defaults to localhost
	-verbose                   Enables verbose output
	`

	return helpText[1 : len(helpText)-1]
}
