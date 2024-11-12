package language

import (
	"bytes"
	"go/printer"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsddsa(t *testing.T) {
	assumption := NewUniversal(
		[]string{"e0"}, NewGoExpression("e0.value < 0"),
	)
	guarantee := NewUniversal(
		[]string{"e0"}, NewGoExpression("e0.ret0 > 0"),
	)

	factory := NewGoAGContractFactory("sopher", "RetainModel")
	contract := factory.Constructor([]Node{assumption}, []Node{guarantee})

	var buffer bytes.Buffer
	printer.Fprint(&buffer, token.NewFileSet(), contract)

	assert.Equal(t, `NewAGContract.sopher([]RetainModel{sopher.NewUniversalMonitor(0, 1, sopher.NewPredicateMonitor(func(assignments []RetainModel) bool {
	e0 := assignments[0]
	return e0.value < 0
}))}, []RetainModel{sopher.NewUniversalMonitor(1, 1, sopher.NewPredicateMonitor(func(assignments []RetainModel) bool {
	e0, e0 := assignments[0], assignments[1]
	return e0.ret0 > 0
}))})`, buffer.String())
}
