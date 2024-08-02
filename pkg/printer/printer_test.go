package printer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrint(t *testing.T) {
	testCases := []struct {
		name           string
		printerFuncMap PrinterFuncMap
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := bytes.NewBufferString("")

			opts := DefaultPrinterOptions()
			opts.Writer = b
			err := Print(tc.printerFuncMap, opts)

			if tc.errMsg != "" {
				require.Equal(t, tc.errMsg, err.Error(), "error msg "+tc.name)
			}

			if tc.exp != "" {
				require.Equal(t, tc.exp, b.String(), "output "+tc.name)
			}
		})
	}
}
