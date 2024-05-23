package vault

// Option list of available options for modifying the output.
type Option func(*FormatOptions)

// Printer struct that holds all options used for displaying the secrets.
type FormatOptions struct {
	ShowDiff       bool
	OnlyKeys       bool
	MaskSecrets    bool
	MaxValueLength int
}

// OnlyKeys flag for only showing secrets keys.
func OnlyKeys() Option {
	return func(p *FormatOptions) {
		p.OnlyKeys = true
	}
}

// hMaskSecrets flag for only showing secrets keys.
func MaskSecrets() Option {
	return func(p *FormatOptions) {
		p.MaskSecrets = true
	}
}

// WithMaskSecrets flag for only showing secrets keys.
func ShowDiff() Option {
	return func(p *FormatOptions) {
		p.ShowDiff = true
	}
}

// NewFormatOptions return a new printer struct.
func NewFormatOptions(opts ...Option) *FormatOptions {
	fOpts := &FormatOptions{}

	for _, opt := range opts {
		opt(fOpts)
	}

	fOpts.MaxValueLength = 12

	return fOpts
}
