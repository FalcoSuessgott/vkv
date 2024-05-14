package engine

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"

	"github.com/FalcoSuessgott/vkv/pkg/regex"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
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
	ErrInvalidFormat = errors.New("invalid format (valid options: base, yaml, json)")
)

// Option list of available options for modifying the output.
type Option func(*Printer)

// Printer struct that holds all options used for displaying the secrets.
type Printer struct {
	format   OutputFormat
	regex    string
	nsPrefix bool
	writer   io.Writer
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
		p.regex = r
	}
}

// WithNSPrefix print engines with their ns prefix.
func WithNSPrefix(b bool) Option {
	return func(p *Printer) {
		p.nsPrefix = b
	}
}

// ToFormat sets the output format of the printer.
func ToFormat(format OutputFormat) Option {
	return func(p *Printer) {
		p.format = format
	}
}

// NewPrinter return a new printer struct.
func NewEnginePrinter(opts ...Option) *Printer {
	p := &Printer{
		writer: defaultWriter,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Out prints out engines.
// nolint: cyclop
func (p *Printer) Out(secrets interface{}) error {
	engines, ok := secrets.(vault.Engines)

	if !ok {
		return fmt.Errorf("invalid engines type: %T", engines)
	}

	engineList := p.buildEngineList(engines)

	if len(engineList) == 0 {
		return errors.New("no engines found")
	}

	if p.regex != "" {
		var err error

		engineList, err = p.applyRegex(engineList)
		if err != nil {
			return err
		}
	}

	sort.Strings(engineList)

	switch p.format {
	case YAML:
		out, err := utils.ToYAML(map[string]interface{}{"engines": utils.RemoveDuplicates(engineList)})
		if err != nil {
			return err
		}

		fmt.Fprint(p.writer, string(out))
	case JSON:
		out, err := utils.ToJSON(map[string]interface{}{"engines": utils.RemoveDuplicates(engineList)})
		if err != nil {
			return err
		}

		fmt.Fprint(p.writer, string(out))
	case Base:
		for _, k := range utils.RemoveDuplicates(engineList) {
			fmt.Fprintln(p.writer, k)
		}
	default:
		return ErrInvalidFormat
	}

	return nil
}

func (p *Printer) buildEngineList(engines map[string][]string) []string {
	engineList := make([]string, 0)

	for ns, eng := range engines {
		for _, e := range eng {
			if p.nsPrefix {
				engineList = append(engineList, path.Join(ns, e))
			} else {
				engineList = append(engineList, e)
			}
		}
	}

	return engineList
}

func (p *Printer) applyRegex(engines []string) ([]string, error) {
	engineListRegex := make([]string, 0)

	for _, e := range engines {
		match, err := regex.MatchRegex(p.regex, e)
		if err != nil {
			return nil, err
		}

		if !match {
			continue
		}

		engineListRegex = append(engineListRegex, e)
	}

	return engineListRegex, nil
}
