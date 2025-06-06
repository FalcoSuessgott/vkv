package secret

import (
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/savioxavier/termlink"
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
	ErrInvalidFormat = errors.New("invalid format (valid options: base, yaml, json, export, markdown, template, policy)")
)

// Option list of available options for modifying the output.
type Option func(*Printer)

// Printer struct that holds all options used for displaying the secrets.
type Printer struct {
	ctx            context.Context
	enginePath     string
	format         OutputFormat
	writer         io.Writer
	onlyKeys       bool
	onlyPaths      bool
	mergePaths     bool
	showVersion    bool
	showValues     bool
	showMetadata   bool
	withHyperLinks bool
	valueLength    int
	template       string
	vaultClient    *vault.Vault
}

// WithContext option for passing a custom context.
// nolint: fatcontext
func WithContext(ctx context.Context) Option {
	return func(p *Printer) {
		p.ctx = ctx
	}
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
		p.onlyKeys = b
	}
}

// WithHyperLinks for enabling hyperlinks.
func WithHyperLinks(b bool) Option {
	return func(p *Printer) {
		if b {
			p.withHyperLinks = termlink.SupportsHyperlinks()
		}
	}
}

// OnlyPaths flag for only printing kv paths.
func OnlyPaths(b bool) Option {
	return func(p *Printer) {
		p.onlyPaths = b
	}
}

// MergePaths flag for only printing kv paths.
func MergePaths(b bool) Option {
	return func(p *Printer) {
		p.mergePaths = b
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

// ShowValues flag for displaying the secrets version.
func ShowValues(b bool) Option {
	return func(p *Printer) {
		p.showValues = b
	}
}

// ShowVersion flag for unmasking secrets in output.
func ShowVersion(b bool) Option {
	return func(p *Printer) {
		p.showVersion = b
	}
}

// ShowMetadata flag for unmasking secrets in output.
func ShowMetadata(b bool) Option {
	return func(p *Printer) {
		p.showMetadata = b
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
			str, err := fs.ReadFile(path)
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

func WithEnginePath(path string) Option {
	return func(p *Printer) {
		p.enginePath = path
	}
}

// NewPrinter return a new printer struct.
func NewSecretPrinter(opts ...Option) *Printer {
	p := &Printer{
		writer:      defaultWriter,
		valueLength: MaxValueLength,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Update update printer applies the given options.
func Update(p *Printer, opts ...Option) {
	for _, opt := range opts {
		opt(p)
	}
}

// Out prints out the secrets according all configured options.
// nolint: cyclop
func (p *Printer) Out(secrets interface{}) error {
	secretMap := utils.ToMapStringInterface(secrets)

	for k, v := range secretMap {
		if !p.showValues {
			secretMap[k] = p.maskValues(utils.ToMapStringInterface(v))
		}

		if p.onlyPaths {
			secretMap[k] = p.printOnlyPaths(utils.ToMapStringInterface(v))
		}

		if p.onlyKeys {
			secretMap[k] = p.printOnlykeys(utils.ToMapStringInterface(v))
		}

		if p.mergePaths {
			secretMap = p.printMergePaths(utils.ToMapStringInterface(v), k)
		}
	}

	switch p.format {
	case YAML:
		return p.printYAML(secretMap)
	case JSON:
		return p.printJSON(secretMap)
	case Export:
		return p.printExport(secretMap)
	case Markdown:
		return p.printMarkdownTable(secretMap)
	case Template:
		return p.printTemplate(secretMap)
	case Base:
		return p.printBase(secretMap)
	case Policy:
		return p.printPolicy(secretMap)
	default:
		return ErrInvalidFormat
	}
}
