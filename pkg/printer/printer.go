package printer

import (
	"fmt"
	"io"
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

var defaultWriter = os.Stdout

// Option list of available options for modifying the output.
type Option func(*Printer)

// Printer struct that holds all options used for displaying the secrets.
type Printer struct {
	secrets     map[string]interface{}
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
func NewPrinter(m map[string]interface{}, opts ...Option) *Printer {
	p := &Printer{
		secrets:     m,
		writer:      defaultWriter,
		valueLength: MaxValueLength,
	}

	for _, opt := range opts {
		opt(p)
	}

	if !p.showSecrets {
		p.maskSecrets()
	}

	if p.onlyPaths {
		p.printOnlyPaths()
	}

	if p.onlyKeys {
		p.printOnlykeys()
	}

	return p
}

// Out prints out the secrets according all configured options.
//nolint:cyclop
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
	case Export:
		for _, s := range utils.SortMapKeys(p.secrets) {
			for _, k := range utils.SortMapKeys(p.secrets[s].(map[string]interface{})) {
				fmt.Fprintf(p.writer, "export %s=%v\n", k, p.secrets[s].(map[string]interface{})[k])
			}
		}
	case Markdown:
		headers, data := p.buildMarkdownTable()

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
	default:
		for _, k := range utils.SortMapKeys(p.secrets) {
			fmt.Fprintf(p.writer, "%s\n", k)
			p.printSecrets(p.secrets[k])
		}
	}

	return nil
}

func (p *Printer) buildMarkdownTable() ([]string, [][]string) {
	data := [][]string{}
	headers := []string{}

	for _, s := range utils.SortMapKeys(p.secrets) {
		//nolint: gocritic
		if p.onlyPaths {
			headers = []string{"paths"}

			data = append(data, []string{s})
		} else if p.onlyKeys {
			headers = []string{"paths", "keys"}

			for _, k := range utils.SortMapKeys(p.secrets[s].(map[string]interface{})) {
				data = append(data, []string{s, k})
			}
		} else {
			headers = []string{"paths", "keys", "values"}

			for _, k := range utils.SortMapKeys(p.secrets[s].(map[string]interface{})) {
				data = append(data, []string{s, k, fmt.Sprintf("%v", p.secrets[s].(map[string]interface{})[k])})
			}
		}
	}

	return headers, data
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
