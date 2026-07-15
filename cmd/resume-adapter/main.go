package main

import (
	"log"
	"os"

	"github.com/SimonFremont-dev/cv-builder/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
