package language

import (
	"testing"

	"github.com/hyperproperties/sopher/pkg/quick"
	"github.com/stretchr/testify/assert"
)

func TestTokenClassString(t *testing.T) {
	tests := []struct {
		description string
		class       TokenClass
		expected    string
	}{
		{
			description: "region",
			class:       RegionToken,
			expected:    "region",
		},
		{
			description: "forall",
			class:       ForallToken,
			expected:    "forall",
		},
		{
			description: "exists",
			class:       ExistsToken,
			expected:    "exists",
		},
		{
			description: "assume",
			class:       AssumeToken,
			expected:    "assume",
		},
		{
			description: "guarantee",
			class:       GuaranteeToken,
			expected:    "guarantee",
		},
		{
			description: "identifier",
			class:       IdentifierToken,
			expected:    "identifier",
		},
		{
			description: "probability",
			class:       ProbabilityToken,
			expected:    "probability",
		},
		{
			description: "expression",
			class:       GoExpressionToken,
			expected:    "expression",
		},
		{
			description: ":",
			class:       ScopeDelimiterToken,
			expected:    ":",
		},
		{
			description: ";",
			class:       GoExpressionDelimiterToken,
			expected:    ";",
		},
		{
			description: "(",
			class:       LeftParenthesisToken,
			expected:    "(",
		},
		{
			description: ")",
			class:       RightParenthesisToken,
			expected:    ")",
		},
		{
			description: "&&",
			class:       LogicalConjunctionToken,
			expected:    "&&",
		},
		{
			description: "||",
			class:       LogicalDisjunctionToken,
			expected:    "||",
		},
		{
			description: "->",
			class:       LogicalImplicationToken,
			expected:    "->",
		},
		{
			description: "<->",
			class:       LogicalBiimplicationToken,
			expected:    "<->",
		},
		{
			description: "!",
			class:       LogicalNegationToken,
			expected:    "!",
		},
		{
			description: ")",
			class:       RightParenthesisToken,
			expected:    ")",
		},
		{
			description: "end",
			class:       EndToken,
			expected:    "end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actual := tt.class.String()
			assert.Equal(t, tt.expected, actual, tt.description)
		})
	}
}

func TestNewTokenProperty(t *testing.T) {
	for counter := 0; counter < 10000; counter++ {
		class := quick.New[TokenClass]()
		lexeme := "the quick brown fox jumps over the lazy dog"

		if class < RegionToken || class > EndToken {
			assert.Panics(t, func() {
				NewToken(class, lexeme)
			})
		} else {
			token := NewToken(class, lexeme)
			assert.Equal(t, class, token.class)
			assert.Equal(t, lexeme, token.lexeme)
		}
	}
}
