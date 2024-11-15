package main

import (
	"log"
)

const (
	APP_TITLE   = "The Launcher"
	APP_VERSION = "v0.1"
	CONFIG_FILE = "config.toml"
)

func main() {
	// read config from file
	config, err := NewConfig(CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	// make sure to write it back at the end
	defer config.Write(CONFIG_FILE)

	GUI_Start(config.Rules)
}
