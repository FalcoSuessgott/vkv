package main

import (
	"context"
	"os"
	"syscall"

	"github.com/FalcoSuessgott/vkv/cmd"
	"github.com/charmbracelet/fang"
)

var version string

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := fang.Execute(ctx, cmd.NewRootCmd(),
		fang.WithNotifySignal(syscall.SIGINT, syscall.SIGTERM),
		fang.WithVersion(version),
	); err != nil {
		os.Exit(1)
	}
}
