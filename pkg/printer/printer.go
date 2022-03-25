package printer

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/olekukonko/tablewriter"
)

// OutputFormat enum of valid output formats.
type OutputFormat int

const (
	maskChar = "*"

	// MaxValueLength maximum length of passwords.
	MaxValueLength = 12

	// Base prints the secrets in the default format.
	Base OutputFormat = iota

	// YAML prints the secrets in yaml format.
	YAML

	// JSON prints the secrets in json format.
	JSON

	// Export prints the secrets in export (export "key=value") format.
	Export

	// Markdown prints the secrets in markdowntable format.
	Markdown
)

var (
	defaultWriter = os.Stdout

	// ErrInvalidFormat invalid output format.
	ErrInvalidFormat = fmt.Errorf("invalid format (valid options: base, yaml, json, export, markdown)")
)

// Option list of available options for modifying the output.
type Option func(*Printer)

// Printer struct that holds all options used for displaying the secrets.
type Printer struct {
	format      OutputFormat
	writer      io.Writer
	onlyKeys    bool
	onlyPaths   bool
	showSecrets bool
	valueLength int
}

// CustomValueLength option for trimming down the output of secrets.
func CustomValueLength(length int) Option {
	return func(p *Printer) {
		p.valueLength = length
	}
}

// OnlyKeys flag for only showing secrets keys.
func OnlyKeys(b bool) Option {
	return func(p *Printer) {
		if b {
			p.onlyKeys = true
		}
	}
}

// OnlyPaths flag for only printing kv paths.
func OnlyPaths(b bool) Option {
	return func(p *Printer) {
		if b {
			p.onlyPaths = true
		}
	}
}

// ToFormat sets the output format of the printer.
func ToFormat(format OutputFormat) Option {
	return func(p *Printer) {
		p.format = format
	}
}

// WithWriter option for passing a custom io.Writer.
func WithWriter(w io.Writer) Option {
	return func(p *Printer) {
		p.writer = w
	}
}

// ShowSecrets flag for unmasking secrets in output.
func ShowSecrets(b bool) Option {
	return func(p *Printer) {
		if b {
			p.showSecrets = true
		}
	}
}

// NewPrinter return a new printer struct.
func NewPrinter(opts ...Option) *Printer {
	p := &Printer{
		writer:      defaultWriter,
		valueLength: MaxValueLength,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Out prints out the secrets according all configured options.
//nolint:cyclop
func (p *Printer) Out(secrets map[string]interface{}) error {
	if !p.showSecrets {
		p.maskSecrets(secrets)
	}

	if p.onlyPaths {
		p.printOnlyPaths(secrets)
	}

	if p.onlyKeys {
		p.printOnlykeys(secrets)
	}

	switch p.format {
	case YAML:
		out, err := utils.ToYAML(secrets)
		if err != nil {
			return err
		}

		fmt.Fprintf(p.writer, "%s", string(out))
	case JSON:
		out, err := utils.ToJSON(secrets)
		if err != nil {
			return err
		}

		fmt.Fprintf(p.writer, "%s", string(out))
	case Export:
		for _, s := range utils.SortMapKeys(secrets) {
			m, ok := secrets[s].(map[string]interface{})
			if !ok {
				log.Fatalf("cannot convert %s to map[string]interface", secrets[s])
			}

			for _, k := range utils.SortMapKeys(m) {
				fmt.Fprintf(p.writer, "export %s=%v\n", k, m[k])
			}
		}
	case Markdown:
		headers, data := p.buildMarkdownTable(secrets)

		table := tablewriter.NewWriter(p.writer)
		table.SetHeader(headers)
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.AppendBulk(data)

		// only merge cell for paths column
		if headers[0] == "paths" {
			table.SetAutoMergeCellsByColumnIndex([]int{0})
		}

		table.Render()
	case Base:
		for _, k := range utils.SortMapKeys(secrets) {
			fmt.Fprintf(p.writer, "%s\n", k)
			p.printSecrets(secrets[k])
		}
	default:
		return ErrInvalidFormat
	}

	return nil
}

func (p *Printer) buildMarkdownTable(secrets map[string]interface{}) ([]string, [][]string) {
	data := [][]string{}
	headers := []string{}

	for _, s := range utils.SortMapKeys(secrets) {
		m, ok := secrets[s].(map[string]interface{})
		if !ok {
			log.Fatalf("cannot convert %s to map[string]interface", secrets[s])
		}

		//nolint: gocritic
		if p.onlyPaths {
			headers = []string{"paths"}

			data = append(data, []string{s})
		} else if p.onlyKeys {
			headers = []string{"paths", "keys"}

			for _, k := range utils.SortMapKeys(m) {
				data = append(data, []string{s, k})
			}
		} else {
			headers = []string{"paths", "keys", "values"}

			for _, k := range utils.SortMapKeys(m) {
				data = append(data, []string{s, k, fmt.Sprintf("%v", m[k])})
			}
		}
	}

	return headers, data
}

func (p *Printer) printOnlykeys(secrets map[string]interface{}) {
	for k := range secrets {
		m, ok := secrets[k].(map[string]interface{})
		if !ok {
			continue
		}

		for k := range m {
			m[k] = ""
		}
	}
}

func (p *Printer) printOnlyPaths(secrets map[string]interface{}) {
	for k := range secrets {
		secrets[k] = nil
	}
}

func (p *Printer) maskSecrets(secrets map[string]interface{}) {
	for k := range secrets {
		m, ok := secrets[k].(map[string]interface{})
		if !ok {
			continue
		}

		for k := range m {
			secret := fmt.Sprintf("%v", m[k])

			if len(secret) > p.valueLength && p.valueLength != -1 {
				m[k] = strings.Repeat(maskChar, p.valueLength)
			} else {
				m[k] = strings.Repeat(maskChar, len(secret))
			}
		}
	}
}

func (p *Printer) printSecrets(s interface{}) {
	m, ok := s.(map[string]interface{})
	if ok {
		for _, k := range utils.SortMapKeys(m) {
			if p.onlyKeys {
				fmt.Fprintf(p.writer, "\t%s\n", k)
			}

			if !p.onlyKeys && !p.onlyPaths {
				fmt.Fprintf(p.writer, "\t%s=%v\n", k, m[k])
			}
		}
	}
}
