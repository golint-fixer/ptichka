package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/danil/ptichka/libs/ptichka"
)

func main() {
	version := flag.Bool("version", false, "display version information and exit")
	pathToConfig := flag.String("config", "", "path to config")

	flag.Parse()

	if *version {
		fmt.Println(ptichka.Version)
		os.Exit(0)
	}

	errs := ptichka.Ptichka(
		*pathToConfig,
		"",
		&ptichka.AnacondaFetcher{},
		&ptichka.SMTPSender{})

	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
