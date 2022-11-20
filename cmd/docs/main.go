package main

import (
	"io"
	"log"

	"github.com/FalcoSuessgott/vkv/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd := cmd.NewRootCmd("", io.Discard)

	err := doc.GenMarkdownTree(cmd, ".")
	if err != nil {
		log.Fatal(err)
	}
}
