package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskSecrets(t *testing.T) {
	testCases := []struct {
		name    string
		options []Option
		input   map[string]interface{}
		output  map[string]interface{}
	}{
		{
			name:    "test: normal secrets",
			options: nil,
			input: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "value", "user": "password"},
			},
			output: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "*****", "user": "********"},
			},
		},
		{
			name:    "test: default options",
			options: nil,
			input: map[string]interface{}{
				"key_1": map[string]interface{}{"key": 12, "user": false},
			},
			output: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "**", "user": "*****"},
			},
		},
		{
			name:    "test: hit password length",
			options: []Option{CustomValueLength(3)},
			input: map[string]interface{}{
				"key_1": map[string]interface{}{"key": 12, "user": "12345"},
			},
			output: map[string]interface{}{
				"key_1": map[string]interface{}{"key": "**", "user": "***"},
			},
		},
	}

	for _, tc := range testCases {
		p := NewPrinter(tc.options...)

		p.maskValues(tc.input)

		assert.Equal(t, tc.output, tc.input, tc.name)
	}
}
