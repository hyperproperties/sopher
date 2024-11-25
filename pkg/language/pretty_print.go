package language

import "strings"

func Print(ast Node) string {
	var builder strings.Builder

	var recursive func(ast Node)
	recursive = func(ast Node) {
		switch cast := ast.(type) {
		case Contract:
			for idx := range cast.regions {
				recursive(cast.regions[idx])
			}
		case Region:
			if builder.Len() > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString("region")
			for idx := range cast.name {
				builder.WriteString(" ")
				builder.WriteString(cast.name[idx])
			}
			builder.WriteString(": ")

			for idx := range cast.assumptions {
				recursive(cast.assumptions[idx])
			}

			for idx := range cast.guarantees {
				recursive(cast.guarantees[idx])
			}
		case Universal:
			builder.WriteString("forall")
			for idx := range cast.variables {
				builder.WriteString(" ")
				builder.WriteString(cast.variables[idx])
			}
			builder.WriteString(". ")
			recursive(cast.assertion)
		case Existential:
			builder.WriteString("exists")
			for idx := range cast.variables {
				builder.WriteString(" ")
				builder.WriteString(cast.variables[idx])
			}
			builder.WriteString(". ")
			recursive(cast.assertion)
		case Assumption:
			builder.WriteString("assume")
			builder.WriteString(": ")
			recursive(cast.assertion)
		case Guarantee:
			builder.WriteString("guarantee")
			builder.WriteString(": ")
			recursive(cast.assertion)
		case BinaryExpression:
			recursive(cast.lhs)
			switch cast.operator {
			case LogicalConjunction:
				builder.WriteString(" && ")
			case LogicalDisjunction:
				builder.WriteString(" || ")
			case LogicalImplication:
				builder.WriteString(" -> ")
			case LogicalBiimplication:
				builder.WriteString(" <-> ")
			default:
				panic("unknown binary operator")
			}
			recursive(cast.rhs)
		case UnaryExpression:
			switch cast.operator {
			case LogicalNegation:
				builder.WriteString("!")
			default:
				panic("unknown unary operator")
			}
			recursive(cast.operand)
		case GoExpresion:
			builder.WriteString(cast.code)
			builder.WriteString(";")
		case Group:
			builder.WriteRune('(')
			recursive(cast.node)
			builder.WriteRune(')')
		default:
			panic("unknown node")
		}
	}

	recursive(ast)

	return builder.String()
}
