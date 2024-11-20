package language

// A logical type which can take three values: true (T), unknown (U), and false (F).
// The implemented logic follows strong logic of indeterminacy.
type LiftedBoolean uint8

const (
	LiftedFalse   = LiftedBoolean(0b0000)
	LiftedUnknown = LiftedBoolean(0b0001)
	LiftedTrue    = LiftedBoolean(0b0011)
)

// Lifts a boolean value to a lifted boolean which can then be used for three value logic.
// If true it returns LiftedTrue. Otherwise, returns LiftedFalse.
func LiftBoolean(value bool) LiftedBoolean {
	if value {
		return LiftedTrue
	}
	return LiftedFalse
}

// Performs logical disjunction on two lifted booleans.
// The operation if symmetrical such that OR(x, y) = OR(y, x)
//
//	OR(F, F) = F
//	OR(F, U) = U
//	OR(F, T) = T
//	OR(U, U) = U
//	OR(U, T) = T
//	OR(T, T) = T
func (lhs LiftedBoolean) Or(rhs LiftedBoolean) LiftedBoolean {
	return LiftedBoolean(lhs | rhs)
}

// Performs logical conjunction on two lifted booleans.
// The operation if symmetrical such that AND(x, y) = AND(y, x)
//
//	AND(F, F) = F
//	AND(F, U) = F
//	AND(F, T) = F
//	AND(U, U) = U
//	AND(U, T) = U
//	AND(T, T) = T
func (lhs LiftedBoolean) And(rhs LiftedBoolean) LiftedBoolean {
	return LiftedBoolean(lhs & rhs)
}

// Returns true if and only if the two lifted booleans have the same value.
//
//	IFF(T, T) = T
//	IFF(T, U) = U
//	IFF(T, F) = F
//	IFF(U, U) = U
//	IFF(F, T) = F
//	IFF(F, U) = U
//	IFF(F, F) = T
func (lhs LiftedBoolean) Iff(rhs LiftedBoolean) LiftedBoolean {
	return lhs.If(rhs).And(rhs.If(lhs))
}

// Returns the implication of p -> q.
//
//	IF(T, T) = T
//	IF(T, U) = U
//	IF(T, F) = F
//	IF(U, T) = T
//	IF(U, U) = U
//	IF(U, F) = U
//	IF(F, T) = T
//	IF(F, U) = T
//	IF(F, F) = T
func (lhs LiftedBoolean) If(rhs LiftedBoolean) LiftedBoolean {
	return lhs.Not().Or(rhs)
}

// Negates the lifted boolean.
//
//	NOT(F) = T
//	NOT(U) = U
//	NOT(T) = T
func (boolean LiftedBoolean) Not() LiftedBoolean {
	if boolean == LiftedUnknown {
		return LiftedUnknown
	}
	return boolean ^ LiftedTrue
}

// Returns true if the lifted boolean is LiftedFalse.
func (boolean LiftedBoolean) IsFalse() bool {
	return boolean == LiftedFalse
}

// Returns true if the lifted boolean is LiftedUnknown.
func (boolean LiftedBoolean) IsUnknown() bool {
	return boolean == LiftedUnknown
}

// Returns true if the lifted boolean is LiftedTrue.
func (boolean LiftedBoolean) IsTrue() bool {
	return boolean == LiftedTrue
}

// Returns the stringifyed version of the lifted boolean.
func (boolean LiftedBoolean) String() string {
	if boolean == LiftedTrue {
		return "true"
	} else if boolean == LiftedFalse {
		return "false"
	}
	return "unknown"
}
