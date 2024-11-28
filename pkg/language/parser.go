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

func (parser *Parser) tryConsume(classes ...TokenClass) bool {
	_, exists := parser.consume(classes...)
	return exists
}

func (parser *Parser) Parse() Contract {
	return parser.contract()
}

func (parser *Parser) contract() Contract {
	var regions []Region

	if parser.match(AssumeToken, GuaranteeToken) {
		assumptions, guarantees := parser.obligations()
		region := NewRegion(nil, assumptions, guarantees)
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

	assumptions, guarantees := parser.obligations()

	return Region{
		name:        name,
		assumptions: assumptions,
		guarantees:  guarantees,
	}
}

func (parser *Parser) obligations() (assumptions, guarantees []Node) {
	for {
		if parser.match(AssumeToken) {
			assumptions = append(assumptions, parser.assumption())
		} else if parser.match(GuaranteeToken) {
			guarantees = append(guarantees, parser.guarantee())
		} else {
			break
		}
	}

	return assumptions, guarantees
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
	return parser.biimplication()
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

func (parser *Parser) biimplication() Node {
	lhs := parser.implication()

	for parser.tryConsume(LogicalBiimplicationToken) {
		rhs := parser.implication()
		lhs = NewBinaryExpression(lhs, LogicalBiimplication, rhs)
	}

	return lhs
}

func (parser *Parser) implication() Node {
	lhs := parser.conjunction()

	for parser.tryConsume(LogicalImplicationToken) {
		rhs := parser.conjunction()
		lhs = NewBinaryExpression(lhs, LogicalImplication, rhs)
	}

	return lhs
}

func (parser *Parser) conjunction() Node {
	lhs := parser.disjunction()

	for parser.tryConsume(LogicalConjunctionToken) {
		rhs := parser.disjunction()
		lhs = NewBinaryExpression(lhs, LogicalConjunction, rhs)
	}

	return lhs
}

func (parser *Parser) disjunction() Node {
	lhs := parser.negation()

	for parser.tryConsume(LogicalDisjunctionToken) {
		rhs := parser.negation()
		lhs = NewBinaryExpression(lhs, LogicalDisjunction, rhs)
	}

	return lhs
}

func (parser *Parser) negation() Node {
	if parser.tryConsume(LogicalNegationToken) {
		operand := parser.negation()
		return NewUnaryExpression(LogicalNegation, operand)
	}

	return parser.quantifier()
}

func (parser *Parser) quantifier() Node {
	if token, ok := parser.consume(ForallToken, ExistsToken); ok {
		variables := parser.variables()
		if _, ok := parser.consume(ScopeDelimiterToken); !ok {
			panic("expected scope delimiter token")
		}

		assertion := parser.quantifier()

		switch token.class {
		case ForallToken:
			return NewUniversal(variables, assertion)
		case ExistsToken:
			return NewExistential(variables, assertion)
		}
	}

	return parser.expression()
}

func (parser *Parser) expression() Node {
	if parser.match(LeftParenthesisToken) {
		return parser.group()
	}

	// Starting delimiter is optional.
	parser.tryConsume(GoExpressionDelimiterToken)

	code, ok := parser.consume(GoExpressionToken)
	if !ok {
		panic("go expression expected expression token")
	}

	if ok := parser.tryConsume(GoExpressionDelimiterToken); !ok {
		panic("missing ending expression delimiter token")
	}

	return NewGoExpression(code.lexeme)
}

func (parser *Parser) group() Group {
	if _, ok := parser.consume(LeftParenthesisToken); !ok {
		panic("group expected left parenthesis")
	}

	assertion := parser.assertion()

	if _, ok := parser.consume(RightParenthesisToken); !ok {
		panic("group expected right parenthesis")
	}
	return NewGroup(assertion)
}
