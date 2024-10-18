package language

import "fmt"

type TokenClass uint8

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
		return "scope delimiter"
	case ExpressionDelimiterToken:
		return "eexpression delimiter"
	case LeftParenthesis:
		return "("
	case RightParenthesis:
		return ")"
	case EofToken:
		return "eof"
	}
	return fmt.Sprintf("%v", uint8(class))
}

const (
	RegionToken = TokenClass(iota)
	ForallToken
	ExistsToken
	AssumeToken
	GuaranteeToken
	IdentifierToken
	ProbabilityToken
	ExpressionToken
	ScopeDelimiterToken
	ExpressionDelimiterToken
	LeftParenthesis
	RightParenthesis
	EofToken
)

type Token struct {
	class  TokenClass
	lexeme string
}

func NewToken(class TokenClass, lexeme string) Token {
	return Token{
		class:  class,
		lexeme: lexeme,
	}
}
