package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
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

// This function is to get data do display in the UI.
// The given rule is split using the input in order to check
// what part of the rule has been matched with the input.
// The result list is of variable size and contains a list of strings
// where any odd index contains a string that matched with input
// and any even index a string that did not.
// Example : ["match", "not match", "match"]
func (r *Rule) GetDisplayStrings(input string) []string {
	result := []string{}
	var tmp string

	// if the input is empty, return ["", all_the_text]
	if input == "" {
		result = append(result, "")
		result = append(result, fmt.Sprintf("%v - %v", r.Match, r.Description))
		return result
	}

	quoted_input := regexp.QuoteMeta(input)
	rule_match := fmt.Sprintf("^(?i)(%v)(.*)", quoted_input)
	rule_desc := fmt.Sprintf("^(?i)(.*)(%v)(.*)", quoted_input)

	re_match := regexp.MustCompile(rule_match)
	re_desc := regexp.MustCompile(rule_desc)

	res_match := re_match.FindStringSubmatch(r.Match)
	res_desc := re_desc.FindStringSubmatch(r.Description)

	if len(res_match) == 0 {
		// if there is no match in the Match part of the rule

		// insert an empty string for the match
		result = append(result, "")
		// and put all the rule match in a tmp string
		tmp = r.Match

	} else if len(res_match) == 3 {
		// if there is a match, there should be 3 elements

		// insert the part that matched (index 0 is the whole rule)
		result = append(result, res_match[1])
		// and put the part that did not match in the tmp string
		tmp = res_match[2]

	} else {
		// this should not happen ... panic
		log.Panicf("error while parsing Match of rule %v", r)
	}

	// Add the separator to the tmp string
	tmp += " - "

	if len(res_desc) == 0 {
		// if there is no match in the Description part of the rule

		// insert an the whole description with the tmp string
		result = append(result, tmp+r.Description)

	} else if len(res_desc) == 4 {
		// if there is a match, there should be 4 elements

		// insert the first element (not matched) with the tmp string
		result = append(result, tmp+res_desc[1])

		// then insert the part that matched
		result = append(result, res_desc[2])

		// finaly the remaining part that did not match (only if it is not empty)
		if len(res_desc[3]) > 0 {
			result = append(result, res_desc[3])
		}

	} else {
		// this should not happen ... panic
		log.Panicf("error while parsing Match of rule %v", r)
	}

	return result
}

func FilterRules(rules []*Rule, input string) []*Rule {
	var result []*Rule

	lower_input := strings.ToLower(input)
	var lower_match string
	var lower_desc string

	for _, rule := range rules {
		lower_match = strings.ToLower(rule.Match)
		lower_desc = strings.ToLower(rule.Description)

		// Check if lower_match starts with lower_input
		// both strings are lowered to ignore case
		// if input is an empty string, it will always match
		if strings.HasPrefix(lower_match, lower_input) || strings.Contains(lower_desc, lower_input) {
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
