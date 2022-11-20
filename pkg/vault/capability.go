package vault

import (
	"fmt"
	"log"
)

const (
	has    = "✔"
	hasNot = "✖"

	capCreate = "create"
	capRead   = "read"
	capUpdate = "update"
	capDelete = "delete"
	capList   = "list"
	capRoot   = "root"

	capabilities = "sys/capabilities-self"
)

// Capability represents a tokens caps for a specific path.
type Capability struct {
	Create bool
	Read   bool
	Update bool
	Delete bool
	List   bool
	Root   bool
}

// GetCapabilities returns the current authenticated tokens capabilities for a given path.
func (v *Vault) GetCapabilities(path string) (*Capability, error) {
	options := map[string]interface{}{
		"paths": []string{path},
	}

	res, err := v.Client.Logical().Write(capabilities, options)
	if err != nil {
		return nil, err
	}

	caps, ok := res.Data["capabilities"].([]interface{})
	if !ok {
		log.Fatal("could not read capabilities from response.")
	}

	//nolint predeclared
	cap := &Capability{}

	for _, c := range caps {
		switch c.(string) { //nolint forcetypeassert
		case capCreate:
			cap.Create = true
		case capRead:
			cap.Read = true
		case capUpdate:
			cap.Update = true
		case capDelete:
			cap.Delete = true
		case capList:
			cap.List = true
		case capRoot:
			cap.Root = true
		}
	}

	return cap, nil
}

func (c *Capability) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n",
		resolveCap(c.Create),
		resolveCap(c.Read),
		resolveCap(c.Update),
		resolveCap(c.Delete),
		resolveCap(c.List),
		resolveCap(c.Root))
}

func resolveCap(v bool) string {
	if v {
		return has
	}

	return hasNot
}
