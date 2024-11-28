package examples

func SlicesEqual[S1, S2 ~[]E, E comparable](s1 S1, s2 S2) bool {
	if len(s1) != len(s2) {
		return false
	}

	for idx := range s1 {
		elem1, elem2 := s1[idx], s2[idx]
		if elem1 != elem2 {
			return false
		}
	}

	return true
}

type Pin []uint

func (pin Pin) Valid() bool {
	return len(pin) == 4 &&
		pin[0] <= 9 && pin[1] <= 9 && pin[2] <= 9 && pin[3] <= 9
}

// assume: forall e. pin.Valid() && e.attempt > 0												// Valid PIN and Attempt
// assume: forall e0. exists e1. e0.attempt > 1; -> e1.attempt == e0.attempt - 1				// Continous Attempts
// assume: forall e0 e1. e0._time < e1._time; <-> e0.attempt < e1.attempt						// Attempts Increment on Consecutive Calls
// assume: forall e0 e1. e0._id != e1._id; <-> e0.attempt != e1.attempt							// Unique Attempts
// guarantee: forall e. e.ret0; <-> e.attempt <= 3 && SlicesEqual(e.pin, []uint{3, 1, 4, 1})	// Successful Check
// guarantee: forall e0 e1. e0.ret0 && e1.ret0; -> SlicesEqual(e0.pin, e1.pin)					// Exactly One Correct PIN
// guarantee: forall e0 e1. math.Abs(e0._duration - e1._duration) <= 0.1 * time.Second			// No Timing Side Channel
func CheckPIN(attempt uint, pin Pin) bool {
	panic("not implemented yet")
}

// assume: forall e. pin.Valid() && e.attempt > 0												// Valid PIN and Attempt
// assume: forall e0. exists e1. e0.attempt > 1; -> e1.attempt == e0.attempt - 1				// Continous Attempts
// assume: forall e0 e1. e0._time + time.Minute <= e1._time; -> e1.attempt == 1					// Reset After 1 Minute
// guarantee: forall e. e.ret0; <-> e.attempt <= 3 && SlicesEqual(e.pin, []uint{3, 1, 4, 1})	// Successful Check
// guarantee: forall e0 e1. e0.ret0 && e1.ret0; -> SlicesEqual(e0.pin, e1.pin)					// Exactly One Correct PIN
// guarantee: forall e0 e1. math.Abs(e0._duration - e1._duration) <= 0.1 * time.Second			// No Timing Side Channel
func WithResetCheckPIN(attempt uint, pin Pin) bool {
	panic("not implemented yet")
}

// assume: forall e. pin.Valid() && e.attempt > 0										// Valid PIN and Attempt
// assume: forall e0. exists e1. e0.attempt > 1; -> e1.attempt == e0.attempt - 1		// Continous Attempts
// assume: forall e0 e1. e0._time < e1._time; <-> e0.attempt < e1.attempt				// Attempts Increment on Consecutive Calls
// assume: forall e0 e1. e0._id != e1._id; <-> e0.attempt != e1.attempt					// Unique Attempts
// guarantee: forall e0 e1. e0.ret0 && e1.ret0; -> SlicesEqual(e0.pin, e1.pin)			// Exactly One Correct PIN
// guarantee: forall e0 e1. math.Abs(e0._duration - e1._duration) <= 0.1 * time.Second	// No Timing Side Channel
func CheckUnknownPIN(attempt uint, pin Pin) bool {
	panic("not implemented yet")
}

// assume: forall e. pin.Valid()														// Valid PIN
// guarantee: forall e0 e1. math.Abs(e0._duration - e1._duration) <= 0.1 * time.Second	// No Time Side Channel
func CheckContinouslyChangingPIN(pin Pin) bool {
	panic("not implemented yet")
}

// assume: forall e. pin.Valid()														// Valid PIN
// assume: forall e0. exists e1. e0.counter > 0; -> e1.counter == e0.counter - 1		// Continous Counter
// assume: forall e0 e1. e1.counter == e0.counter + 1; -> !SlicesEqual(e0.pin, e1.pin)	// No Immediate Duplicate
// guarantee: forall e0. e0.counter == 0; -> SlicesEqual(e0.pin, []uint{0, 0, 0, 0})	// Inital PIN 0000
// guarantee: forall e0 e1. math.Abs(e0._duration - e1._duration) <= 0.1 * time.Second	// No Timing Side Channel
// guarantee: "TODO: No pin must be the reversal of another."							//
// guarantee: forall e0. e0.ret0 -> forall e1. e1.ret0 && e1._time < e0._time;
//
//	-> e1._time + 15 * time.Minute <= e0._time			// Atleast 15 Minute Between Successful Changes
func ChangePIN(counter uint, pin Pin) bool {
	panic("not implemented yet")
}

// contract:
//
//	assume: forall e. e.attempt >= 0
//	assume: forall e0. e0.attempt > 1; -> exists e1. e1.attempt == e0.attempt - 1
//	assume: forall e0 e1. e0._time < e1._time; <-> e0.attempt < e1.attempt
//	assume: forall e0 e1. e0._id != e1._id; <-> e0.attempt != e1.attempt
//	guarantee: exists e. e.ret0 && e.ret1 == nil
//	guarantee: forall e. e.ret0; -> e.ret1 == nil
//	guarantee: forall e. e.attempt > 3; -> !e.re0 && e.ret1 != nil
//	guarantee: forall e. !e.pin.Valid(); -> !e.ret0 && e.ret1 == nil
//	guarantee: forall e. e.ret0; <-> e.attempt <= 3 && SlicesEqual(e.pin, []uint{3, 1, 4, 1})
//	guarantee: forall e0 e1. e0.ret0 && e1.ret0; -> SlicesEqual(e0.pin, e1.pin)
//	guarantee: forall e0 e1. math.Abs(e0._duration - e1._duration) <= 0.1 * time.Second
//	assignable: attempt
/*func (pin Pin) CheckV2() (bool, error) {
	if attempt < 0 {
		panic("a negative attempt cannot exist")
	}

	attempt++

	if attempt > 3 {
		return false, nil
	}

	if !pin.Valid() {
		return false, ErrInvalidPin
	}

	return SlicesEqual(pin, []uint{3, 1, 4, 1}), nil
}*/
