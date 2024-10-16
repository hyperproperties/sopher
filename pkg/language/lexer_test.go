package language

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
	"unicode"

	"github.com/hyperproperties/sopher/pkg/iterx"
	"github.com/stretchr/testify/assert"
)

func TestConsumeWord(t *testing.T) {
	prefix := "assume"
	skip := unicode.IsSpace
	body := func(str string) bool {
		return str != ":"
	}
	suffix := ":"

	str := "assume:"
	lexer := NewLexer(iterx.FromSlice([]rune(str)))

	found, lookahead, _ := lexer.peekWord(prefix, skip, body, suffix)
	assert.True(t, found)
	assert.Equal(t, len(str), lookahead)

	lexer.consumeWord(prefix, skip, body, suffix)
	_, exists := lexer.next()
	assert.False(t, exists)
}

func TestPeekWord(t *testing.T) {
	tests := []struct {
		description string
		word        string
		found       bool
		lookahead   int
	}{
		{
			description: "",
			word:        "ass",
			found:       false,
			lookahead:   3,
		},
		{
			description: "",
			word:        "assume:",
			found:       true,
			lookahead:   7,
		},
		{
			description: "",
			word:        "assume      :",
			found:       true,
			lookahead:   13,
		},
		{
			description: "",
			word:        "assume a b a :",
			found:       true,
			lookahead:   14,
		},
		{
			description: "",
			word:        "assume a b c d e",
			found:       false,
			lookahead:   16,
		},
	}

	prefix := "assume"
	skip := unicode.IsSpace
	body := func(str string) bool {
		return str != ":"
	}
	postfix := ":"

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			lexer := NewLexer(iterx.FromSlice([]rune(tt.word)))
			found, lookahead, _ := lexer.peekWord(prefix, skip, body, postfix)
			assert.Equal(t, tt.found, found)
			assert.Equal(t, tt.lookahead, lookahead)
		})
	}
}

func TestLexerPeekAdvance(t *testing.T) {
	str := "hello, world!"
	lexer := NewLexer(iterx.FromSlice([]rune(str)))

	var builder strings.Builder
	for {
		pCharacter, pOk := lexer.peek(1)
		aCharacter, aOk := lexer.next()

		assert.Equal(t, pCharacter, aCharacter)
		assert.Equal(t, pOk, aOk)

		if !pOk || !aOk {
			break
		}

		builder.WriteRune(aCharacter)
	}

	assert.Equal(t, builder.String(), str)
}

func TestLexMultiLine(t *testing.T) {
	source := `
	region  Positive
assume: forall e0.in >= 0
			guarantee: forall e0.
				ret0 >= 0
	region Negative
		assume
			forall e0. in <= 0
		guarantee
			exists e0
				ret0 <= 0
		`
	expected := []TokenClass{
		RegionToken, IdentifierToken, ScopeDelimiterToken,
		AssumeToken, ScopeDelimiterToken,
		ForallToken, IdentifierToken, ScopeDelimiterToken,
		ExpressionToken, ExpressionDelimiterToken,
		GuaranteeToken, ScopeDelimiterToken,
		ForallToken, IdentifierToken, ScopeDelimiterToken,
		ExpressionToken, ExpressionDelimiterToken,
		RegionToken, IdentifierToken, ScopeDelimiterToken,
		AssumeToken, ScopeDelimiterToken,
		ForallToken, IdentifierToken, ScopeDelimiterToken,
		ExpressionToken, ExpressionDelimiterToken,
		GuaranteeToken, ScopeDelimiterToken,
		ExistsToken, IdentifierToken, ScopeDelimiterToken,
		ExpressionToken, ExpressionDelimiterToken,
		EofToken,
	}

	tokens := iterx.Collect(LexString(source))
	classes := make([]TokenClass, len(tokens))
	for idx := range tokens {
		classes[idx] = tokens[idx].class
	}
	assert.ElementsMatch(t, classes, expected)
}

