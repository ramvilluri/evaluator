package parser

import "testing"

func TestEvaluate(t *testing.T) {

	tt := []struct {
		name       string
		expression string
		result     bool
	}{
		{"exp1", "(a&c)&(b&e)", false},
		{"exp2", "(a&c)&(b|e)", true},
		{"exp3", "(a&c)|(b&e)", true},
		{"exp4", "(a|b)", true},
		{"exp5", "(a&c)", true},
		{"exp6", "a&b", false},
		{"exp7", "(a-c)", true},
		{"exp8", "(a&c) & (b|e) & (a&b)", false},
		{"exp9", "(a&c) & (b|e) & (a|b)", true},
		{"exp10", "(a&b) & (b|e) & (a|b)", false},
		{"exp11", "b", false},
		{"exp12", "c", true},
		{"exp13", "a|b&c", true},
		{"exp14", "a&b|c", true},
		{"exp15", "(a&b) & (b|e) & (b|b)", false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := Evaluate(tc.expression)
			if r != tc.result {
				t.Errorf("%v should be %v but got %v", tc.name, tc.result, r)
			}
		})
	}
}
