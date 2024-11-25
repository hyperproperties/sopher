package language

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScope(t *testing.T) {
	scope := NewScope()

	depth := 0
	existentialSize := 0
	universalSize := 0

	for counter := 0; counter < 10000; counter++ {
		switch rand.IntN(5) {
		case 0: // Push a universal quantifier.
			oldDepth := scope.Depth()
			size := 1
			quantifier := NewUniversalQuantifierScope(size)
			scope.Push(quantifier)
			universalSize += size
			depth += 1
			assert.Equal(t, depth, oldDepth+1)
			assert.Equal(t, depth, scope.Depth())
			assert.False(t, scope.OnlyExistential())
		case 1: // Push an existential quantifier.
			oldDepth := scope.Depth()
			size := 1
			quantifier := NewExistentialQuantifierScope(size)
			scope.Push(quantifier)
			existentialSize += size
			depth += 1
			assert.Equal(t, depth, oldDepth+1)
			assert.Equal(t, depth, scope.Depth())
		case 2: // Peek the top of the stack.
			if scope.Depth() > 0 {
				oldDepth := scope.Depth()
				quantifier := scope.Pop()
				if quantifier.quantification == UniversalQuantification {
					universalSize -= quantifier.Size()
				}
				depth -= 1
				assert.Equal(t, oldDepth-1, scope.Depth())
			} else {
				assert.Panics(t, func() {
					scope.Pop()
				})
			}
		case 3:
			expectedOffset := 0
			for offset, quantifier := range scope.Quantifiers() {
				assert.Equal(t, expectedOffset, offset)
				expectedOffset += quantifier.Size()
			}
			assert.Equal(t, expectedOffset, scope.Size())
		case 4:
			for scope.Depth() > 0 {
				scope.Pop()
			}
			assert.Equal(t, 0, scope.Depth())
			assert.Equal(t, 0, scope.UniversalSize())
			assert.Equal(t, 0, scope.Size())
			depth = 0
			universalSize = 0
		case 5:
			assert.Equal(t, universalSize == 0 && existentialSize > 0, scope.OnlyExistential())
		}
	}
}
