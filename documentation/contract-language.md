# Probabilitic Hyper-Assertions
Probabilistic Hyper-Assertions (PHA) are assertions that interact with a continuously growing set of sampled states, along with probabilitistic assignments, universal and exsistential quantification. These assertions can have one of three possible outcomes: _accepted_, _rejected_, or _inconclusive_. 
- An _accepted_ assertion means that the current set of states satisfies the assertion.
- A _rejected_ assertion means the states do not satisfy the assertion.
- If an assertion is _inconclusive_, it behaves like an accepted assertion—it does not report an error and allows the execution to continue.

```math
Θ = ∀π. Θ \; | \; ∃π. Θ \; | \; \neg Θ \; | \; Θ_1 \land Θ_2 \; | \; Φ_1 \geq Φ_2 \; | \; ψ^π
```
```math
Φ = P_{Π}(ψ^π) \; | \; P_{Π}(ψ_1^{π_1} \; | \; ψ_2^{π_2}) - Q
```
> In relaity a quantifier is required for $ψ^π$ to make sense. Otherwise, $π$ would be the empty set of free variables which would then just make $ψ$ constant. This is in the cases where $ψ$ does not rely on chaning state of the program which effects the outcome of $ψ$.

> The probability $P_{Π}(ψ^π)$ is the probability of choosing any random assignments of variables in $Π$ in such a way that $Θ$ is satisfied.  

> The probability $P_{Π}(ψ_1^{π_1} \; | \; ψ_2^{π_2})$ is the probability of choosing any random assignments of variables in $Π$ in such a way that $ψ_1^{π_1}$ is satisfied given $ψ_2^{π_2}$.  

> WE use $Q$ (rational numbers) because the probabilities are also fractions meaning subtraction does not reduce resolution.

## Sequential Probability Ratio Test
For the PHAs to work in practice the hypothesis testing must be done in sequence and not on a fixed sample set of states. To support this a Sequential Probability Ratio Test (SPRT) is applied. It allows for continuous monitoring of data and makes decisions about hypotheses as data is collected, rather than waiting until a predetermined sample size is reached. This also forces PHAs to have the option of returning _inconclusive_.

# Contract Language
The language used to formulate contracts in Golang is the following. Main aspects of it, is the ability to split the specification into region, quantification over sets of states, and probability function. In short, a contract is satisfied if an input falls into atleast one assumption which is either accepted or inconclusive. Only if the assumption accepts the state its is tested against the guarantee after execution. If the guarantee is not satisfied then we have found a violation of the contract. In addition, if the input is rejected by all assumptions then it is not accepted as a valid input.
- _Assumptions_: An input must satisfy at least one assumption (accepted or inconclusive) to be considered valid. For all the satisfied assumptions the guarantees must be satisfied.
- _Guarantees_: Only checked if the corresponding assumption accepts the state, ensuring that guarantees are contextually relevant.
- _Violations_: Clear pathways for detecting contract violations or rejecting invalid inputs based on assumption satisfaction.

```ebnf
Contract    = ( Obligations | Region ) { Region } .
Region      = "region" { Identifier }  "." Obligations .

Obligations = { ( Assumption | Guarantee ) } .
Assumption  = "assume" Assertion .
Guarantee   = "guarantee" Assertion .

Assertion   = Group | Expression | Assertion ⊕ Assertion | Probability ⧠ Probability .
Group       = "(" Quantifier ")" | Quantifier .
Quantifier  = ( "forall" | "exists" ) Variables "." Assertion .
Probability = "probability" "(" Expression ")" | "probability" "(" Expression "|" Expression  ")" [ ⋈ Number ] .
Expression  = GoExpression ";" .
Variables   = Identifier { Identifier } .
```
> `⊕ ∈ {&&, ||, ->, <->}`, `⧠ ∈ {<=, <, >, >=}`, `⋈ ∈ {+, -}`

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