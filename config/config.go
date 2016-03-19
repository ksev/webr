package config

import (
	"flag"
	"os"
)

// Config Struct containing arg data
type Config struct {
	Bind              string
	TargetBindArgName string
	TargetPassTrough  []string
}

// Current the current config values
var Current = Config{}

func init() {
	flag.StringVar(&Current.Bind, "bind", ":8080", "Bind proxy server to this host:port")
	flag.StringVar(&Current.TargetBindArgName, "target-bind-arg-name", "-bind", "The the argument name for passing bind spec to target")
	Current.TargetPassTrough = flag.Args()

	help := flag.Bool("h", false, "Show help message")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
}
