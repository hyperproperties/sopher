# Non-interference
_Non-interference_: A program satisfies noninterference when the outputs observed by low-security users are the same as they would be in the absence of inputs submitted by high-security users. 

```math
∀π.∀π'. π =^H_{in} π' → π =^L_{out} π'
```

```go
// forall: t0 t1
// assume: t0.high == t1.high
// guarantee: t0.ret == t2.ret
func Foo(low, high int) int { ... }
```

or

```go
// forall: t0 t1
// guarantee: !(t0.high == t1.high) || t0.ret == t2.ret
func Foo(low, high int) int { ... }
```

The former allows an easier generation of paired execution using the `assume` which allows the generation of high values in such a way that any pair following the assume is relevant for the guarantee.

A probabilistic interpretation of the hyperproperty could be that for any two random assignments there is a probability of non-interference to be satisfied some percentage of the time. However, we will allow the property to be broken some times.

```math
P_{π, π'}(π =^L_{in} π' | π =^H_{out} π') > 0.99
```

> The probability of choosing random assignments to `π` and `π'` which satisfies the non-interference hyperproperty is 99%. Allowing 1% of all random assignments to break non-interference.

```go
// forall: t0 t1
// guarantee: probability({t2}, t2.high = t0.high && t2.ret == t1.ret) > 0.99
func Foo(low, high int) int { ... }
```


# Generalized Non-interference
_Generalized noninterference_: Allows for non-determinism in low-observable behavior while ensuring that low-security outputs remain unchanged in response to high-security inputs. This can be seen as the same high inputs only ones has to have the same low return value. Therfore, it in some way, relaxes the non-interference requirement and allows non-determinism.

```math
∀π.∀π'.∃π''. π'' =^H_{in} π ∧ π'' =^L_{out} π'
```

```go
// forall: t0 t1
// exists: t2
// guarantee: t2.high = t0.high && t2.ret == t1.ret
func Foo(low, high int) int {
    ...
}
```

A probabilistic interpretation of this hyperproperty could be that for any random assignment to `π''` some percentage of them inhibits a non-interference behaviour.

```math
∀π.∀π'. P_{π''}(π'' =^H_{in} π | π'' =^L_{out} π') > 0.95
```

> For all pairs of  `π` and  `π'` the probability of choosing a random assignment to `π''` which satisfies the hyperproperty is 95%.

```go
// forall: t0 t1
// guarantee: probability({t2}, t2.high = t0.high, t2.ret == t1.ret) > 95%
func Foo(low, high int) int {
    ...
}
```

# Observational Determinism
_Observational determinism_: A (nondeterministic) program satisfies observational determinism if every pair with the same low inputs remain indistinguishable for low users. That is, the program appears to be deterministic to low users.

```math
∀π.∀π'. π =^L_{in} π' → π =^L_{out} π'
```

```go
// forall: t0 t1
// assume: t0.low == t1.low
// guarantee: t0.ret == t2.ret
func Foo(low, high int) int { ... }
```

or

```go
// forall: t0 t1
// guarantee: !(t0.low == t1.low) || t0.ret == t2.ret
func Foo(low, high int) int { ... }
```

The former allows an easier generation of paired execution using the `assume` which allows the generation of high values in such a way that any pair following the assume is relevant for the guarantee.

A probabilistic interpretation of this hyperproperty could be that for all `π` the probability of randomly choosing a `π'` where observational determinism is satisfied is 80%. Meaning that the chance of any execution to be observational deterministic is atleast 80%.

```math
P_{π, π'}(π =^L_{in} π' | π =^L_{out} π') > 0.8
```

> for all executions a random assignment to `π'` has at least 80% chance of satisfying observational determinism.

```go
// guarantee: probability({t0, t1}, t0.low == t1.low, t0.ret == t1.ret) > 80%
func Foo(low, high int) int { ... }
```

# Declassification
_Declassification_: Some programs need to reveal secret information to fulfill functional requirements. For example, a password checker must reveal whether the entered password is correct or not. If the low and declassification inputs are the same the output is the same.

```math
∀π.∀π'.(π =^L_{in} π' ∧ π =^D_{in} π') → π =^L_{out} π'
```

```go
// forall: t0 t1
// guarantee: !(t0.user == t1.user && t0.password == t1.password) || t0.ret == t1.ret
func Authenticate(user, password string) bool { ... }
```

or

```go
// forall: t0 t1
// assume: t0.user == t1.user && t0.password == t1.password
// guarantee: t0.ret == t1.ret
func Authenticate(user, password string) bool { ... }
```

A probabilistic interpretation of declassification could be an example where the probability of declassification is correlated the age of the data to determine declassification with.

```math
∀π.P_{π'}(π =^L_{in} π' ∧ π =^D_{in} π' | π =^L_{out} π') > 0.8
```

> for all executions a random assignment to `π'` has at least 80% chance of satisfying declassification. 

__I cannot construct a meaningful example of probabilistic declassification__

# Maximum Mean Response Time
_Maximum Mean Response Time_: This is a common type of service level agreement (SLA) where a service is required to respond, on average, within a specified time limit (upper bound). Unlike traditional properties, which describe individual system behaviors, hyperproperties allow reasoning over the mean response time across multiple executions of the system. This enables probabilistic analysis of response times and similar performance metrics. With probabilities we relax the mean to a probability which is sufficient since we rarely what an equality check of a SLA.

Provided som condition determined by `C(π)` the probability `P_{π}(C(π))` determines the ratio of which the condition is true, and `⋈` is then the inequality comparison operator.

```math
P_{π}(C(π)) ⋈ p
```

Atleast 50% of all responses does not exceed a response time of 0.5 seconds.
```go
// guarantee: probability({t}, t.time <= 0.5) >= 0.5
func Request() []byte { ... }
```

Atleast 95% of all responses does not exceed a response time of 0.1 seconds.  
```go
// guarantee: probability({t}, t.time <= 0.1) >= 0.95
func Request() []byte { ... }
```

The slowest 5% does not exceed a response time of 1 second.
```go
// guarantee: probability({t}, t.time <= 1) <= 0.05
func Request() []byte { ... }
```

It is more likely to get a response in 0.1 seconds than more than 2 seconds.
```go
// guarantee: probability({t}, t.time <= 0.1) > probability({t}, t.time > 2)
func Request() []byte { ... }
```

If the response length is less than 100 bytes then more than half of the response times will be less than 0.2 seconds.
```go
// guarantee: probability({t}, t.time <= 0.2, len(t.ret0) < 100) >= 0.5
func Request() []byte { ... }
```

# Erasure (TODO: Maybe requires custom composition operation to be done well?)
_Erasure_: Refers to the process of completely and irretrievably deleting data or information from a storage medium to prevent its recovery or access. This process often involves overwriting the original data with random values or zeros, ensuring that any remnants of the original content cannot be reconstructed. Effective erasure is crucial for protecting sensitive information and maintaining data privacy in various applications, including personal computing and enterprise data management.

```go
func Read(path string) ([]byte, bool) { ... }

// compose: r Read, e Erase
// forall: t0 t1
// guarantee: (t0.r.path == t0.r.path && ) 
func Erase(path string) { ... }
```