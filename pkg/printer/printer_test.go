package printer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrint(t *testing.T) {
	testCases := []struct {
		name           string
		printerFuncMap PrinterFuncMap
		sanitizerFuncs []SanitizerFunc
		exp            string
		errMsg         string
	}{
		{
			name:           "no printer available",
			printerFuncMap: PrinterFuncMap{},
			errMsg:         "no printer available",
		},
		{
			name: "basic",
			printerFuncMap: PrinterFuncMap{
				"default": func() ([]byte, error) {
					return []byte("unit testing is fun"), nil
				},
			},
			exp: "unit testing is fun\n",
		},
		{
			name: "invalid format",
			printerFuncMap: PrinterFuncMap{
				"invalid": func() ([]byte, error) {
					return []byte("unit testing is fun"), nil
				},
			},
			errMsg: "\"default\" invalid format. Available formats: [invalid]",
		},
		{
			name: "sanitizer",
			printerFuncMap: PrinterFuncMap{
				"default": func() ([]byte, error) {
					return []byte("unit testing is fun"), nil
				},
			},
			sanitizerFuncs: []SanitizerFunc{
				func() error {
					return nil
				},
			},
			exp: "unit testing is fun\n",
		},
		{
			name: "sanitizer error",
			printerFuncMap: PrinterFuncMap{
				"default": func() ([]byte, error) {
					return []byte("unit testing is fun"), nil
				},
			},
			sanitizerFuncs: []SanitizerFunc{
				func() error {
					return fmt.Errorf("error")
				},
			},
			errMsg: "error while sanitizing output: error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := bytes.NewBufferString("")

			opts := DefaultPrinterOptions()
			opts.Writer = b

			err := Print(tc.printerFuncMap, opts, tc.sanitizerFuncs...)

			if tc.errMsg != "" {
				require.Equal(t, tc.errMsg, err.Error(), "error msg")
			}

			if tc.exp != "" {
				require.Equal(t, tc.exp, b.String(), "output")
			}
		})
	}
}
