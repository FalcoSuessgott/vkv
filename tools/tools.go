//go:build tools

package tools

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install mvdan.cc/gofumpt
//go:generate go install github.com/daixiang0/gci
//go:generate go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt
import (
	// gci
	_ "github.com/daixiang0/gci"
	// golangci-lint
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	// gotestfmt
	_ "github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt"
	// gofumpt
	_ "mvdan.cc/gofumpt"
)
