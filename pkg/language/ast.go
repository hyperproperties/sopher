package language

type Node interface{}

type Contract struct {
	regions []Region
}

type Region struct {
	name        []string
	assumptions []Node
	guarantees  []Node
}

type Universal struct {
	variables []string
	assertion Node
}

type Existential struct {
	variables []string
	assertion Node
}

type Assumption struct {
	assertion Node
}

type Guarantee struct {
	assertion Node
}

type GoExpresion struct {
	code string
}

type BinaryExpression struct {
	lhs, rhs Node
	operator BinaryOperator
}

type UnaryExpression struct {
	operator UnaryOperator
	operand  Node
}

type PredicateExpression[T any] struct {
	predicate func(assignments []T) bool
}

type ProbabilisticQuantifier struct {
	event Node // P( event )
}

type ConditionalProbabilityQuantifier struct {
	event, given Node // P( event | given )
}

type ConstantNumber struct {
	value float32
}

type Group struct {
	node Node
}

func NewContract(regions ...Region) Contract {
	return Contract{
		regions: regions,
	}
}

func NewRegion(name []string, assumptions, guarantees []Node) Region {
	return Region{
		name:        name,
		assumptions: assumptions,
		guarantees:  guarantees,
	}
}

func NewUniversal(variables []string, assertion Node) Universal {
	return Universal{
		variables: variables,
		assertion: assertion,
	}
}

func NewExistential(variables []string, assertion Node) Existential {
	return Existential{
		variables: variables,
		assertion: assertion,
	}
}

func NewAssumption(assertion Node) Assumption {
	return Assumption{
		assertion: assertion,
	}
}

func NewGuarantee(assertion Node) Guarantee {
	return Guarantee{
		assertion: assertion,
	}
}

func NewProbabilisticQuantifier(body Node) ProbabilisticQuantifier {
	return ProbabilisticQuantifier{
		event: body,
	}
}

func NewConditionalProbabilityQuantifier(event, given Node) ConditionalProbabilityQuantifier {
	return ConditionalProbabilityQuantifier{
		event: event,
		given: given,
	}
}

func NewGoExpression(code string) GoExpresion {
	return GoExpresion{
		code: code,
	}
}

func Number(value float32) ConstantNumber {
	return ConstantNumber{value}
}

func NewPredicateExpression[T any](predicate func(assignments []T) bool) PredicateExpression[T] {
	return PredicateExpression[T]{
		predicate: predicate,
	}
}

func NewGroup(node Node) Group {
	return Group{
		node: node,
	}
}

func NewBinaryExpression(lhs Node, operator BinaryOperator, rhs Node) BinaryExpression {
	return BinaryExpression{
		lhs:      lhs,
		operator: operator,
		rhs:      rhs,
	}
}

func NewUnaryExpression(operator UnaryOperator, operand Node) UnaryExpression {
	return UnaryExpression{
		operator: operator,
		operand:  operand,
	}
}
