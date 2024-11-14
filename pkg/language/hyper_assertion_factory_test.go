package language

import (
	"bytes"
	"go/printer"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPredicateMonitorCall(t *testing.T) {
	factory := NewGoMonitorFactory("sopher", "ExecutionModel")
	expression := NewGoExpression("!(e0.high == e1.high) || (e0.ret0 == e2.ret0)")
	call := factory.NewPredicateMonitorCall(expression)
	var buffer bytes.Buffer
	printer.Fprint(&buffer, token.NewFileSet(), call)
	assert.Equal(t, `sopher.NewPredicateMonitor(func(assignments []ExecutionModel) bool {
	return !(e0.high == e1.high) || (e0.ret0 == e2.ret0)
})`, buffer.String())
}

func TestNewUniversalMonitorCall(t *testing.T) {
	factory := NewGoMonitorFactory("sopher", "ExecutionModel")
	expression := NewGoExpression("!(e0.high == e1.high) || (e0.ret0 == e2.ret0)")
	forall := NewUniversal([]string{"e0", "e1"}, expression)
	call := factory.Create(forall)
	var buffer bytes.Buffer
	printer.Fprint(&buffer, token.NewFileSet(), call)
	assert.Equal(t, `sopher.NewUniversalMonitor(0, 2, sopher.NewPredicateMonitor(func(assignments []ExecutionModel) bool {
	e0, e1 := assignments[0], assignments[1]
	return !(e0.high == e1.high) || (e0.ret0 == e2.ret0)
}))`, buffer.String())
}

func TestNewExistentialMonitorCall(t *testing.T) {
	factory := NewGoMonitorFactory("sopher", "ExecutionModel")
	expression := NewGoExpression("e0.ret > 0")
	exists := NewExistential([]string{"e0"}, expression)
	call := factory.Create(exists)
	var buffer bytes.Buffer
	printer.Fprint(&buffer, token.NewFileSet(), call)
	assert.Equal(t, `sopher.NewExistentialMonitor(0, 1, sopher.NewPredicateMonitor(func(assignments []ExecutionModel) bool {
	e0 := assignments[0]
	return e0.ret > 0
}))`, buffer.String())
}

func TestNewUniversalaNDExistentialMonitorCall(t *testing.T) {
	factory := NewGoMonitorFactory("sopher", "ExecutionModel")
	expression := NewGoExpression("!(e0.high == e1.high) || (e0.ret0 == e2.ret0)")
	exists := NewExistential([]string{"e2"}, expression)
	forall := NewUniversal([]string{"e0", "e1"}, exists)
	call := factory.Create(forall)
	var buffer bytes.Buffer
	printer.Fprint(&buffer, token.NewFileSet(), call)
	t.Log(buffer.String())
	assert.Equal(t, `sopher.NewUniversalMonitor(0, 2, sopher.NewExistentialMonitor(2, 1, sopher.NewPredicateMonitor(func(assignments []ExecutionModel) bool {
	e0, e1, e2 := assignments[0], assignments[1], assignments[2]
	return !(e0.high == e1.high) || (e0.ret0 == e2.ret0)
})))`, buffer.String())
}
