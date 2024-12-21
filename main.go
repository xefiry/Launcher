package main

import (
	"log"
	"os"

	"github.com/xefiry/launcher/launcher"
)

const (
	CONFIG_FILE = "config.toml"
	LOG_FILE    = "launcher.log"
)

func main() {
	// Open log file
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_APPEND, 0644)

	// in case of error ... whatever
	if err != nil {
		log.Print("Could not open log file")
	} else { // otherwise, use it for logs
		defer file.Close()
		log.SetOutput(file)
	}

	// read config from file
	config, err := launcher.NewConfig(CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	// make sure to write it back at the end
	defer config.Write(CONFIG_FILE)

	launcher.GUI_Start(config)
}
