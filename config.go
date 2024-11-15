package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Rules []*Rule
}

func NewConfig(filepath string) (*Config, error) {
	var config Config
	var undefined_time time.Time // unset var to get default value of Time

	data, err := toml.DecodeFile(filepath, &config)
	if err != nil {
		return nil, err
	}

	// There should not be any invalid rule in the file
	if undecoded := data.Undecoded(); len(undecoded) != 0 {
		msg := fmt.Sprintf("invalid keys found in config file: %v", undecoded)

		return nil, errors.New(msg)
	}

	// But there has to be rules
	if len(config.Rules) == 0 {
		return nil, errors.New("no rules found in config file")
	}

	// Check if all rules are valid
	valid := true
	for i, rule := range config.Rules {
		err := rule.Check()
		if err != nil {
			valid = false
			log.Printf("Rule n°%v (%v) : %v\n", i, rule, err)
		}
	}
	if !valid {
		return nil, errors.New("invalid rules detected")
	}

	// Loop on all rules. If the time was not defined, set it to epoch time 0
	for _, rule := range config.Rules {
		if rule.LastUse == undefined_time {
			rule.LastUse = time.Unix(0, 0)
		}
	}

	return &config, nil

}

func (config *Config) Write(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	err = toml.NewEncoder(file).Encode(config)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}
