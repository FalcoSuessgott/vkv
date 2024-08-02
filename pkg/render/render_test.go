package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testTemplate = []byte("This is a {{ .Test }} template which replaces certain {{ .Values }}!")

func TestRenderTemplate(t *testing.T) {
	result, err := Apply(testTemplate, map[string]interface{}{
		"Test":   "test",
		"Values": "values",
	})

	require.NoError(t, err)
	assert.Equal(t, "This is a test template which replaces certain values!", string(result))
}

func TestRenderInvalidString(t *testing.T) {
	_, err := Apply([]byte("{{ invalid }"), map[string]interface{}{
		"Test":   "test",
		"Values": "values",
	})

	require.Error(t, err)
}

func TestRenderExpectError(t *testing.T) {
	_, err := Apply(testTemplate, map[string]interface{}{
		"Test":       "test",
		"WrongValue": "values",
	})

	require.Error(t, err)
}
