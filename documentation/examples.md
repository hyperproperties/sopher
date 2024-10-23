# Non-interference
_Non-interference_: A program satisfies noninterference when the outputs observed by low-security users are the same as they would be in the absence of inputs submitted by high-security users. 

```math
∀π.∀π'. π =^H_{in} π' → π =^L_{out} π'
```

```go
// guarantee: forall t0 t1. (t0.high == t1.high) -> (t0.ret == t1.ret)
func Foo(low, high int) int { ... }
```

or

```go
// guarantee: forall t0 t1. !(t0.high == t1.high) || t0.ret == t1.ret
func Foo(low, high int) int { ... }
```

A probabilistic interpretation of the hyperproperty could be that, for any two random assignments, there is a certain probability non-interference will be satisfied.

```math
P_{π, π'}(π =^L_{in} π' | π =^L_{out} π') > 0.99
```

> The probability of choosing random assignments to `π` and `π'` which satisfies the non-interference hyperproperty is atleast 99%. Allowing around 1% of all random assignments to break non-interference.

```go
// guarantee: forall t0 t1. probability t2. t2.high = t0.high && t2.ret == t1.ret; > 0.99
func Foo(low, high int) int { ... }
```


# Generalized Non-interference
_Generalized noninterference_: Allows for non-determinism in low-observable behavior while ensuring that low-security outputs remain unchanged in response to high-security inputs. This can be seen as the same high inputs only ones has to have the same low return value. Therfore, it in some way, relaxes the non-interference requirement and allows non-determinism.

```math
∀π.∀π'.∃π''. π'' =^H_{in} π ∧ π'' =^L_{out} π'
```

```go
// guarantee: forall t0 t1. exists t2. t2.high = t0.high && t2.ret == t1.ret
func Foo(low, high int) int {
    ...
}
```

# Observational Determinism
_Observational determinism_: A non-deterministic program satisfies observational determinism if every pair with the same low inputs remain indistinguishable for low users. That is, the program appears to be deterministic to low users.

```math
∀π.∀π'. π =^L_{in} π' → π =^L_{out} π'
```

```go
// guarantee: forall t0 t1. (t0.low == t1.low) -> (t0.ret == t2.ret)
func Foo(low, high int) int { ... }
```

or

```go
// guarantee: forall t0 t1. !(t0.low == t1.low) || t0.ret == t2.ret
func Foo(low, high int) int { ... }
```

A probabilistic interpretation of this hyperproperty could be that for all `π` the probability of randomly choosing a `π'` where observational determinism is satisfied is atleast 80%.

```math
∀π. P_{π'}(π =^L_{in} π' | π =^L_{out} π') > 0.8
```

```go
// guarantee: foraal t0. probability t1. t0.low == t1.low; | t0.ret == t1.ret; > 80%
func Foo(low, high int) int { ... }
```

# Declassification
_Declassification_: Some programs need to reveal secret information to fulfill functional requirements. For example, a password checker must reveal whether the entered password is correct or not. If the low and declassification inputs are the same the output is the same.

```math
∀π.∀π'.(π =^L_{in} π' ∧ π =^D_{in} π') → π =^L_{out} π'
```

```go
// guarantee: forall t0 t1. (t0.user == t1.user && t0.password == t1.password) -> (t0.ret == t1.ret)
func Authenticate(user, password string) bool { ... }
```

or

```go
// guarantee: forall t0 t1. !(t0.user == t1.user && t0.password == t1.password) || t0.ret == t1.ret
func Authenticate(user, password string) bool { ... }
```

A probabilistic interpretation of declassification could be an example where the probability of declassification is correlated the age of the data to determine declassification with.

```math
∀π.P_{π'}(π =^L_{in} π' ∧ π =^D_{in} π' | π =^L_{out} π') > 0.8
```

# Maximum Mean Response Time
_Maximum Mean Response Time_: This is a common type of service level agreement (SLA) where a service is required to respond, on average, within a specified time limit (upper bound). Unlike traditional properties, which describe individual system behaviors, hyperproperties allow reasoning over the mean response time across multiple executions of the system. This enables probabilistic analysis of response times and similar performance metrics. With probabilities we relax the mean to a probability which is sufficient since we rarely what an equality check of a SLA.

Atleast 50% of all responses does not exceed a response time of 0.5 seconds.
```go
// guarantee: probability t. t.time <= 0.5; >= 0.5
func Request() []byte { ... }
```

Atleast 95% of all responses does not exceed a response time of 0.1 seconds.  
```go
// guarantee: probability t. t.time <= 0.1; >= 0.95
func Request() []byte { ... }
```

The slowest 5% does not exceed a response time of 1 second.
```go
// guarantee: probability t. t.time <= 1; <= 0.05
func Request() []byte { ... }
```

It is more likely to get a response in 0.1 seconds than more than 2 seconds.
```go
// guarantee: probability t. t.time <= 0.1; > probability t. t.time > 2
func Request() []byte { ... }
```

If the response length is less than 100 bytes then more than half of the response times will be less than 0.2 seconds.
```go
// guarantee: probability t. t.time <= 0.2; | len(t.ret0) < 100; >= 0.5
func Request() []byte { ... }
```

# Time-Sensitive Side-Channel
_Time-Sensitive Side-Channel_: Like non-interference we dont want information of confidential material leaked though a low observable channel. In some case we want some information leaked as in the case of declassification - we want to be able to tell whether the password was correct or not. However, in the case of incorrect passwords, we dont want to leak how correct the incorrect password was and at what character of the password did the password become incorrect. A naive password checker would act like a string compare and return a result immediately when the password was found to be incorrect. However, in the case of very long passwords and where the time of comparing is measureable so would the time to reach the result of a incorrect password also be measureable. This highlights one case where time is necessary for a secure information flow policy (hyperproperty).

```go
// guarantee: forall e0 e1. e0.time - e1.time < time.Second
func Authenticate(user, password string) bool { ... }
```

There cannot be more than 1 second difference between any pair of executions no matter whether they were correctly authenticated or not.

# Erasure (TODO: Maybe requires custom composition operation to be done well?)
_Erasure_: Refers to the process of completely and irretrievably deleting data or information from a storage medium to prevent its recovery or access. This process often involves overwriting the original data with random values or zeros, ensuring that any remnants of the original content cannot be reconstructed. Effective erasure is crucial for protecting sensitive information and maintaining data privacy in various applications, including personal computing and enterprise data management.

```go
func Read(path string) ([]byte, bool) { ... }

// compose: r Read, e Erase
// guarantee forall t0 t1. (t0.r.path == t0.r.path && ) 
func Erase(path string) { ... }
```