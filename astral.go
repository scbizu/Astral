package main

import (
	"log"

	"github.com/scbizu/Astral/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	if err := cmd.LaunchCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
