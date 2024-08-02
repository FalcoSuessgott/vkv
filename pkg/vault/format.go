package vault

import (
	"log"

	"github.com/FalcoSuessgott/vkv/pkg/fs"
)

// Option list of available options for modifying the output.
type Option func(*FormatOptions)

// Printer struct that holds all options used for displaying the secrets.
type FormatOptions struct {
	showDiff    bool
	maskSecrets bool
	template    []byte
}

// hMaskSecrets flag for only showing secrets keys.
func MaskSecrets() Option {
	return func(p *FormatOptions) {
		p.maskSecrets = true
	}
}

// WithMaskSecrets flag for only showing secrets keys.
func ShowDiff() Option {
	return func(p *FormatOptions) {
		p.showDiff = true
	}
}

// WithTemplate sets the template file.
func WithTemplate(str, path string) Option {
	return func(p *FormatOptions) {
		if str != "" {
			p.template = []byte(str)

			return
		}

		if path != "" {
			out, err := fs.ReadFile(path)
			if err != nil {
				log.Fatalf("error reading %s: %s", path, err.Error())
			}

			p.template = out

			return
		}
	}
}

// NewFormatOptions return a new printer struct.
func NewFormatOptions(opts ...Option) *FormatOptions {
	fOpts := &FormatOptions{}

	for _, opt := range opts {
		opt(fOpts)
	}

	return fOpts
}
