package main

import (
	"fmt"
	"os"

	"github.com/bluebrown/blueprint/pkg/types"
	flag "github.com/spf13/pflag"
)

func init() {
	flag.Usage = func() {
		fmt.Print(`Usage:
        blueprint [source] [destination] [flags]

Flags:
    -h, --help          show the help help text
        --no-hooks      disable the pre and post hooks
        --set list      set one or more values seperated by comma. i.e. a=1,b=2. Flag can be specified multiple times
    -f, --values list   set values from a file or url. Flag can be specified multiple times
    -v, --version       show the version
`,
		)
	}
}

func parseFlags() (sets types.StringSlice, values types.StringSlice, noHooks bool) {
	sets = types.StringSlice{}
	values = types.StringSlice{}
	flag.Var(&sets, "set", "set one or more values seperated by comma. i.e. a=1,b=2. Flag can be specified multiple times")
	flag.VarP(&values, "values", "f", "set values from a file or url. Flag can be specified multiple times")
	flag.BoolVar(&noHooks, "no-hooks", false, "disable the pre and post hooks")

	var showHelp bool
	flag.BoolVarP(&showHelp, "help", "h", false, "show the help help text")

	var showVersion bool
	flag.BoolVarP(&showVersion, "version", "v", false, "show the version")

	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("blueprint %s %s\n", version, buildDate)
		os.Exit(0)
	}

	return
}
