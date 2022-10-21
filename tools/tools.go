//go:build tools

package tools

// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install mvdan.cc/gofumpt
//go:generate go install github.com/daixiang0/gci

import (
	_ "github.com/daixiang0/gci"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "mvdan.cc/gofumpt"
)
