package utils

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
)

const (
	// Delimiter / delimiter for splitting a path.
	Delimiter = "/"
)

// Keys type for receiving all keys of a map.
type Keys []string

// SplitPath splits a given path by / and returns the first element and the joined rest paths.
func SplitPath(path string) (string, string) {
	parts := removeEmptyElements(strings.Split(path, Delimiter))

	if len(parts) >= 2 {
		return parts[0], strings.Join(parts[1:], Delimiter)
	}

	return strings.Join(parts, Delimiter), ""
}

func removeEmptyElements(s []string) []string {
	r := []string{}

	for _, e := range s {
		if e != "" {
			r = append(r, e)
		}
	}

	return r
}

// ToJSON marshalls a given map to json.
func ToJSON(m map[string]interface{}) ([]byte, error) {
	out, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// ToYAML marshalls a given map to yaml.
func ToYAML(m map[string]interface{}) ([]byte, error) {
	out, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// SortMapKeys sorts the keys of a map.
func SortMapKeys(m map[string]interface{}) []string {
	keys := make(Keys, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Sort(keys)

	return keys
}

// Len returns the length of Keys.
func (k Keys) Len() int {
	return len(k)
}

// Swap swaps keys alphabetically.
func (k Keys) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

// Less compares keys alphabetically.
func (k Keys) Less(i, j int) bool {
	k1 := strings.ReplaceAll(k[i], "/", "\x00")
	k2 := strings.ReplaceAll(k[j], "/", "\x00")

	return k1 < k2
}
