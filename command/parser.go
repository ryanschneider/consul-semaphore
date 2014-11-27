package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Parser struct {
	Consul  string
	Path    string
	Holder  string
	Verbose bool

	flags *flag.FlagSet
}

func newParser(name string, args []string, addFlags func(*flag.FlagSet)) (parser *Parser, err error) {
	parser = new(Parser)

	// Parse the flags and options
	parser.flags = flag.NewFlagSet(name, flag.ContinueOnError)

	parser.flags.StringVar(&parser.Consul, "consul", "127.0.0.1:8500",
		"address of the Consul instance")
	parser.flags.StringVar(&parser.Path, "path", "global/semaphore",
		"KV path to the semaphore to use")
	parser.flags.StringVar(&parser.Holder, "holder", "",
		"the holder of the semaphore (default hostname)")
	parser.flags.BoolVar(&parser.Verbose, "verbose", false, "enables verbose output")

	//call setupFunc if supplied
	if addFlags != nil {
		addFlags(parser.flags)
	}

	// If there was a parser error, stop
	if len(args) > 0 {
		if err := parser.flags.Parse(args); err != nil {
			return nil, errors.New(fmt.Sprintf("Error parsing flags: %s", err))
		}
	}

	if parser.Holder == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error determining hostname for default holder: %s", err))
		}

		parser.Holder = hostname
	}

	return
}

func commonHelp() string {
	helpText := `
	-path                      KV path to the semaphore to use
	-holder                    The name of the holder (defaults to hostname)
	-consul                    Consul server to use, defaults to localhost:8500
	-verbose                   Enables verbose output
	`

	return helpText[1 : len(helpText)-1]
}
