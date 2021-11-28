package printer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
)

type OutputFormat int

const (
	separator = " "
	maskChar  = "*"
	tabChar   = '\t'
	minWidth  = 0
	tabWidth  = 8
	padding   = 1

	YAML OutputFormat = iota
	JSON
)

var defaultWriter = os.Stdout

// Option list of available options for modifying the output.
type Option func(*Printer)

// Printer struct that holds all options used for displaying the secrets.
type Printer struct {
	secrets   map[string]interface{}
	format    OutputFormat
	writer    io.Writer
	tabWriter *tabwriter.Writer
	onlyKeys  bool
	onlyPaths bool
}

// OnlyKeys flag for only showing secrets keys.
func OnlyKeys(b bool) Option {
	return func(p *Printer) {
		if b {
			p.printOnlykeys()
		}
	}
}

// OnlyPaths flag for only printing kv paths.
func OnlyPaths(b bool) Option {
	return func(p *Printer) {
		if b {
			p.printOnlyPaths()
		}
	}
}

// ToYAML outputformat to yaml.
func ToYAML(b bool) Option {
	return func(p *Printer) {
		if b {
			p.format = YAML
		}
	}
}

// ToJSON outputformat to yaml.
func ToJSON(b bool) Option {
	return func(p *Printer) {
		if b {
			p.format = JSON
		}
	}
}

// WithWriter option for passing a custom io.Writer.
func WithWiter(w io.Writer) Option {
	return func(p *Printer) {
		p.writer = w
	}
}

// ShowSecrets flag for unmasking secrets in output.
func ShowSecrets(b bool) Option {
	return func(p *Printer) {
		if !b {
			p.maskSecrets()
		}
	}
}

// NewPrinter return a new printer struct.
func NewPrinter(m map[string]interface{}, opts ...Option) *Printer {
	p := &Printer{
		secrets: m,
		writer:  defaultWriter,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.tabWriter = tabwriter.NewWriter(p.writer, minWidth, tabWidth, padding, tabChar, tabwriter.AlignRight)

	return p
}

func (p *Printer) Out() error {
	switch p.format {
	case YAML:
		out, err := utils.ToYAML(p.secrets)
		if err != nil {
			return err
		}

		fmt.Fprintf(p.writer, "%s", string(out))
	case JSON:
		out, err := utils.ToJSON(p.secrets)
		if err != nil {
			return err
		}

		fmt.Fprintf(p.writer, "%s", string(out))
	default:
		for _, k := range utils.SortMapKeys(p.secrets) {
			// nolint
			if p.onlyKeys {
				fmt.Fprintf(p.tabWriter, "%s\t%v\n", k, printMap(p.secrets[k]))
			} else if p.onlyPaths {
				fmt.Fprintf(p.tabWriter, "%s\n", k)
			} else {
				fmt.Fprintf(p.tabWriter, "%s\t%v\n", k, printMap(p.secrets[k]))
			}
		}

		p.tabWriter.Flush()
	}

	return nil
}

func (p *Printer) printOnlykeys() {
	for k := range p.secrets {
		m, ok := p.secrets[k].(map[string]interface{})
		if !ok {
			continue
		}

		for k := range m {
			m[k] = ""
		}
	}
}

func (p *Printer) printOnlyPaths() {
	for k := range p.secrets {
		p.secrets[k] = nil
	}
}

func (p *Printer) maskSecrets() {
	for k := range p.secrets {
		m, ok := p.secrets[k].(map[string]interface{})
		if !ok {
			continue
		}

		for k := range m {
			secret := fmt.Sprintf("%v", m[k])
			m[k] = strings.Repeat(maskChar, len(secret))
		}
	}
}

func printMap(m interface{}) string {
	out := ""

	secrets, ok := m.(map[string]interface{})
	if !ok {
		return ""
	}

	for _, k := range utils.SortMapKeys(secrets) {
		out += k

		if secrets[k] == "" {
			out += separator
		} else {
			out += fmt.Sprintf("=%v ", secrets[k])
		}
	}

	return strings.TrimSuffix(out, separator)
}
