package language

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/hyperproperties/sopher/pkg/dstx"
)

type AssertionFactory struct {
	packageName string
	modelName   string
	variables   Stack[string]
}

func NewAssertionFactory(packageName, modelName string) AssertionFactory {
	return AssertionFactory{
		packageName: packageName,
		modelName:   modelName,
		variables:   make([]string, 0),
	}
}

func (factory *AssertionFactory) Create(node Node) *dst.CallExpr {
	switch cast := node.(type) {
	case GoExpresion:
		return factory.NewPredicate(cast)
	case Universal:
		return factory.NewUniversal(cast)
	case Existential:
		return factory.NewExistential(cast)
	case Guarantee:
		return factory.Create(cast.assertion)
	case Assumption:
		return factory.Create(cast.assertion)
	case BinaryExpression:
		return factory.NewBinary(cast)
	case Group:
		return factory.Create(cast.node)
	case UnaryExpression:
		return factory.NewUnary(cast)
	}
	panic(fmt.Sprintf("unknown node type %t", node))
}

func (factory *AssertionFactory) NewBinary(binary BinaryExpression) *dst.CallExpr {
	operators := map[BinaryOperator]string{
		LogicalConjunction:   "LogicalConjunction",
		LogicalDisjunction:   "LogicalDisjunction",
		LogicalImplication:   "LogicalImplication",
		LogicalBiimplication: "LogicalBiimplication",
	}
	var operator dst.Expr
	if name, ok := operators[binary.operator]; ok {
		operator = dstx.SelectS(name).FromS(factory.packageName)
	} else {
		panic("unknown binary operator")
	}

	dstx.AppendEnd(operator, fmt.Sprintf("/* %s */", binary.operator.String()))

	call := dstx.Call(
		dstx.IndexS(factory.modelName).Of(
			dstx.SelectS("NewBinaryHyperAssertion").FromS(factory.packageName),
		),
	).Pass(factory.Create(binary.lhs), operator, factory.Create(binary.rhs))

	return dstx.NewLineAround(call)
}

func (factory *AssertionFactory) NewUnary(unary UnaryExpression) *dst.CallExpr {
	operators := map[UnaryOperator]string{
		LogicalNegation: "LogicalNegation",
	}
	var operator dst.Expr
	if name, ok := operators[unary.operator]; ok {
		operator = dstx.SelectS(name).FromS(factory.packageName)
	} else {
		panic("unknown unary operator")
	}

	dstx.AppendEnd(operator, fmt.Sprintf("/* %s */", unary.operator.String()))

	call := dstx.Call(
		dstx.IndexS(factory.modelName).Of(
			dstx.SelectS("NewUnaryHyperAssertion").FromS(factory.packageName),
		),
	).Pass(operator, factory.Create(unary.operand))

	return dstx.NewLineAround(call)
}

func (factory *AssertionFactory) NewPredicate(expression GoExpresion) *dst.CallExpr {
	// Construct the execution variable definitions.
	// e0, e1, e2 := assignments[0], assignments[1], assignments[2]
	definitions := dstx.
		DefineS(factory.variables...).
		As(dstx.Construct[[]dst.Expr](factory.variables, func(i int, _ string) dst.Expr {
			return dstx.Index(dstx.BasicInt(i)).OfS("assignments")
		})...)

	// Construct anonymous execution variable assignments
	// This ensures all variables are in use and that the compiler wont complain.
	// _, _, _ = e0, e1, e2
	anonymous := dstx.
		AssignS(dstx.RepeatS("_", len(factory.variables))...).
		ToS(factory.variables...)

	// Parse the code of the expression to an ast node.
	expr, err := parser.ParseExpr(expression.code)
	if err != nil {
		panic(err)
	}

	// Convert the ast node to a decorated dst node.
	decorator := decorator.NewDecorator(token.NewFileSet())
	node, err := decorator.DecorateNode(expr)
	if err != nil {
		panic(err)
	}

	// Create the return statement of the hyper assertion predicate.
	returnStmt := dstx.Return(node.(dst.Expr))

	// Construct the predicate of the hyper assertion.
	predicate := dstx.Function(
		dstx.
			TakingN(dstx.FieldS("assignments").
				Type(dstx.SliceTypeS(factory.modelName))).
			ResultsN(dstx.Field().TypeS("bool")),
		dstx.Block(
			dstx.NewLineAround(definitions),
			dstx.NewLineAround(anonymous),
			dstx.NewLineAround(returnStmt),
		),
	)

	call := dstx.
		Call(dstx.SelectS("NewPredicateHyperAssertion").
			FromS(factory.packageName)).
		Pass(dstx.NewLineAround(predicate))

	return dstx.NewLineAround(call)
}

func (factory *AssertionFactory) NewUniversal(universal Universal) *dst.CallExpr {
	factory.variables.Push(universal.variables...)
	call := dstx.Call(
		dstx.SelectS("NewUniversalHyperAssertion").FromS(factory.packageName),
	).Pass(
		dstx.BasicInt(len(universal.variables)),
		factory.Create(universal.assertion),
	)
	factory.variables.PopN(len(universal.variables))
	
	return dstx.NewLineAround(call)
}

func (factory *AssertionFactory) NewExistential(existential Existential) *dst.CallExpr {
	factory.variables.Push(existential.variables...)
	call := dstx.Call(
		dstx.SelectS("NewExistentialHyperAssertion").FromS(factory.packageName),
	).Pass(
		dstx.BasicInt(len(existential.variables)),
		factory.Create(existential.assertion),
	)
	factory.variables.PopN(len(existential.variables))
	
	return dstx.NewLineAround(call)
}
