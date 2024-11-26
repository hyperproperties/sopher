package main

import (
	"errors"
	"fmt"
	"log"
)

func main() {
	PrintPin(1, 2, 3, 4)
	PrintPin(123, 312)
	PrintPin(3, 1, 4, 1)
	PrintPin(3, 1, 4, 1)
}

func PrintPin(pin ...uint) {
	pass, err := NewPin(pin...).Check()
	if message := recover(); message != nil {
		log.Println("pin:", pin, "paniced with", message)
	} else {
		log.Println("pin:", pin, "passed?", fmt.Sprintf("%v,", pass), "err?", fmt.Sprintf("%v", err))
	}
}

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

func NewPin(pin ...uint) Pin {
	return pin
}

func (pin Pin) Valid() bool {
	return len(pin) == 4 &&
		pin[0] <= 9 && pin[1] <= 9 && pin[2] <= 9 && pin[3] <= 9
}

var ErrInvalidPin = errors.New("invalid pin")

var attempt int = 0

// assume: forall e. e.attempt >= 0																// Valid Attempt
// assume: forall e0. e0.attempt > 1; -> exists e1. e1.attempt == e0.attempt - 1				// Continous Attempts
// assume: forall e0 e1. e0._time < e1._time; <-> e0.attempt < e1.attempt						// Attempts Increment on Consecutive Calls
// assume: forall e0 e1. e0._id != e1._id; <-> e0.attempt != e1.attempt							// Unique Attempts
// guarantee: exists e. e.ret0 && e.ret1 == nil													// There Is A Check Which Passes
// guarantee: forall e. e.ret0; -> e.ret1 == nil												// All Passes does not produce an error
// guarantee: forall e. e.attempt > 3; -> !e.re0 && e.ret1 != nil								// Exceeds Attempt
// guarantee: forall e. !e.pin.Valid(); -> !e.ret0 && e.ret1 == nil								// Invalid Pin
// guarantee: forall e. e.ret0; <-> e.attempt <= 3 && SlicesEqual(e.pin, []uint{3, 1, 4, 1})	// Successful Check
// guarantee: forall e0 e1. e0.ret0 && e1.ret0; -> SlicesEqual(e0.pin, e1.pin)					// Exactly One Correct PIN
// guarantee: forall e0 e1. math.Abs(e0._duration - e1._duration) <= 0.1 * time.Second			// No Timing Side Channel
// assignable: attempt
func (pin Pin) Check() (bool, error) {
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
func (pin Pin) CheckV2() (bool, error) {
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
}
