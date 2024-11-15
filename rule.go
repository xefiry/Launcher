package main

import (
	"errors"
	"log"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Rule struct {
	Match       string
	Description string
	Exe         string
	Args        []string
	LastUse     time.Time
}

func (r *Rule) Execute() {
	r.LastUse = time.Now()

	cmd := exec.Command(r.Exe, r.Args...)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func (r *Rule) Check() error {
	if len(r.Match) == 0 {
		return errors.New("invalid rule, Match field is empty")
	}
	if len(r.Description) == 0 {
		return errors.New("invalid rule, Description field is empty")
	}
	if len(r.Exe) == 0 {
		return errors.New("invalid rule, Exe field is empty")
	}

	return nil
}

func FilterRules(rules []*Rule, input string) []*Rule {
	var result []*Rule

	lower_input := strings.ToLower(input)
	var lower_match string

	for _, rule := range rules {
		lower_match = strings.ToLower(rule.Match)

		// Check if lower_match starts with lower_input
		// both strings are lowered to ignore case
		// if input is an empty string, it will always match
		if strings.HasPrefix(lower_match, lower_input) {
			result = append(result, rule)
		}
	}

	return result
}

func SortRules(rules []*Rule) {

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].LastUse.After(rules[j].LastUse)
	})
}

func RulesToAray(rules []*Rule) []string {
	var out []string

	for _, rule := range rules {
		out = append(out, rule.Match)
	}

	return out
}
