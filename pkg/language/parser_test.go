package language

import (
	"testing"

	"github.com/hyperproperties/sopher/pkg/iterx"
	"github.com/stretchr/testify/assert"
)

func TestParseAssertion(t *testing.T) {
	tests := []struct {
		description string
		tokens      []Token
		assertion  Node
	}{
		{
			description: "true;",
			tokens: []Token{
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(EndToken, ""),
			},
			assertion: NewGoExpression("true"),
		},
		{
			description: "!true;",
			tokens: []Token{
				NewToken(LogicalNegationToken, "!"),
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(EndToken, ""),
			},
			assertion: NewUnaryExpression(
				LogicalNegation, NewGoExpression("true"),
			),
		},
		{
			description: "!!!true;",
			tokens: []Token{
				NewToken(LogicalNegationToken, "!"),
				NewToken(LogicalNegationToken, "!"),
				NewToken(LogicalNegationToken, "!"),
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(EndToken, ""),
			},
			assertion: NewUnaryExpression(
				LogicalNegation,
				NewUnaryExpression(
					LogicalNegation, 
					NewUnaryExpression(
						LogicalNegation,
						NewGoExpression("true"),
					),
				),
			),
		},
		{
			description: "true; && !false; || true;",
			tokens: []Token{
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalConjunctionToken, "&&"),
				NewToken(LogicalNegationToken, "!"),
				NewToken(GoExpressionToken, "false"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalDisjunctionToken, "||"),
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(EndToken, ""),
			},
			assertion: NewBinaryExpression(
				NewGoExpression("true"),
				LogicalConjunction,
				NewBinaryExpression(
					NewUnaryExpression(
						LogicalNegation,
						NewGoExpression("false"),
					),
					LogicalDisjunction,
					NewGoExpression("true"),
				),
			),
		},
		{
			description: "(true; && !false;) || true;",
			tokens: []Token{
				NewToken(LeftParenthesisToken, "("),
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalConjunctionToken, "&&"),
				NewToken(LogicalNegationToken, "!"),
				NewToken(GoExpressionToken, "false"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(RightParenthesisToken, ")"),
				NewToken(LogicalDisjunctionToken, "||"),
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(EndToken, ""),
			},
			assertion: NewBinaryExpression(
				NewGroup(
					NewBinaryExpression(
						NewGoExpression("true"),
						LogicalConjunction,
						NewUnaryExpression(
							LogicalNegation,
							NewGoExpression("false"),
						),
					),
				),
				LogicalDisjunction,
				NewGoExpression("true"),
			),
		},
		{
			description: "true; && false; -> false;",
			tokens: []Token{
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalConjunctionToken, "&&"),
				NewToken(GoExpressionToken, "false"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalImplicationToken, "->"),
				NewToken(GoExpressionToken, "false"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(EndToken, ""),
			},
			assertion: NewBinaryExpression(
				NewBinaryExpression(
					NewGoExpression("true"),
					LogicalConjunction,
					NewGoExpression("false"),
				),
				LogicalImplication,
				NewGoExpression("false"),
			),
		},
		{
			description: "true; -> false; <-> true; && true;",
			tokens: []Token{
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalImplicationToken, "->"),
				NewToken(GoExpressionToken, "false"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalBiimplicationToken, "<->"),
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalConjunctionToken, "&&"),
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(EndToken, ""),
			},
			assertion: NewBinaryExpression(
				NewBinaryExpression(
					NewGoExpression("true"),
					LogicalImplication,
					NewGoExpression("false"),
				),
				LogicalBiimplication,
				NewBinaryExpression(
					NewGoExpression("true"),
					LogicalConjunction,
					NewGoExpression("true"),
				),
			),
		},
		{
			description: "forall a. exists b. forall c. P(a, b, c);",
			tokens: []Token{
				NewToken(ForallToken, "forall"),
				NewToken(IdentifierToken, "a"),
				NewToken(ScopeDelimiterToken, "."),
				NewToken(ExistsToken, "exists"),
				NewToken(IdentifierToken, "b"),
				NewToken(ScopeDelimiterToken, "."),
				NewToken(ForallToken, "forall"),
				NewToken(IdentifierToken, "c"),
				NewToken(ScopeDelimiterToken, "."),
				NewToken(GoExpressionToken, "P(a, b, c)"),
				NewToken(GoExpressionDelimiterToken, ";"),
			},
			assertion: NewUniversal(
				[]string{"a"},
				NewExistential(
					[]string{"b"},
					NewUniversal(
						[]string{"c"},
						NewGoExpression("P(a, b, c)"),
					),
				),
			),
		},
		{
			description: "forall a. true; || (exists b. false; && forall c. true;)",
			tokens: []Token{
				NewToken(ForallToken, "forall"),
				NewToken(IdentifierToken, "a"),
				NewToken(ScopeDelimiterToken, "."),
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalDisjunctionToken, "||"),
				NewToken(LeftParenthesisToken, "("),
				NewToken(ExistsToken, "exists"),
				NewToken(IdentifierToken, "b"),
				NewToken(ScopeDelimiterToken, "."),
				NewToken(GoExpressionToken, "false"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(LogicalConjunctionToken, "&&"),
				NewToken(ForallToken, "forall"),
				NewToken(IdentifierToken, "c"),
				NewToken(ScopeDelimiterToken, "."),
				NewToken(GoExpressionToken, "true"),
				NewToken(GoExpressionDelimiterToken, ";"),
				NewToken(RightParenthesisToken, ")"),
			},
			assertion: NewBinaryExpression(
				NewUniversal(
					[]string{"a"},
					NewGoExpression("true"),
				),
				LogicalDisjunction,
				NewGroup(
					NewBinaryExpression(
						NewExistential(
							[]string{"b"},
							NewGoExpression("false"),
						),
						LogicalConjunction,
						NewUniversal(
							[]string{"c"},
							NewGoExpression("true"),
						),
					),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			parser := NewParser(iterx.Forward(tt.tokens))
			expression := parser.assertion()
			assert.Equal(t, tt.assertion, expression)
		})
	}
}
