package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
		fmt.Fprintf(os.Stderr, "Error on load config: %v", err)
		os.Exit(1)
	}
	errCh := make(chan error)
	for _, config := range configs.Accounts {

		var infHandler, errHandler io.Writer

		if len(config.LogFile) > 0 {
			logFile, err := os.OpenFile(
				config.LogFile,
				os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"Error on open log file %s %v",
					config.LogFile, err)
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer func() { _ = logFile.Close() }()
			infHandler, errHandler = logFile, logFile
		} else {
			infHandler = os.Stdout
			errHandler = os.Stderr
		}

		if !config.Verbose {
			infHandler = ioutil.Discard
		}

		l := config
		go ptichka.Fly(
			&l,
			errCh,
			log.New(infHandler, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			log.New(errHandler, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile))
	}

	var errs []error
	for range configs.Accounts {
		err := <-errCh
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
		}
		os.Exit(1)
	}
}
