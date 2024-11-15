package main

import (
	"testing"

	"github.com/BurntSushi/toml"
)

const VALID_RULES = `
[[Rules]]
  Match = "Match"
  Description = "Description"
  Exe = "Exe"
  Args = ["1", "2", "3"]

[[Rules]]
  Match = "Match"
  Description = "Description"
  Exe = "Exe"
  Args = []

[[Rules]]
  Match = "Match"
  Description = "Description"
  Exe = "Exe"
`

const INVALID_RULES = `
[[Rules]]
  # empty Match
  Match = ""
  Description = "Description"
  Exe = "Exe"
  Args = ["1", "2", "3"]

  [[Rules]]
  # no Match
  Description = "Description"
  Exe = "Exe"
  Args = ["1", "2", "3"]

  [[Rules]]
  # empty Description
  Match = "Match"
  Description = ""
  Exe = "Exe"
  Args = ["1", "2", "3"]

  [[Rules]]
  # no Description
  Match = "Match"
  Exe = "Exe"
  Args = ["1", "2", "3"]

  [[Rules]]
  # empty Exe
  Match = "Match"
  Description = "Description"
  Exe = ""
  Args = ["1", "2", "3"]

  [[Rules]]
  # no Exe
  Match = "Match"
  Description = "Description"
  Args = ["1", "2", "3"]
`

func TestNewConfigValid(t *testing.T) {
	var config Config

	_, err := toml.Decode(VALID_RULES, &config)
	if err != nil {
		t.Error(err)
	}

	for _, rule := range config.Rules {
		if err := rule.Check(); err != nil {
			t.Errorf("Rule %v sould be valid", rule)
		}
	}
}

func TestNewConfigInvalid(t *testing.T) {
	var config Config

	_, err := toml.Decode(INVALID_RULES, &config)
	if err != nil {
		t.Error(err)
	}

	for _, rule := range config.Rules {
		if err := rule.Check(); err == nil {
			t.Errorf("Rule %v sould be invalid", rule)
		}
	}
}
