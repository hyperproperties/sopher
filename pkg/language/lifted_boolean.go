package language

type LiftedBoolean uint8

const (
	LiftedFalse   = LiftedBoolean(0b0000)
	LiftedUnknown = LiftedBoolean(0b0001)
	LiftedTrue    = LiftedBoolean(0b0011)
)

func LiftBoolean(value bool) LiftedBoolean {
	if value {
		return LiftedTrue
	}
	return LiftedFalse
}

func (lhs LiftedBoolean) Or(rhs LiftedBoolean) LiftedBoolean {
	return LiftedBoolean(lhs | rhs)
}

func (lhs LiftedBoolean) And(rhs LiftedBoolean) LiftedBoolean {
	return LiftedBoolean(lhs & rhs)
}

func (lhs LiftedBoolean) Iff(rhs LiftedBoolean) LiftedBoolean {
	return LiftBoolean(lhs == rhs)
}

func (boolean LiftedBoolean) Not() LiftedBoolean {
	if boolean == LiftedUnknown {
		return LiftedUnknown
	}
	return boolean ^ LiftedTrue
}

func (boolean LiftedBoolean) IsFalse() bool {
	return boolean == LiftedFalse
}

func (boolean LiftedBoolean) IsUnknown() bool {
	return boolean == LiftedUnknown
}

func (boolean LiftedBoolean) IsTrue() bool {
	return boolean == LiftedTrue
}

func (boolean LiftedBoolean) String() string {
	if boolean == LiftedTrue {
		return "true"
	} else if boolean == LiftedFalse {
		return "false"
	}
	return "unknown"
}
