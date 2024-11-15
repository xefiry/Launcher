package main

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

// Tests if []Rule obtained through Filter gives a reference of a rule
// Changing the obtained rule should change the original one.
func TestRuleFilterNoCopy(t *testing.T) {
	// Create a rule list
	var rules1 = []*Rule{
		{"xxx", "ChangeMe", "dummy.exe", nil, time.Unix(0, 0)},
	}

	// This function should not interfere
	SortRules(rules1)

	// Get a rule by filtering
	rules2 := FilterRules(rules1, "")

	// Modify a field of the rule in the filtered list
	rules2[0].Description = "Modified"

	// It should change the original list
	if rules1[0].Description != "Modified" {
		t.Errorf("Rule description has not changed in the original Rule list")
	}

	// Check if the modification is applied (just in case)
	if rules2[0].Description != "Modified" {
		t.Errorf("Rule description was not modified")
	}
}

func TestRuleFilter(t *testing.T) {
	var rules = []*Rule{
		{"Demo 1", "Demo rule 1", "dummy.exe", nil, time.Unix(0, 0)},
		{"demo 2", "Demo rule 2", "dummy.exe", nil, time.Unix(0, 0)},
		{"r/(a-z)+", "Sub test 1", "dummy.exe", nil, time.Unix(0, 0)},
	}

	var tests = []struct {
		name  string
		input string
		want  []string
	}{
		{"empty", "", []string{"Demo 1", "demo 2", "r/(a-z)+"}},
		{"static 1", "demo", []string{"Demo 1", "demo 2"}},
		{"static 2", "demo 1", []string{"Demo 1"}},
		{"static 3", "Demo 2", []string{"demo 2"}},
		{"static 4", "emo", []string{}},

		// ToDo: behaviour yet to be defined for regexp
		/*{"regexp basic 1", "r", []string{"r/(a-z)+"}},
		{"regexp basic 2", "r/", []string{"r/(a-z)+"}},
		{"regexp basic 3", "r/3", []string{}},
		{"regexp basic 4", "r/(", []string{}},
		{"regexp single 1", "r/sub", []string{"r/(a-z)+"}},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test the function
			ans := FilterRules(rules, tt.input)

			// compare result length
			if len(ans) != len(tt.want) {
				t.Errorf("got %d, want %d", len(ans), len(tt.want))
			}

			// sort results and deep compare (only if expected result is not empty)
			sort.Strings(tt.want)
			ans2 := RulesToAray(ans)
			if len(tt.want) != 0 && !reflect.DeepEqual(ans2, tt.want) {
				t.Errorf("got %v, want %v", ans2, tt.want)
			}
		})
	}
}
