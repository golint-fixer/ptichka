package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danil/ptichka"
)

func main() {
	version := flag.Bool("version", false, "display version information and exit")
	pathToConfig := flag.String("config", "", "path to config")

	flag.Parse()

	if *version {
		fmt.Println(ptichka.Version)
		os.Exit(0)
	}

	configs, err := ptichka.LoadConfig(*pathToConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on LoadConfig(%s): %v", *pathToConfig, err)
		os.Exit(1)
	}

	ch := make(chan error)
	for _, config := range configs.Accounts {
		go ptichka.Fly(&config, ch)
	}

	var errors []error
	for range configs.Accounts {
		err = <-ch
		if err != nil {
			errors = append(errors, <-ch)
		}
	}

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
		}
		os.Exit(1)
	}
}
