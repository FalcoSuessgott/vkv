package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidInteger(t *testing.T) {
	i, ok := ToInteger("2")
	assert.Equal(t, 2, i)
	assert.Nil(t, ok, "expect error to be nil")
}

func TestInValidIntegers(t *testing.T) {
	_, ok := ToInteger("string")
	assert.Error(t, ok)
}
