# Probabilitic Hyper-Assertions
Probabilistic Hyper-Assertions (PHA) are assertions that interact with a continuously growing set of sampled states, along with probabilitistic assignments, universal and exsistential quantification. These assertions can have one of three possible outcomes: _accepted_, _rejected_, or _inconclusive_. 
- An _accepted_ assertion means that the current set of states satisfies the assertion.
- A _rejected_ assertion means the states do not satisfy the assertion.
- If an assertion is _inconclusive_, it behaves like an accepted assertion—it does not report an error and allows the execution to continue.

```math
Θ = ∀v. Θ \; | \; ∃v. Θ \; | \; \neg Θ \; | \; Θ_1 \land Θ_2 \; | \; Φ_1 \geq Φ_2 \; | \; ψ
```
```math
Φ = Φ - C \; | \; P_{Π}(Θ) \; | \; P_{Π}(Θ_1 \; | \; Θ_2)
```
> The probability $P_{Π}(Θ)$ is the probability of choosing any random assignments of variables in $Π$ in such a way that $Θ$ is satisfied.  

> The probability $P_{Π}(Θ_1 | Θ_2)$ is the probability of choosing any random assignments of variables in $Π$ in such a way that $Θ_1$ is satisfied given $Θ_2$.  

> Where $ψ$ is a valid boolean expression in Go. over the assignments in its scope and $C∈Q$ (It is a rational because the probabilities are also fractions meaning subtraction does not reduce resolution).

## Sequential Probability Ratio Test
For the PHAs to work in practice the hypothesis testing must be done in sequence and not on a fixed sample set of states. To support this a Sequential Probability Ratio Test (SPRT) is applied. It allows for continuous monitoring of data and makes decisions about hypotheses as data is collected, rather than waiting until a predetermined sample size is reached. This also forces PHAs to have the option of returning _inconclusive_.

# Contract Language
The language used to formulate contracts in Golang is the following. Main aspects of it, is the ability to split the specification into region, quantification over sets of states, and probability function. In short, a contract is satisfied if an input falls into atleast one assumption which is either accepted or inconclusive. Only if the assumption accepts the state its is tested against the guarantee after execution. If the guarantee is not satisfied then we have found a violation of the contract. In addition, if the input is rejected by all assumptions then it is not accepted as a valid input.
- _Assumption Handling_: An input must satisfy at least one assumption (accepted or inconclusive) to be considered valid.
- _Guarantee Verification_: Guarantees are only checked if the corresponding assumption accepts the state, ensuring that guarantees are contextually relevant.
- _Violation and Rejection Handling_: Clear pathways for detecting contract violations or rejecting invalid inputs based on assumption satisfaction.

```ebnf
Contract     = ( Obligations | Region ) { Region } .
Name         = Identifier { Identifier } .
Region       = "region" ":" [ Name ] Obligations .
Obligations  = Quantifier { ( Assumption | Guarantee ) } [ Obligations ] .
Quantifier   = Universal | Exsistential .
Universal    = "forall" ":" Variables .
Exsistential = "exists" ":" Variables .
Variables    = Variable { Variable } .
Variable     = Identifier .

Assumption   = "assume" ":" Expression .
Guarantee    = "guarantee" ":" Expression .

Assertion    = "!" Assertion | Assertion ⊕ Assertion | Expression.
Probability  = "probability" "(" "{" Variables "}" "," Expression [ "," Expression ] ")" .
Expression   = ... .
```
> `⊕` is the binary operators which are allowed `⊕ ∈ {<, <=, >, >=, &&, ||, ->, <->}` where implication and biimplication is specifially a contractual operator and not a part of Go.  

- _Assumption:_ Probabilistic hyper-assertions on state excluding time and return values.  
- _Guarantee:_ Probabilistic hyper-assertions on state including time and return values.

## Examples
There are multiple example for very common cases of hyperproperties. To highlight the usage of hypercontracts with both probabilities and time-sensitivity of the execution in case where service level agreements are important.

If the request is by an admin then if the response length is less than 100 bytes then more than half of the response times will be less than 0.2 seconds. However, if is not an admin, then we have no guarantee on anything.
```go
// region: Admin
// forall: e
// assume: e.admin
// guarantee: probability({t}, t.time <= 0.2, len(t.ret0) < 100) >= 0.5
// region: User
// forall: e
// assume: !e.admin
func Request(admin bool) []byte { ... }
```