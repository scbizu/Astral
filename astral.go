package main

import (
	"cmd"
	"log"
)

func main() {
	if err := cmd.LaunchCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
