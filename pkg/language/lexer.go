package language

import (
	"fmt"
	"iter"
	"strings"
	"unicode"

	"github.com/hyperproperties/sopher/pkg/iterx"
)

type Lexer struct {
	next func() (rune, bool)
	peek func(lookahread int) (rune, bool)
	stop func()
}

func NewLexer(runes iter.Seq[rune]) Lexer {
	next, peek, stop := iterx.BufferedPull(runes)
	return Lexer{
		next: next,
		peek: peek,
		stop: stop,
	}
}

func LexComments(docs []string) iter.Seq[Token] {
	var builder strings.Builder

	for _, text := range docs {
		if rest, hasPrefix := strings.CutPrefix(text, "/*"); hasPrefix {
			if rest, hasSuffix := strings.CutSuffix(rest, "*/"); hasSuffix {
				builder.WriteString(rest)
			} else {
				panic("expected go multi-line comment to has both prefix and suffix")
			}
		} else if rest, hasPrefix := strings.CutPrefix(text, "//"); hasPrefix {
			builder.WriteString(rest)
			builder.WriteRune('\n')
		}
	}

	return LexString(builder.String())
}

func LexString(str string) iter.Seq[Token] {
	return LexRunes([]rune(str))
}

func LexRunes(runes []rune) iter.Seq[Token] {
	lexer := NewLexer(iterx.Forward(runes))
	return lexer.Scan()
}

func (lexer *Lexer) peekWords(
	prefix string,
	skip func(rune) bool,
	body func(string) bool,
	suffixs ...string,
) (bool, int, []string) {
	var words []string

	lookahead := 1
	runes := []rune(prefix)
	for offset := 0; offset < len(prefix); offset++ {
		character, ok := lexer.peek(lookahead)
		if !ok {
			return false, lookahead - 1, words
		}

		if runes[offset] != character {
			return false, lookahead, words
		}

		lookahead++
	}

	words = append(words, prefix)
	var builder strings.Builder

	for ; ; lookahead++ {
		character, ok := lexer.peek(lookahead)

		// Append the current work we are building if there are
		// no more runes or we encountered a skip.
		if !ok || skip(character) {

			// We dont have to consider suffix as it is handle after.
			if builder.Len() > 0 {
				word := builder.String()
				if body(word) {
					words = append(words, word)
				}

				builder.Reset()
			}

			if !ok {
				break
			}
		}

		if !skip(character) {
			builder.WriteRune(character)
			word := builder.String()

			for _, suffix := range suffixs {
				if prefix, hasSuffix := strings.CutSuffix(word, suffix); hasSuffix {
					if len(prefix) > 0 {
						if body(prefix) {
							words = append(words, prefix)
							words = append(words, suffix)
							return true, lookahead, words
						}
					} else {
						words = append(words, suffix)
						return true, lookahead, words
					}
				}
			}
		}
	}

	return false, lookahead - 1, words
}

func (lexer *Lexer) consumeWords(
	prefix string,
	skip func(rune) bool,
	body func(string) bool,
	suffixs ...string,
) (bool, []string) {
	if found, lookahead, words := lexer.peekWords(prefix, skip, body, suffixs...); found {
		var builder strings.Builder
		for idx := 0; idx < lookahead; idx++ {
			character, ok := lexer.next()
			if !ok {
				panic("consuming in lookahead failed")
			}
			builder.WriteRune(character)
		}
		return found, words
	}

	return false, []string{}
}

func (lexer *Lexer) peekWord(word string) bool {
	lookahead := 1
	runes := []rune(word)
	for offset := 0; offset < len(word); offset++ {
		character, ok := lexer.peek(lookahead)
		if !ok {
			return false
		}

		if runes[offset] != character {
			return false
		}

		lookahead++
	}

	return true
}

func (lexer *Lexer) consumeWord(word string) bool {
	if found := lexer.peekWord(word); found {
		for idx := 0; idx < len(word); idx++ {
			lexer.next()
		}
		return found
	}
	return false
}

func (lexer *Lexer) isIdentifier(str string) bool {
	for idx, character := range []rune(str) {
		if idx == 0 && unicode.IsNumber(character) {
			return false
		}
		if !(unicode.IsLetter(character) || character == '_' || unicode.IsNumber(character)) {
			return false
		}

	}
	return true
}

func (lexer *Lexer) region(fields []string) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		length := len(fields)
		if length < 2 {
			panic("region expected no less than two fields")
		}

		if fields[0] != "region" {
			panic("region expected to start with \"region\"")
		}

		if !yield(NewToken(RegionToken, fields[0])) {
			return
		}

		for idx := 1; idx < length-1; idx++ {
			if !yield(NewToken(IdentifierToken, fields[idx])) {
				return
			}
		}

		yield(NewToken(ScopeDelimiterToken, fields[length-1]))
	}
}

func (lexer *Lexer) assume(fields []string) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		if len(fields) != 2 {
			panic("assume expected only two fields")
		}

		if fields[0] != "assume" {
			panic("assume expected to start with \"assume\"")
		}

		if !yield(NewToken(AssumeToken, fields[0])) {
			return
		}

		if !yield(NewToken(ScopeDelimiterToken, fields[1])) {
			return
		}
	}
}