func TestLexClasses(t *testing.T) {
	tests := []struct {
		description string
		input       string
		classes     []TokenClass
	}{
		{
			description: "region with name",
			input:       "region dasa ba123 c :",
			classes:     []TokenClass{RegionToken, IdentifierToken, IdentifierToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "assumption",
			input:       "assume   :   ",
			classes:     []TokenClass{AssumeToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "guarantee",
			input:       "guarantee:",
			classes:     []TokenClass{GuaranteeToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "forall",
			input:       "forall a  b d   .",
			classes:     []TokenClass{ForallToken, IdentifierToken, IdentifierToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "guarantee forall",
			input:       "guarantee: forall a  b d   .",
			classes:     []TokenClass{GuaranteeToken, ScopeDelimiterToken, ForallToken, IdentifierToken, IdentifierToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "guarantee forall exists",
			input:       "guarantee: forall ab  basd c   . exists d.",
			classes:     []TokenClass{GuaranteeToken, ScopeDelimiterToken, ForallToken, IdentifierToken, IdentifierToken, IdentifierToken, ScopeDelimiterToken, ExistsToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "empty expression",
			input:       "   ;  ",
			classes:     []TokenClass{ExpressionDelimiterToken, EofToken},
		},
		{
			description: "right parenthesis",
			input:       "   )  ",
			classes:     []TokenClass{RightParenthesis, EofToken},
		},
		{
			description: "left forall parenthesis",
			input:       "   (  forall a. )",
			classes:     []TokenClass{LeftParenthesis, ForallToken, IdentifierToken, ScopeDelimiterToken, RightParenthesis, EofToken},
		},
		{
			description: "left forall parenthesis",
			input:       "   (  forall a. forall b. )",
			classes:     []TokenClass{LeftParenthesis, ForallToken, IdentifierToken, ScopeDelimiterToken, ForallToken, IdentifierToken, ScopeDelimiterToken, RightParenthesis, EofToken},
		},
		{
			description: "empty input",
			input:       "",
			classes:     []TokenClass{EofToken},
		},
		{
			description: "grouped empty expression",
			input:       "( ; )",
			classes:     []TokenClass{LeftParenthesis, ExpressionDelimiterToken, RightParenthesis, EofToken},
		},
		{
			description: "Single unnamed region",
			input:       "region :",
			classes:     []TokenClass{RegionToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "Single named region",
			input:       "region Positive :",
			classes:     []TokenClass{RegionToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "Two named regions",
			input:       "region Positive : region Negative :",
			classes:     []TokenClass{RegionToken, IdentifierToken, ScopeDelimiterToken, RegionToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "No region just a quantifier",
			input:       "forall e0 .",
			classes:     []TokenClass{ForallToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "No region two quantifiers",
			input:       "forall e0 . exists e1 .",
			classes:     []TokenClass{ForallToken, IdentifierToken, ScopeDelimiterToken, ExistsToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "assumption with quantifier",
			input:       "assume: forall e0.",
			classes:     []TokenClass{AssumeToken, ScopeDelimiterToken, ForallToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "two assumptions with a quantifier and no region",
			input:       "assume: forall e0. assume: exists e0.",
			classes:     []TokenClass{AssumeToken, ScopeDelimiterToken, ForallToken, IdentifierToken, ScopeDelimiterToken, AssumeToken, ScopeDelimiterToken, ExistsToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "region without identifiers",
			input:       "region :",
			classes:     []TokenClass{RegionToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "region without identifiers newline seperated",
			input:       "region:",
			classes:     []TokenClass{RegionToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "region with an identifier",
			input:       "region a :",
			classes:     []TokenClass{RegionToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "region with two identifiers",
			input:       "region a b :",
			classes:     []TokenClass{RegionToken, IdentifierToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "forall without identifiers",
			input:       "forall a.",
			classes:     []TokenClass{ForallToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "empty exists",
			input:       "exists a.",
			classes:     []TokenClass{ExistsToken, ScopeDelimiterToken, IdentifierToken, EofToken},
		},
		{
			description: "assume with true",
			input:       "assume: true ;",
			classes:     []TokenClass{AssumeToken, ScopeDelimiterToken, ExpressionToken, ExpressionDelimiterToken, EofToken},
		},
		{
			description: "assume with large expression",
			input:       "assume: a == b && check(a, b) ;",
			classes:     []TokenClass{AssumeToken, ScopeDelimiterToken, ExpressionToken, ExpressionDelimiterToken, EofToken},
		},
		{
			description: "guarantee with large expression",
			input:       "guarantee: a == b && a == c ;",
			classes:     []TokenClass{GuaranteeToken, ScopeDelimiterToken, ExpressionToken, ExpressionDelimiterToken, EofToken},
		},
		{
			description: "assume with expression delimiter but no expression",
			input:       "assume: ;",
			classes:     []TokenClass{AssumeToken, ScopeDelimiterToken, ExpressionDelimiterToken, EofToken},
		},
		{
			description: "empty region, assumption, and forall",
			input:       "region : assume : ; forall a.",
			classes:     []TokenClass{RegionToken, ScopeDelimiterToken, AssumeToken, ScopeDelimiterToken, ExpressionDelimiterToken, ForallToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "region with name newline delimiter",
			input:       "region SomeName:",
			classes:     []TokenClass{RegionToken, IdentifierToken, ScopeDelimiterToken, EofToken},
		},
		{
			description: "named region, forall, and assumption",
			input:       "region SomeName: forall a . assume: false && true",
			classes:     []TokenClass{RegionToken, IdentifierToken, ScopeDelimiterToken, ForallToken, IdentifierToken, ScopeDelimiterToken, AssumeToken, ScopeDelimiterToken, ExpressionToken, ExpressionDelimiterToken, EofToken},
		},
		{
			description: "multiple named regions",
			input:       "region positive: assume: true; region negative: assume: false",
			classes:     []TokenClass{RegionToken, IdentifierToken, ScopeDelimiterToken, AssumeToken, ScopeDelimiterToken, ExpressionToken, ExpressionDelimiterToken, RegionToken, IdentifierToken, ScopeDelimiterToken, AssumeToken, ScopeDelimiterToken, ExpressionToken, ExpressionDelimiterToken, EofToken},
		},
		{
			description: "multiple named regions",
			input:       "assume: false",
			classes:     []TokenClass{AssumeToken, ScopeDelimiterToken, ExpressionToken, ExpressionDelimiterToken, EofToken},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tokens := iterx.Collect(LexString(tt.input))
			classes := make([]TokenClass, len(tokens))
			for idx := range tokens {
				classes[idx] = tokens[idx].class
			}
			assert.ElementsMatch(t, classes, tt.classes)
		})
	}
}

func TestLexGo(t *testing.T) {
	sourceMulti := `
package main

/*
region Positive
assume: forall e0. e0.in >= 0
guarantee: forall e0. e0.ret0 >= 0
region Negative
assume: forall e0. e0.in < 0
guarantee: forall e0. e0.ret0 < 0
*/
func Self(in a) int {
	return in
}
`
	sourceSingle := `
package main

// region Positive
// assume: forall e0. e0.in >= 0
// guarantee: 
// forall e0. e0.ret0 >= 0
// region Negative
// assume: forall e0. e0.in < 0
// guarantee: forall e0.
// 			e0.ret0 < 0
func Self(in a) int {
	return in
}
`

	fset := token.NewFileSet()

	fileMulti, err := parser.ParseFile(fset, "", sourceMulti, parser.ParseComments)
	assert.Nil(t, err)
	tokensMulti := iterx.Collect(LexGo(fileMulti.Decls[0].(*ast.FuncDecl).Doc))

	fileSingle, err := parser.ParseFile(fset, "", sourceSingle, parser.ParseComments)
	assert.Nil(t, err)
	tokensSingle := iterx.Collect(LexGo(fileSingle.Decls[0].(*ast.FuncDecl).Doc))

	assert.ElementsMatch(t, tokensMulti, tokensSingle)
}

/*func FuzzLexString(f *testing.F) {
	f.Add("guarantee a == b && a == c ;")
	f.Add("region SomeName \n forall a . assume false && true")
	f.Add("forall &&")
	f.Add("guarantee;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;")
	f.Add("\x9b")
	f.Add("�")
	f.Add("region .")
	f.Add("region Positive .")
	f.Add("region Positive . region Negative .")
	f.Add("forall e0 .")
	f.Add("forall e0 . exists e1 .")
	f.Add("assume forall e0 .")
	f.Add("assume forall e0 . assume exists e0 .")
	f.Add("region .")
	f.Add("region \n")
	f.Add("region a .")
	f.Add("region a b .")
	f.Add("forall .")
	f.Add("exists .")
	f.Add("assume true ;")
	f.Add("assume a == b && check(a, b) ;")
	f.Add("guarantee a == b && a == c ;")
	f.Add("assume ;")
	f.Add("region . assume ; forall .")
	f.Add("region SomeName \n")
	f.Add("region SomeName \n forall a . assume false && true")
	f.Add("region positive . assume true ; \n region negative . assume false")
	f.Add("assume false")
	f.Fuzz(func(t *testing.T, str string) {
		iterx.Collect(LexString(str))
	})
}*/
