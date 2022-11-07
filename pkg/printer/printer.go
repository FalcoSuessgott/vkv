package printer

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
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

	// Policy displays the current token policy capabilities for each path in a matrix.
	Policy
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
	vaultClient *vault.Vault
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

// WithVaultClient inject a vault client.
func WithVaultClient(v *vault.Vault) Option {
	return func(p *Printer) {
		p.vaultClient = v
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
		return p.printYAML(secrets)
	case JSON:
		return p.printJSON(secrets)
	case Export:
		return p.printExport(secrets)
	case Markdown:
		return p.printMarkdownTable(secrets)
	case Template:
		return p.printTemplate(secrets)
	case Base:
		return p.printBase(secrets)
	case Policy:
		return p.printPolicy(secrets)
	default:
		return ErrInvalidFormat
	}
}
