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
	showValues  bool
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

// ShowValues flag for unmasking secrets in output.
func ShowValues(b bool) Option {
	return func(p *Printer) {
		if b {
			p.showValues = true
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
	for k, v := range secrets {
		if !p.showValues {
			secrets[k] = p.maskValues(utils.ToMapStringInterface(v))
		}

		if p.onlyPaths {
			secrets[k] = p.printOnlyPaths(utils.ToMapStringInterface(v))
		}

		if p.onlyKeys {
			secrets[k] = p.printOnlykeys(utils.ToMapStringInterface(v))
		}
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
			for _, v := range utils.ToMapStringInterface(secrets[s]) {
				m, ok := v.(map[string]interface{})
				if !ok {
					log.Fatalf("cannot convert %T to map[string]interface", secrets[s])
				}

				for _, k := range utils.SortMapKeys(m) {
					fmt.Fprintf(p.writer, "export %s=\"%v\"\n", k, m[k])
				}
			}
		}
	case Markdown:
		headers, data := p.buildMarkdownTable(secrets)

		table := tablewriter.NewWriter(p.writer)
		table.SetHeader(headers)
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
		table.AppendBulk(data)
		// merge mounts and paths colunmn
		table.SetAutoMergeCellsByColumnIndex([]int{0, 1})
		table.Render()
	case Base:
		for _, k := range utils.SortMapKeys(secrets) {
			fmt.Fprintf(p.writer, "%s\n", k)
			m := utils.ToMapStringInterface(secrets[k])

			for _, i := range utils.SortMapKeys(m) {
				fmt.Fprintf(p.writer, "%s\n", i)
				p.printSecrets(utils.ToMapStringInterface(m[i]))
			}
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
		v := utils.ToMapStringInterface(secrets[s])
		for _, k := range utils.SortMapKeys(v) {
			m, ok := v[k].(map[string]interface{})
			if !ok {
				log.Fatalf("cannot convert %s to map[string]interface", secrets[s])
			}

			//nolint: gocritic
			if p.onlyPaths {
				headers = []string{"mount", "paths"}

				data = append(data, []string{s, k})
			} else if p.onlyKeys {
				headers = []string{"mount", "paths", "keys"}

				for _, j := range utils.SortMapKeys(m) {
					data = append(data, []string{s, k, j})
				}
			} else {
				headers = []string{"mount", "paths", "keys", "values"}

				for _, j := range utils.SortMapKeys(m) {
					data = append(data, []string{s, k, j, fmt.Sprintf("%v", m[j])})
				}
			}
		}
	}

	return headers, data
}

func (p *Printer) printOnlykeys(secrets map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}

	for k, v := range secrets {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		for j := range m {
			m[j] = ""
		}

		res[k] = m
	}

	return res
}

func (p *Printer) printOnlyPaths(secrets map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}

	for k, v := range secrets {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		for j := range m {
			m[j] = nil
		}

		res[k] = m
	}

	return res
}

func (p *Printer) maskValues(secrets map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}

	for k, v := range secrets {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		for j, v := range m {
			secret := fmt.Sprintf("%v", v)

			if len(secret) > p.valueLength && p.valueLength != -1 {
				m[j] = strings.Repeat(maskChar, p.valueLength)
			} else {
				m[j] = strings.Repeat(maskChar, len(secret))
			}
		}

		res[k] = m
	}

	return res
}

func (p *Printer) printSecrets(s map[string]interface{}) {
	for _, k := range utils.SortMapKeys(s) {
		if p.onlyKeys {
			fmt.Fprintf(p.writer, "\t%s\n", k)
		}

		if !p.onlyKeys && !p.onlyPaths {
			fmt.Fprintf(p.writer, "\t%s=%v\n", k, s[k])
		}
	}
}
