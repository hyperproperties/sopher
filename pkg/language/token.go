package language

import "fmt"

// A class which are used to group lexemes by.
type TokenClass uint8

const (
	// The beginning of a region declaration.
	RegionToken = TokenClass(iota)
	// The beginning of a universal quantifier declaration.
	ForallToken
	// The beginning of an existential quantifier declaration.
	ExistsToken
	// Represents the end of declaring quantifier variables.
	ScopeDelimiterToken
	// The beginning of a probability check "quantifier" declaration.
	ProbabilityToken
	// The declaration of an assumption.
	AssumeToken
	// The declaration of a guarantee.
	GuaranteeToken
	// An identifier which is declared by a quantifier and not in an expression.
	IdentifierToken
	// Represents a golang expression which are the hyper assertion predicates.
	ExpressionToken
	// The delimiter for golang expressions.
	ExpressionDelimiterToken
	// Left parenthesis "(".
	LeftParenthesis
	// Left parenthesis ")".
	RightParenthesis
	// Represents the end of the contract.
	EndToken
)

// Returns a stringifyed version of the token class.
func (class TokenClass) String() string {
	switch class {
	case RegionToken:
		return "region"
	case ForallToken:
		return "forall"
	case ExistsToken:
		return "exists"
	case AssumeToken:
		return "assume"
	case GuaranteeToken:
		return "guarantee"
	case IdentifierToken:
		return "identifier"
	case ProbabilityToken:
		return "probability"
	case ExpressionToken:
		return "expression"
	case ScopeDelimiterToken:
		return ":"
	case ExpressionDelimiterToken:
		return ";"
	case LeftParenthesis:
		return "("
	case RightParenthesis:
		return ")"
	case EndToken:
		return "end"
	}
	return fmt.Sprintf("%v", uint8(class))
}

// Represents the smallest sequence of meaningful characters
// bundled together and grouped by class.
type Token struct {
	class  TokenClass
	lexeme string
}

// Returns a new token which is required to have a valid token class.
func NewToken(class TokenClass, lexeme string) Token {
	if class < RegionToken || class > EndToken {
		panic("invalid token class")
	}

	return Token{
		class:  class,
		lexeme: lexeme,
	}
}

// Returns the class of the token.
func (token Token) Class() TokenClass {
	return token.class
}

// Returns the lexeme of the token.
func (token Token) Lexeme() string {
	return token.lexeme
}
