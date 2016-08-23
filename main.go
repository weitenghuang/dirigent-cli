package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/weitenghuang/dirigent-cli/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Errorln(err)
		os.Exit(-1)
	}
}
