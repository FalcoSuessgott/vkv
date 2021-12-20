package printer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
)

type outputFormat int

const (
	separator = " "
	tabChar   = '\t'
	minWidth  = 0
	tabWidth  = 8
	padding   = 1

	// MaxPasswordLength maximum length of passwords.
	MaxPasswordLength = 12
	maskChar          = "*"

	yaml outputFormat = iota
	json
)

var defaultWriter = os.Stdout

// Option list of available options for modifying the output.
type Option func(*Printer)

// Printer struct that holds all options used for displaying the secrets.
type Printer struct {
	passwordLength int
	format         outputFormat

	writer         io.Writer
	tabWriter *tabwriter.Writer

	showSecrets bool
	showMetadata bool
	onlyKeys       bool
	onlyPaths      bool
}

// OnlyKeys flag for only showing secrets keys.
func OnlyKeys(b bool) Option {
	return func(p *Printer) {
		p.onlyKeys = b
	}
}

// OnlyPaths flag for only printing kv paths.
func OnlyPaths(b bool) Option {
	return func(p *Printer) {
		p.onlyPaths = b
	}
}

// ToYAML outputformat to yaml.
func ToYAML(b bool) Option {
	return func(p *Printer) {
		if b {
			p.format = yaml
		}
	}
}

// ToJSON outputformat to yaml.
func ToJSON(b bool) Option {
	return func(p *Printer) {
		if b {
			p.format = json
		}
	}
}

// WithWriter option for passing a custom io.Writer.
func WithWriter(w io.Writer) Option {
	return func(p *Printer) {
		p.writer = w
	}
}

// CustomPasswordLength option for trimming down the output of secrets.
func CustomPasswordLength(length int) Option {
	return func(p *Printer) {
		p.passwordLength = length
	}
}

// ShowSecrets flag for unmasking secrets in output.
func ShowSecrets(b bool) Option {
	return func(p *Printer) {
		p.showSecrets = b
	}
}

// ShowMetadata flag for unmasking secrets in output.
func ShowMetadata(b bool) Option {
	return func(p *Printer) {
		p.showMetadata = b
	}
}

// NewPrinter return a new printer struct.
func NewPrinter(opts ...Option) *Printer {
	p := &Printer{
		writer:         defaultWriter,
		passwordLength: MaxPasswordLength,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.tabWriter = tabwriter.NewWriter(p.writer, minWidth, tabWidth, padding, tabChar, tabwriter.AlignRight)

	return p
}

// Out prints out the secrets according all configured options.
func (p *Printer) Out(s map[string]*vault.Secret) error {
	if p.onlyKeys {
		p.printOnlyKeys(s)
	}

	if p.onlyPaths {
		p.printOnlyPaths(s)
	}

	if !p.showSecrets {
		p.maskSecrets(s)
	}

	switch p.format {
	case yaml:
		out, err := utils.ToYAML(s)
		if err != nil {
			return err
		}

		fmt.Fprintf(p.writer, "%s", string(out))
	case json:
		out, err := utils.ToJSON(s)
		if err != nil {
			return err
		}

		fmt.Fprintf(p.writer, "%s", string(out))
	default:
		for k, v := range s {
			fmt.Fprintf(p.tabWriter, "%s\n", k)
			p.printSecrets(v)
		}
	}

	return nil
}

func (p *Printer) printOnlyKeys(secrets map[string]*vault.Secret) {
	for k := range secrets {
		for i := range secrets[k].Entries {
			secrets[k].Entries[i] = ""
		}
	}
}

func (p *Printer) printOnlyPaths(secrets map[string]*vault.Secret) {
	for k := range secrets {
		secrets[k].Entries[k] = nil
	}
}

func (p *Printer) maskSecrets(secrets map[string]*vault.Secret) {
	for k := range secrets {
		for i := range secrets[k].Entries {
			secret := fmt.Sprintf("%v", secrets[k].Entries[i])
			if len(secret) > p.passwordLength {
				secrets[k].Entries[i] = strings.Repeat(maskChar, p.passwordLength -2 ) + ".."
			} else {
				secrets[k].Entries[i] = strings.Repeat(maskChar, len(secret))
			}
		}
	}
}

func (p *Printer) printSecrets(s *vault.Secret) {
	for k, v := range  s.Entries {
		if p.onlyKeys {
			fmt.Fprintf(p.writer, "\t%s\n", k)
		}

		if p.showMetadata {
			fmt.Fprintf(p.tabWriter, "\t%s=%v%s\n", k, v, p.printMetadata(s.Metadata))
		}

		if (!p.onlyKeys && !p.onlyPaths) && !p.showMetadata {
			fmt.Fprintf(p.writer, "\t%s=%v\n", k, v)
		}
	}
}

func (p *Printer) printMetadata(m *vault.Metadata) string{
	s := ""

	s += fmt.Sprintf("\tversion: %s", m.Version)
	s += "\tcustom_metadata: "

	for k, v := range m.Metadata {
		s += fmt.Sprintf("%s=%v ", k, v)
	}


	return s
}