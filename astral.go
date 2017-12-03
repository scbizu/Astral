package main

import (
	"log"

	"github.com/scbizu/Astral/cmd"
)

func main() {
	if err := cmd.LaunchCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
