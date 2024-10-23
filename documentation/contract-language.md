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

> The probability $P_{Π}(ψ^π)$ is the probability of choosing any random assignments of variables in $Π$ in such a way that $ψ^π$ is satisfied with the assignments $π$.  

> The probability $P_{Π}(ψ_1^{π_1} \; | \; ψ_2^{π_2})$ is the probability of choosing any random assignments of variables in $Π$ in such a way that given $ψ_2^{π_2}$, $ψ_1^{π_1}$ is satisfied.  

> Rational numbers $Q$ are used because the probabilities are also fractions meaning subtraction does not reduce resolution.

# Contract Language
Main aspects of the contract structure, is probabilistic hyper assertions and the ability to split the specification into regions. An execution falls into a region if the region's assumption is either accepted or inconclusive. The contract is said to be breached, if for any regions an execution falls into, the guarantee rejects it.

```ebnf
Contract    = ( Obligations | Region ) { Region } .
Region      = "region" { Identifier }  "." Obligations .

Obligations = { ( Assumption | Guarantee ) } .
Assumption  = "assume" ":" Assertion .
Guarantee   = "guarantee" ":" Assertion .

Assertion   = Group | Expression | Assertion ⊕ Assertion | Probability ⧠ Probability .
Group       = "(" Quantifier ")" | Quantifier .
Quantifier  = ( "forall" | "exists" ) Variables "." Assertion .
Probability = "probability" Variables "." Expression | "probability" Variables "." Expression "|" Expression [ ⋈ Number ] .
Expression  = GoExpression ";" .
Variables   = Identifier { Identifier } .
```
> `⊕ ∈ {&&, ||, ->, <->}`, `⧠ ∈ {<=, <, >, >=}`, `⋈ ∈ {+, -}`

## Sequential Probability Ratio Test
For the PHAs to work in practice the hypothesis testing must be done in sequence and not on a fixed sample set of states. To support this a Sequential Probability Ratio Test (SPRT) is applied. It allows for continuous monitoring of data and makes decisions about hypotheses as data is collected, rather than waiting until a predetermined sample size is reached. This also forces PHAs to have the option of returning _inconclusive_.