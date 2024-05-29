package printer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/samber/lo"
)

var ErrInvalidFormat = errors.New("invalid format")

// Printer a vkv printer needs to print out its data in yaml, json, markdown and a base format.
type Printer interface {
	Print() ([]byte, error)
}

// PrinterOptions options for the printer.
type PrinterOptions struct {
	Format string
	Writer io.Writer
}

// DefaultPrinterOptions returns the default printer options.
func DefaultPrinterOptions() PrinterOptions {
	return PrinterOptions{
		Writer: os.Stdout,
		Format: "default",
	}
}

// PrinterFunc function that returns the data to be printed.
type PrinterFunc func() ([]byte, error)

// PrinterFunc implements Printer interface.
func (pf PrinterFunc) Print() ([]byte, error) {
	return pf()
}

type PrinterFuncMap map[string]PrinterFunc

// SanitizerFunc function that allow customizing the output, for example masking certain values.
type SanitizerFunc func() error

// Print prints out the entities and applies the given sanitizer functions.
func Print(pf PrinterFuncMap, opts PrinterOptions, sf ...SanitizerFunc) error {
	if len(pf) == 0 {
		return fmt.Errorf("no printer available")
	}

	// find specific printer function
	p, ok := pf[strings.ToLower(opts.Format)]
	if !ok {
		return fmt.Errorf("\"%s\" %w. Available formats: %v", opts.Format, ErrInvalidFormat, lo.Keys(pf))
	}

	// apply any custom sanitizer funcs
	for _, f := range sf {
		if err := f(); err != nil {
			return fmt.Errorf("error while sanitizing output: %w", err)
		}
	}

	// print
	out, err := p.Print()
	if err != nil {
		return err
	}

	// write to writer
	fmt.Fprintln(opts.Writer, string(out))

	return nil
}
