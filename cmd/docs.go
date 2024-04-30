package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
hide:
  - toc
title: "%s"
---
`

var cmdDocPath = "./docs/cmd"

func NewDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "docs",
		Short:             "Generate the documentation for the CLI commands.",
		Hidden:            true,
		PersistentPreRunE: nil,
		RunE:              docsRun,
	}

	return cmd
}

func docsRun(cmd *cobra.Command, args []string) error {
	if err := os.MkdirAll(cmdDocPath, 0o777); err != nil {
		return err
	}

	err := doc.GenMarkdownTreeCustom(cmd.Root(), cmdDocPath, frontmatterPrepender, linkHandler)
	if err != nil {
		return err
	}

	return nil
}

func frontmatterPrepender(filename string) string {
	name := filepath.Base(filename)
	base := strings.TrimSuffix(name, path.Ext(name))
	title := strings.ReplaceAll(base, "_", " ")

	return fmt.Sprintf(fmTemplate, title)
}

func linkHandler(name string) string {
	base := strings.TrimSuffix(name, path.Ext(name))

	return strings.ToLower(base) + ".md"
}
