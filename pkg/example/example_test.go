package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	assert.Equal(t, 5, Add(2, 3))
}

func TestMultiply(t *testing.T) {
	assert.Equal(t, 6, Multiply(2, 3))
}
