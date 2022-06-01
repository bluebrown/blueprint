package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

func setUsage() {
	flag.CommandLine.SortFlags = false

	flag.Usage = func() {
		fmt.Print(`Usage:
  blueprint [source] [destination] [flags]

Flags:
      --set stringArray      set one or more values separated by comma. i.e. a=1,b=2. Flag can be specified multiple times
  -f, --values stringArray   set values from a file or url. Flag can be specified multiple times
      --no-hooks             disable the pre and post hooks
  -h, --help                 show the help text
  -v, --version              show the version
`,
		)
	}
}

func parseFlags() (sets *[]string, values *[]string, noHooks bool) {
	setUsage()
	sets = flag.StringArray("set", []string{}, "set one or more values separated by comma. i.e. a=1,b=2. Flag can be specified multiple times")
	values = flag.StringArrayP("values", "f", []string{}, "set values from a file or url. Flag can be specified multiple times")
	flag.BoolVar(&noHooks, "no-hooks", false, "disable the pre and post hooks")

	var showHelp bool
	flag.BoolVarP(&showHelp, "help", "h", false, "show the help text")

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
