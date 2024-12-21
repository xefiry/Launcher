package launcher

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestRuleGetDisplayStrings(t *testing.T) {
	rules := []*Rule{
		// regular case
		{"Demo rule", "Description", "dummy.exe", nil, time.Unix(0, 0)},

		// edge case (regexp metacharacters)
		{"Edge rule", "{}[]()^$.|*+?", "dummy.exe", nil, time.Unix(0, 0)},
	}

	var tests = []struct {
		name  string
		rule  *Rule
		input string
		want  []string
	}{
		{"empty", rules[0], "", []string{"", "Demo rule - Description"}},
		{"match 1", rules[0], "demo", []string{"Demo", " rule - Description"}},
		{"match 2", rules[0], "Demo ", []string{"Demo ", "rule - Description"}},
		{"match 3", rules[0], "Demo rule", []string{"Demo rule", " - Description"}},
		{"desc 1", rules[0], "Des", []string{"", "Demo rule - ", "Des", "cription"}},
		{"desc 2", rules[0], "crip", []string{"", "Demo rule - Des", "crip", "tion"}},
		{"desc 3", rules[0], "tion", []string{"", "Demo rule - Descrip", "tion"}},
		{"desc 4", rules[0], "Description", []string{"", "Demo rule - ", "Description"}},
		{"both 1", rules[0], "De", []string{"De", "mo rule - ", "De", "scription"}},

		{"edge 0", rules[1], "XXX", []string{"", "Edge rule - {}[]()^$.|*+?"}},
		{"edge 1", rules[1], "{", []string{"", "Edge rule - ", "{", "}[]()^$.|*+?"}},
		{"edge 2", rules[1], "}", []string{"", "Edge rule - {", "}", "[]()^$.|*+?"}},
		{"edge 3", rules[1], "[", []string{"", "Edge rule - {}", "[", "]()^$.|*+?"}},
		{"edge 4", rules[1], "]", []string{"", "Edge rule - {}[", "]", "()^$.|*+?"}},
		{"edge 5", rules[1], "(", []string{"", "Edge rule - {}[]", "(", ")^$.|*+?"}},
		{"edge 6", rules[1], ")", []string{"", "Edge rule - {}[](", ")", "^$.|*+?"}},
		{"edge 7", rules[1], "^", []string{"", "Edge rule - {}[]()", "^", "$.|*+?"}},
		{"edge 8", rules[1], "$", []string{"", "Edge rule - {}[]()^", "$", ".|*+?"}},
		{"edge 9", rules[1], ".", []string{"", "Edge rule - {}[]()^$", ".", "|*+?"}},
		{"edge 10", rules[1], "|", []string{"", "Edge rule - {}[]()^$.", "|", "*+?"}},
		{"edge 11", rules[1], "*", []string{"", "Edge rule - {}[]()^$.|", "*", "+?"}},
		{"edge 12", rules[1], "+", []string{"", "Edge rule - {}[]()^$.|*", "+", "?"}},
		{"edge 13", rules[1], "?", []string{"", "Edge rule - {}[]()^$.|*+", "?"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test the function
			ans := tt.rule.GetDisplayStrings(tt.input, true)

			// compare result length
			if len(ans) != len(tt.want) {
				t.Errorf("got %d, want %d", len(ans), len(tt.want))
			}

			// deep compare (only if expected result is not empty)
			if len(tt.want) != 0 && !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}

			// check that the split does not modify the values
			join_ans := strings.Join(ans, "")
			join_rule := fmt.Sprintf("%v - %v", tt.rule.Match, tt.rule.Description)

			if join_ans != join_rule {
				t.Errorf("got '%v', want '%v'", join_ans, join_rule)
			}
		})
	}
}

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
	rules2 := FilterRules(rules1, "", true)

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
		{"Demo 1", "Description 1", "dummy.exe", nil, time.Unix(0, 0)},
		{"demo 2", "Description 2", "dummy.exe", nil, time.Unix(0, 0)},
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
		{"desc 1", "Desc", []string{"Demo 1", "demo 2"}},
		{"desc 2", "Description 1", []string{"Demo 1"}},
		{"desc fuzy 1", "cription", []string{"Demo 1", "demo 2"}},
		{"desc fuzy 2", "cription 2", []string{"demo 2"}},

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
			ans := FilterRules(rules, tt.input, true)

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
