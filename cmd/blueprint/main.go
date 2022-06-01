package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"
)

var (
	version   string = "0.1.0"
	buildDate string = "2022-03-27"

	templatesDir    string = "templates"
	configFileName  string = "blueprint.yaml"
	valuesFileName  string = "values.yaml"
	helpersFileName string = "_helpers.tpl"
)

func main() {
	// get the flag values
	sets, values, noHooks := parseFlags()

	// input is the repo containing a templates dir and values.yaml
	input := flag.Arg(0)

	// output is the directory to output the rendered templates to
	output := flag.Arg(1)

	if input == "" || output == "" {
		fmt.Println("source and destination must be specified")
		flag.Usage()
		os.Exit(1)
	}

	// create a context that cancels when the user presses ctrl-c
	ctx, cancel := context.WithCancel(context.Background())

	// cancel on signal
	go func() {
		defer cancel()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		select {
		case <-ch:
		case <-ctx.Done():
		}
	}()

	err := run(ctx, input, output, *sets, *values, noHooks)
	if err != nil {
		cancel()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
