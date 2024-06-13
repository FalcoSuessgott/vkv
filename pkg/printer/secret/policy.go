package secret

import (
	"fmt"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
	"github.com/FalcoSuessgott/vkv/pkg/vault"
	"github.com/juju/ansiterm"
)

const (
	header = "PATH\tCREATE\tREAD\tUPDATE\tDELETE\tLIST\tROOT\n"

	tabChar  = '\t'
	minWidth = 4
	tabWidth = 8
	padding  = 2
)

func (p *Printer) printPolicy(secrets map[string]interface{}) error {
	transformMap := make(map[string]interface{})
	utils.TransformMap(secrets, transformMap, "")

	capMap := make(map[string]*vault.Capability)

	for k := range transformMap {
		c, err := p.vaultClient.GetCapabilities(k)
		if err != nil {
			return err
		}

		capMap[k] = c
	}

	return p.printCapabilities(capMap)
}

func (p *Printer) printCapabilities(caps map[string]*vault.Capability) error {
	t := ansiterm.NewTabWriter(p.writer, minWidth, tabWidth, padding, tabChar, uint(ansiterm.Default))
	fmt.Fprint(t, header)

	for p, c := range caps {
		fmt.Fprintf(t, "%s\t%s", p, c.String())
	}

	if err := t.Flush(); err != nil {
		return err
	}

	return nil
}
