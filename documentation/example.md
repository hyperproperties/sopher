# Pin Checker
_Assumptions:_
- _Valid PIN and Attempt:_ The PIN is valid, and the number of attempts is greater than 0.
- _Continous Attempts:_ If there is an execution which is not the initial, then there must be one before it.
- __With Reset__:
  - _Reset After 1 Minute:_ The attempts can be reset to 1 if there is atleast one minute between them.
- __Without Reset__:
  - _Unique Attempts:_ The attempt must be unique between all executions.
  - _Attempts Increment on Consecutive Calls:_ The attempt must be strictly increasing by one each execution over time.

_Guarantees:_
- _Successful Check:_ The PIN is valid, and the correct pin is [3, 1, 4, 1].
- _Exactly One Correct PIN:_ The two executions are both correct then they must have the same PIN.
- _No Timing Side Channel:_ There must be at most a difference of 100ms between every execution durations.

```go
func CheckPin(attempt uint, pin Pin) bool { ... }
```

# Pin Registration
Specification:
- The function returns true if the pin has been set. Otherwise, false.
- The pin is four numbers in [0, 9].
- The initial pin must always be 0000.
- The same pin cannot be used consequtively.
- The execution time difference between two registrations must not exceed 10ms (Ensures that registering the same pin as the currently correct one does not show).
- No pin must be the reversal of another.
- TODO: Pin changes must be atleast 15 minute after another.

```go
// assume: forall e. e.digits[0] <= 9 && e.digits[1] <= 9 && e.digits[2] <= 9 && e.digits[3] <= 9
// assume: forall e0 e1. e0.digitgs[0] != e1.digits[0] && e0.digitgs[1] != e1.digits[1] && e0.digitgs[2] != e1.digits[2] && e0.digitgs[3] != e1.digits[3]
func RegisterPin(digits [4]uint) bool { ... }
``` 