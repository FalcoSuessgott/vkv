package printer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/samber/lo"
)

// ErrInvalidPrinterFormat error for invalid format.
var ErrInvalidPrinterFormat = errors.New("invalid format")

// Printer a vkv printer needs to print out its data in yaml, json, markdown and a base format.
type Printer interface {
	PrintFormat() ([]byte, error)
}

// PrinterOptions options for the printer.
type PrinterOptions struct {
	Writer io.Writer
	Format string
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
func (pf PrinterFunc) PrintFormat() ([]byte, error) {
	return pf()
}

// PrinterFuncMap is a map of all implemented PrinterFuncs.
type PrinterFuncMap map[string]PrinterFunc

// Print prints out the entities and applies the given sanitizer functions.
func Print(pfm PrinterFuncMap, opts PrinterOptions) error {
	if len(pfm) == 0 {
		return fmt.Errorf("no printer available")
	}

	// find specific printer function
	pf, ok := pfm[strings.ToLower(opts.Format)]
	if !ok {
		return fmt.Errorf("\"%s\" %w. Available formats: %v", opts.Format, ErrInvalidPrinterFormat, lo.Keys(pfm))
	}

	out, err := pf.PrintFormat()
	if err != nil {
		return fmt.Errorf("failed to print format: %w", err)
	}

	// write to writer
	if _, err := fmt.Fprintln(opts.Writer, string(out)); err != nil {
		return fmt.Errorf("failed to write to writer: %w", err)
	}

	return nil
}
