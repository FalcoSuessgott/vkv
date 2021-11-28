package utils

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
)

type Keys []string

func ToJSON(m map[string]interface{}) ([]byte, error) {
	out, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func ToYAML(m map[string]interface{}) ([]byte, error) {
	out, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func SortMapKeys(m map[string]interface{}) []string {
	keys := make(Keys, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Sort(keys)

	return keys
}

func (k Keys) Len() int {
	return len(k)
}

func (k Keys) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func (k Keys) Less(i, j int) bool {
	k1 := strings.ReplaceAll(k[i], "/", "\x00")
	k2 := strings.ReplaceAll(k[j], "/", "\x00")

	return k1 < k2
}
