package vault

import (
	"fmt"

	"github.com/FalcoSuessgott/vkv/pkg/utils"
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

	res, err := v.Client.Logical().WriteWithContext(v.Context, capabilities, options)
	if err != nil {
		return nil, err
	}

	caps, ok := res.Data["capabilities"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("could not read capabilities from response")
	}

	cap := &Capability{}

	for _, c := range caps {
		switch c.(string) { //nolint forcetypeassert
		case "create":
			cap.Create = true
		case "read":
			cap.Read = true
		case "update":
			cap.Update = true
		case "delete":
			cap.Delete = true
		case "list":
			cap.List = true
		case "root":
			// if the token has root capabilities, every capability is set to true
			cap.Root = true
			cap.Create = true
			cap.Read = true
			cap.Update = true
			cap.Delete = true
			cap.List = true
		}
	}

	return cap, nil
}

func (c *Capability) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n",
		utils.ResolveCap(c.Create),
		utils.ResolveCap(c.Read),
		utils.ResolveCap(c.Update),
		utils.ResolveCap(c.Delete),
		utils.ResolveCap(c.List),
		utils.ResolveCap(c.Root),
	)
}
