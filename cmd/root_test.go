package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	v := "v1.1.1"

	c := rootCmd(v)
	b := bytes.NewBufferString("")

	c.SetArgs([]string{"-v"})
	c.SetOut(b)

	err := c.Execute()
	assert.NoError(t, err)

	out, _ := ioutil.ReadAll(b)
	assert.Equal(t, fmt.Sprintf("vkv: %s\n", v), string(out))
}
