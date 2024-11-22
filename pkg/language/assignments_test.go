package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssignmentsExpand(t *testing.T) {
	assignments := NewAssignments[int]()

	from1, to1 := assignments.Expand(2)
	assert.Equal(t, 2, to1-from1)
	assert.Equal(t, 2, len(assignments))

	from2, to2 := assignments.Expand(3)
	assert.Equal(t, 3, to2-from2)
	assert.Equal(t, 5, len(assignments))

	assert.Panics(t, func() {
		assignments.Expand(-1)
	})

	buffer1 := assignments[from1:to1]
	buffer1[0], buffer1[1] = -1, 3
	assert.ElementsMatch(t, assignments[:2], buffer1)

	buffer2 := assignments[from2:to2]
	buffer2[0], buffer2[1], buffer2[2] = 3, 9, 312
	assert.ElementsMatch(t, assignments[2:], buffer2)
}

func TestAssignmentsShrink(t *testing.T) {
	assignments := NewAssignments(0, 1, 2, 3, 4, 5, 6)
	assert.Len(t, assignments, 7)

	assert.Panics(t, func() {
		assignments.Shrink(-1)
	})

	assignments.Shrink(1)
	assert.Len(t, assignments, 6)

	assignments.Shrink(4)
	assert.Len(t, assignments, 2)

	assignments.Shrink(2)
	assert.Len(t, assignments, 0)
	assert.Panics(t, func() {
		assignments.Shrink(1)
	})
}
