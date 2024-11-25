package language

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyPrint(t *testing.T) {
	tests := []struct {
		description string
		source      string
		print       string
	}{
		{
			description: "Guarantee with single quantifier",
			source:      "guarantee: forall e. e >= 0",
			print:       "region: guarantee: forall e. e >= 0;",
		},
		{
			description: "Guarantee with constant grouped expression",
			source:      "guarantee: (e >= 0  ;  )",
			print:       "region: guarantee: (e >= 0  ;)",
		},
		{
			description: "Guarantee with constant grouped quantifier",
			source:      "guarantee: ((forall e. e >= 0;))",
			print:       "region: guarantee: ((forall e. e >= 0;))",
		},
		{
			description: "Assume with nested quantifiers",
			source:      "assume: forall e0 e1. exists e2. true",
			print:       "region: assume: forall e0 e1. exists e2. true;",
		},
		{
			description: "Region with nested quantifiers and checking generalised non-interference",
			source:      "guarantee: forall e0 e1. exists e2. e2.high == e0.high && e2.low == e1.low",
			print:       "region: guarantee: forall e0 e1. exists e2. e2.high == e0.high && e2.low == e1.low;",
		},
		{
			description: "Two named regions with a single quantifier",
			source:      "region Positive: guarantee: forall e. e >= 0; region Negative: guarantee: forall e. e < 0",
			print:       "region Positive: guarantee: forall e. e >= 0; region Negative: guarantee: forall e. e < 0;",
		},
		{
			description: "nested quantifiers in a guarantee",
			source:      "guarantee: forall a. true; && exists b. false",
			print:       "region: guarantee: forall a. true; && exists b. false;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			parser := NewParser(LexString(tt.source))
			node := parser.Parse()
			print := strings.Trim(Print(node), " ")
			assert.Equal(t, tt.print, print)
		})
	}
}
