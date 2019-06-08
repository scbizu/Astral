package main

import (
	"log"

	"github.com/scbizu/Astral/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.ErrorLevel)
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
