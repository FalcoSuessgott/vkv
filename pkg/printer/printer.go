package printer

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/FalcoSuessgott/vkv/pkg/render"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/disiqueira/gotree/v3"
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

	// Template renders a given template string or file.
	Template
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
	template    string
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

// WithTemplate sets the template file.
func WithTemplate(str, path string) Option {
	return func(p *Printer) {
		if str != "" {
			p.template = str

			return
		}

		if path != "" {
			str, err := utils.ReadFile(path)
			if err != nil {
				log.Fatalf("error reading %s: %s", path, err.Error())
			}

			p.template = string(str)

			return
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
//
//nolint:cyclop,gocognit
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
		for _, k := range utils.SortMapKeys(secrets) {
			m := utils.ToMapStringInterface(secrets[k])

			for _, i := range utils.SortMapKeys(m) {
				subMap, ok := m[i].(map[string]interface{})
				if !ok {
					log.Fatalf("cannot convert %T to map[string]interface", m[i])
				}

				for _, j := range utils.SortMapKeys(subMap) {
					fmt.Fprintf(p.writer, "export %s=\"%v\"\n", j, subMap[j])
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
		table.SetAutoMergeCellsByColumnIndex([]int{0, 1}) // merge mounts and paths colunmn
		table.Render()

	case Template:
		type entry struct {
			Path, Key string
			Value     interface{}
		}

		entries := []entry{}

		for _, k := range utils.SortMapKeys(secrets) {
			m := utils.ToMapStringInterface(secrets[k])

			for _, i := range utils.SortMapKeys(m) {
				subMap, ok := m[i].(map[string]interface{})
				if !ok {
					log.Fatalf("cannot convert %T to map[string]interface", m[i])
				}

				for _, j := range utils.SortMapKeys(subMap) {
					entries = append(entries, entry{Path: i, Key: j, Value: subMap[j]})
				}
			}
		}

		output, err := render.String([]byte(p.template), entries)
		if err != nil {
			return err
		}

		fmt.Fprintln(p.writer, output.String())

	case Base:
		for _, k := range utils.SortMapKeys(secrets) {
			tree := gotree.New(k + utils.Delimiter)
			m := utils.ToMapStringInterface(secrets[k])

			for _, i := range utils.SortMapKeys(m) {
				subMap, ok := m[i].(map[string]interface{})
				if !ok {
					log.Fatalf("cannot convert %T to map[string]interface", m[i])
				}

				// remove mount point of path
				path := strings.Join(strings.Split(i, utils.Delimiter)[1:], utils.Delimiter)

				tree.AddTree(p.printTree(path, subMap))
			}

			fmt.Fprint(p.writer, tree.Print())
		}
	default:
		return ErrInvalidFormat
	}

	return nil
}

func (p *Printer) printTree(path string, m map[string]interface{}) gotree.Tree {
	var tree gotree.Tree

	parts := strings.Split(path, utils.Delimiter)
	if len(parts) > 1 {
		tree = gotree.New(parts[0] + utils.Delimiter)
		tree.AddTree(p.printTree(strings.Join(parts[1:], utils.Delimiter), m))
	} else {
		tree = gotree.New(parts[0])
		for _, k := range utils.SortMapKeys(m) {
			if p.onlyKeys {
				tree.Add(k)
			}

			if !p.onlyKeys && !p.onlyPaths {
				tree.Add(fmt.Sprintf("%s=%v", k, m[k]))
			}
		}
	}

	return tree
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