func (lexer *Lexer) guarantee(fields []string) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		if len(fields) != 2 {
			panic("guarantee expected two fields")
		}

		if fields[0] != "guarantee" {
			panic("guarantee expected to start with \"guarantee\"")
		}

		if !yield(NewToken(GuaranteeToken, fields[0])) {
			return
		}

		if !yield(NewToken(ScopeDelimiterToken, fields[1])) {
			return
		}
	}
}

func (lexer *Lexer) quantifier(prefix string, class TokenClass, fields []string) iter.Seq[Token] {
	return func(yield func(Token) bool) {
		length := len(fields)

		if fields[0] == "(" {
			// mimum example "(" "forall" "e" "."
			if len(fields) < 4 {
				panic(fmt.Sprintf("grouped %s has a minimum length of 4", class))
			}
			if fields[1] != prefix {
				panic(fmt.Sprintf("%s expected to have \"%s\" after \"(\"", class, prefix))
			}
		} else {
			// mimum example "forall" "e" "."
			if len(fields) < 3 {
				panic(fmt.Sprintf("%s has a minimum length of 3", class))
			}
			if fields[0] != prefix {
				panic(fmt.Sprintf("%s expected to start with \"%s\"", class, prefix))
			}
		}

		if !yield(NewToken(class, fields[0])) {
			return
		}

		for idx := 1; idx < length-1; idx++ {
			if !yield(NewToken(IdentifierToken, fields[idx])) {
				return
			}
		}

		yield(NewToken(ScopeDelimiterToken, fields[length-1]))
	}
}

func (lexer *Lexer) forall(fields []string) iter.Seq[Token] {
	return lexer.quantifier("forall", ForallToken, fields)
}

func (lexer *Lexer) exists(fields []string) iter.Seq[Token] {
	return lexer.quantifier("exists", ExistsToken, fields)
}

func (lexer *Lexer) expression() iter.Seq[Token] {
	return func(yield func(Token) bool) {
		var builder strings.Builder
		for {
			character, ok := lexer.peek(1)
			if !ok {
				yield(NewToken(GoExpressionToken, builder.String()))
				yield(NewToken(GoExpressionDelimiterToken, ";"))
				return
			}

			lexer.next()

			if character == ';' || character == '\n' {
				if !yield(NewToken(GoExpressionToken, builder.String())) {
					return
				}
				yield(NewToken(GoExpressionDelimiterToken, ";"))
				return
			} else {
				builder.WriteRune(character)
			}
		}
	}
}

func (lexer *Lexer) whitespace() {
	for {
		character, ok := lexer.peek(1)
		if !ok {
			break
		}
		if unicode.IsSpace(character) {
			lexer.next()
			continue
		}
		break
	}
}

func (lexer *Lexer) isSpace(r rune) bool {
	switch r {
	case '\t', '\v', '\f', '\r', ' ', 0x85, 0xA0:
		return true
	}
	return false
}

func (lexer *Lexer) Scan() iter.Seq[Token] {
	keycharacters := map[rune]Token{
		'!': NewToken(LogicalNegationToken, ";"),
		';': NewToken(GoExpressionDelimiterToken, ";"),
		'(': NewToken(LeftParenthesisToken, "("),
		')': NewToken(RightParenthesisToken, ")"),
	}

	keywords := map[string]Token{
		"&&":  NewToken(LogicalConjunctionToken, "&&"),
		"||":  NewToken(LogicalDisjunctionToken, "||"),
		"->":  NewToken(LogicalImplicationToken, "->"),
		"<->": NewToken(LogicalBiimplicationToken, "<->"),
	}

	return func(yield func(Token) bool) {
		for {
			lexer.whitespace()
			character, ok := lexer.peek(1)
			if !ok {
				break
			}

			isKeyword := false
			for keyword, token := range keywords {
				if found := lexer.consumeWord(keyword); found {
					isKeyword = true
					if !yield(token) {
						return
					}
				}
			}

			if isKeyword {
				continue
			}

			if token, found := keycharacters[character]; found {
				lexer.next()
				if !yield(token) {
					return
				}
			} else if found, words := lexer.consumeWords(
				"region",
				lexer.isSpace,
				lexer.isIdentifier,
				":", "\n",
			); found {
				region := lexer.region(words)
				if !iterx.Pipe(region, yield) {
					return
				}
			} else if found, words := lexer.consumeWords(
				"assume",
				lexer.isSpace,
				func(string) bool { return false },
				":", "\n",
			); found {
				assume := lexer.assume(words)
				if !iterx.Pipe(assume, yield) {
					return
				}
			} else if found, words := lexer.consumeWords(
				"guarantee",
				lexer.isSpace,
				func(string) bool { return false },
				":", "\n",
			); found {
				guarantee := lexer.guarantee(words)
				if !iterx.Pipe(guarantee, yield) {
					return
				}
			} else if found, words := lexer.consumeWords(
				"forall",
				lexer.isSpace,
				lexer.isIdentifier,
				".", "\n",
			); found {
				forall := lexer.forall(words)
				if !iterx.Pipe(forall, yield) {
					return
				}
			} else if found, words := lexer.consumeWords(
				"exists",
				lexer.isSpace,
				lexer.isIdentifier,
				".", "\n",
			); found {
				exists := lexer.exists(words)
				if !iterx.Pipe(exists, yield) {
					return
				}
			} else {
				expression := lexer.expression()
				if !iterx.Pipe(expression, yield) {
					return
				}
			}
		}

		yield(NewToken(EndToken, ""))
	}
}
