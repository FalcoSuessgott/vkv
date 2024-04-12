package namespace

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"

	"github.com/FalcoSuessgott/vkv/pkg/regex"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

// OutputFormat enum of valid output formats.
type OutputFormat int

const (
	// Base prints the secrets in the default format.
	Base OutputFormat = iota

	// YAML prints the secrets in yaml format.
	YAML

	// JSON prints the secrets in json format.
	JSON
)

var (
	defaultWriter = os.Stdout

	// ErrInvalidFormat invalid output format.
	ErrInvalidFormat = errors.New("invalid format (valid options: base, yaml, json, export, markdown)")
)

// Option list of available options for modifying the output.
type Option func(*Printer)

// Printer struct that holds all options used for displaying the secrets.
type Printer struct {
	format OutputFormat
	Regex  string
	writer io.Writer
}

// WithWriter option for passing a custom io.Writer.
func WithWriter(w io.Writer) Option {
	return func(p *Printer) {
		p.writer = w
	}
}

// WithRegex namespace regex.
func WithRegex(r string) Option {
	return func(p *Printer) {
		p.Regex = r
	}
}

// ToFormat sets the output format of the printer.
func ToFormat(format OutputFormat) Option {
	return func(p *Printer) {
		p.format = format
	}
}

// NewPrinter return a new printer struct.
func NewPrinter(opts ...Option) *Printer {
	p := &Printer{
		writer: defaultWriter,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Out prits out namespaces in various formats.
// nolint: cyclop
func (p *Printer) Out(ns map[string][]string) error {
	nsList := p.buildNamespaceList(ns)

	if len(ns) == 0 {
		return errors.New("no namespaces found")
	}

	if p.Regex != "" {
		var err error

		nsList, err = p.applyRegex(nsList)
		if err != nil {
			return err
		}
	}

	sort.Strings(nsList)

	switch p.format {
	case YAML:
		out, err := utils.ToYAML(map[string]interface{}{"namespaces": utils.RemoveDuplicates(nsList)})
		if err != nil {
			return err
		}

		fmt.Fprintln(p.writer, string(out))
	case JSON:
		out, err := utils.ToJSON(map[string]interface{}{"namespaces": utils.RemoveDuplicates(nsList)})
		if err != nil {
			return err
		}

		fmt.Fprintln(p.writer, string(out))
	case Base:
		for _, k := range utils.RemoveDuplicates(nsList) {
			fmt.Fprintln(p.writer, k)
		}
	default:
		return ErrInvalidFormat
	}

	return nil
}

func (p *Printer) buildNamespaceList(ns map[string][]string) []string {
	nsList := make([]string, 0)

	for k, v := range ns {
		if k != "" {
			nsList = append(nsList, k)
		}

		for _, i := range v {
			path := path.Join(k, i)

			if k != path {
				nsList = append(nsList, path)
			}
		}
	}

	return nsList
}

func (p *Printer) applyRegex(nsList []string) ([]string, error) {
	nsListRegex := make([]string, 0)

	for _, k := range nsList {
		match, err := regex.MatchRegex(p.Regex, k)
		if err != nil {
			return nil, err
		}

		if !match {
			continue
		}

		nsListRegex = append(nsListRegex, k)
	}

	return nsListRegex, nil
}
