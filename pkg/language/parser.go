package language

import (
	"iter"

	"github.com/hyperproperties/sopher/pkg/iterx"
)

type Parser struct {
	next func() (Token, bool)
	peek func(lookahead int) (Token, bool)
	stop func()
}

func NewParser(tokens iter.Seq[Token]) Parser {
	next, peek, stop := iterx.BufferedPull(tokens)
	return Parser{
		next: next,
		peek: peek,
		stop: stop,
	}
}

func (parser *Parser) any(classes ...TokenClass) (Token, bool) {
	token, exists := parser.peek(1)
	if !exists {
		return Token{}, false
	}

	for _, class := range classes {
		if token.class == class {
			return token, true
		}
	}

	return Token{}, false
}

func (parser *Parser) match(classes ...TokenClass) bool {
	_, exists := parser.any(classes...)
	return exists
}

func (parser *Parser) consume(classes ...TokenClass) (Token, bool) {
	if token, exists := parser.any(classes...); exists {
		parser.next()
		return token, exists
	}
	return Token{}, false
}

func (parser *Parser) Parse() Node {
	return parser.contract()
}

func (parser *Parser) contract() Contract {
	var regions []Region

	if parser.match(AssumeToken, GuaranteeToken) {
		obligations := parser.obligations()
		region := NewRegion(nil, obligations)
		regions = append(regions, region)
	}

	for {
		if !parser.match(RegionToken) {
			break
		}
		regions = append(regions, parser.region())
	}

	return NewContract(regions...)
}

func (parser *Parser) region() Region {
	parser.consume(RegionToken)

	var name []string
	for {
		if identifier, exists := parser.consume(IdentifierToken); exists {
			name = append(name, identifier.lexeme)
		} else {
			break
		}
	}

	if _, exists := parser.consume(ScopeDelimiterToken); !exists {
		panic("expected a scope delimiter token")
	}

	obligations := parser.obligations()

	return Region{
		name:        name,
		obligations: obligations,
	}
}

func (parser *Parser) obligations() (obligations []Node) {
	for {
		if parser.match(AssumeToken) {
			obligations = append(obligations, parser.assumption())
		} else if parser.match(GuaranteeToken) {
			obligations = append(obligations, parser.guarantee())
		} else {
			break
		}
	}

	return obligations
}

func (parser *Parser) assumption() Assumption {
	if _, ok := parser.consume(AssumeToken); !ok {
		panic("assume token expected for assumption")
	}

	if _, ok := parser.consume(ScopeDelimiterToken); !ok {
		panic("assume expected scope delimiter toke")
	}

	assertion := parser.assertion()

	return NewAssumption(assertion)
}

func (parser *Parser) guarantee() Guarantee {
	if _, ok := parser.consume(GuaranteeToken); !ok {
		panic("guarantee token expected for guarantee")
	}

	if _, ok := parser.consume(ScopeDelimiterToken); !ok {
		panic("guarantee expected scope delimiter toke")
	}

	assertion := parser.assertion()

	return NewGuarantee(assertion)
}

func (parser *Parser) assertion() Node {
	switch {
	case parser.match(ForallToken):
		return parser.universal()
	case parser.match(ExistsToken):
		return parser.existential()
	case parser.match(ExpressionToken):
		return parser.expression()
	case parser.match(LeftParenthesis):
		return parser.group()
	}
	panic("unknown assertion")
}

func (parser *Parser) group() Group {
	if _, ok := parser.consume(LeftParenthesis); !ok {
		panic("group expected left parenthesis")
	}

	assertion := parser.assertion()

	if _, ok := parser.consume(RightParenthesis); !ok {
		panic("group expected right parenthesis")
	}
	return NewGroup(assertion)
}

func (parser *Parser) variables() (variables []string) {
	for {
		identifier, exists := parser.consume(IdentifierToken)
		if !exists {
			break
		}
		variables = append(variables, identifier.lexeme)
	}

	return variables
}

func (parser *Parser) universal() Universal {
	if _, ok := parser.consume(ForallToken); !ok {
		panic("forall token expected for universal quantifier")
	}

	variables := parser.variables()

	if _, ok := parser.consume(ScopeDelimiterToken); !ok {
		panic("expected scope delimiter toke")
	}

	assertion := parser.assertion()

	return NewUniversal(variables, assertion)
}

func (parser *Parser) existential() Existential {
	if _, ok := parser.consume(ExistsToken); !ok {
		panic("forall token expected for universal quantifier")
	}

	variables := parser.variables()

	if _, ok := parser.consume(ScopeDelimiterToken); !ok {
		panic("expected scope delimiter toke")
	}

	assertion := parser.assertion()

	return NewExistential(variables, assertion)
}

func (parser *Parser) expression() (expression Node) {
	code, ok := parser.consume(ExpressionToken)
	if !ok {
		panic("go expression expected expression token")
	}

	if _, ok = parser.consume(ExpressionDelimiterToken); !ok {
		panic("missing expression delimiter token")
	}

	return GoExpresion{
		code: code.lexeme,
	}
}
