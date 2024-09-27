# Non-interference
_Non-interference_: A program satisfies noninterference when the outputs observed by low-security users are the same as they would be in the absence of inputs submitted by high-security users. The following HyperLTL formula proposed in "Temporal Logics for Hyperproperties" by Clarkson describes non-interference as a 2-hypersafety property universally quantifying over pairs of traces:

```math
∀π.∀π′. π[0] =^L_{in} π′[0] → π =^H_{out} π′
```

In Sopher traces are instead of a model of an execution which can be extended arbritarily as long as it has a record of the inputs and outputs of the function. E.g., temporal aspects can be added to cusom execution models. The following is an example of a pure execution model using the HyperLTL language for contractual specificaiton non-interference for a function in Go.

```go
// forall: t1
// guarantee: ret == low
// forall: t2
// assume: low' == low
// guarantee: t1.ret == ret
func Retain(low, high int) int {
    if high % 2 == 0 {
        return 0
    }
    return low
}
```

The contract is converted to a contract object which is stored in the global scope but not exported. This enables contracts to store previous executions such that when checking hyperproperties they cna be used without having to execute entirely new ones. In addition, maybea heuristic can be used to prioritise what traces to compare initially.
```go
var retainWrap := func(low, high int) int {
    if high % 2 == 0 {
        return 0
    }
    return low
}

var retainContract := NewContract(retainWrap, ...)

func Retain(low, high int) int {
    returns := retainContract.Check(low, high)
    return returns[0]
}
```