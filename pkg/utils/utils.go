package utils

// nolint: staticcheck
import (
	"bytes"
	"encoding/json"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/ghodss/yaml"
)

const (
	// Delimiter / delimiter for splitting a path.
	Delimiter = "/"
)

// Keys type for receiving all keys of a map.
type Keys []string

// TransformMap takes a multi leveled map and returns a map with its combined paths
// as the keys and the map as its value. Also see TestTransformMap().
func TransformMap(a, b map[string]interface{}, key string) {
	for k, v := range a {
		// if its a map -> go deeper
		m, ok := v.(map[string]interface{})
		if ok {
			TransformMap(m, b, path.Join(key, k))
		} else { // otherwise add the key and value to the map
			b[key] = a
		}
	}
}

// PathMap takes a path like "a/b/c" and returns a map like map[a] -> map[b] -> map[c].
// if isSecretPath is true, then c does not have a / as suffix.
func PathMap(path string, s map[string]interface{}, isSecretPath bool) map[string]interface{} {
	m := map[string]interface{}{}

	parts := strings.Split(path, Delimiter)

	if path == "" {
		return s
	}

	if len(parts) > 1 {
		m[parts[0]+Delimiter] = PathMap(strings.Join(parts[1:], Delimiter), s, isSecretPath)
	} else {
		// if path leads to a vault kv directory, append a "/"
		if !isSecretPath {
			path += Delimiter
		}

		m[path] = s
	}

	return m
}

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

// ToMapStringInterface takes any value and returns the map string interface.
func ToMapStringInterface(i interface{}) map[string]interface{} {
	var m map[string]interface{}

	data, err := json.Marshal(i)
	if err != nil {
		log.Fatalf("cannot convert %v to map[string]interface: %v", i, err)
	}

	if err := json.Unmarshal(data, &m); err != nil {
		log.Fatalf("cannot convert %v to map[string]interface: %v", i, err)
	}

	return m
}

// ToJSON marshalls a given map to json.
func ToJSON(m interface{}) ([]byte, error) {
	// avoid encoding < and > to \u003c and \u003e
	data := &bytes.Buffer{}
	enc := json.NewEncoder(data)
	enc.SetEscapeHTML(false)

	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}

	// pretty print json
	res := &bytes.Buffer{}
	if err := json.Indent(res, data.Bytes(), "", "  "); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

// FromJSON takes a json byte array and marshalls it into a map.
func FromJSON(b []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

// ToYAML marshalls a given map to yaml.
func ToYAML(m interface{}) ([]byte, error) {
	out, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// FromYAML takes a yaml byte array and marshalls it into a map.
func FromYAML(b []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
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

// RemoveExtension removes the extension of a specified filename.
func RemoveExtension(file string) string {
	return strings.TrimSuffix(file, filepath.Ext(file))
}

// RemoveCarriageReturns removes \r mostly used for unit tests on windows OS.
func RemoveCarriageReturns(s string) string {
	return strings.ReplaceAll(s, "\r", "")
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

// RemoveDuplicates removes duplicate elements from a string slice.
func RemoveDuplicates(s []string) []string {
	allKeys := make(map[string]bool)

	list := []string{}

	for _, item := range s {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true

			list = append(list, item)
		}
	}

	return list
}

// DeepMergeMaps takes two maps and deeply merges them together.
// https://stackoverflow.com/questions/62953360/golang-merge-deeply-two-maps/62954592#62954592
func DeepMergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}

	for k, v := range b {
		// If you use map[string]interface{}, ok is always false here.
		// Because yaml.Unmarshal will give you map[interface{}]interface{}.
		if v, ok := v.(map[string]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]interface{}); ok {
					out[k] = DeepMergeMaps(bv, v)

					continue
				}
			}
		}

		out[k] = v
	}

	return out
}

// HandleEnginePath handles the engine path if one is specified.
func HandleEnginePath(enginePath, path string) (string, string) {
	// if engine path has been specified use that value as the root path and append the path
	if enginePath != "" {
		return enginePath, path
	}

	return SplitPath(path)
}

func ParseEnvs(prefix string, i interface{}) error {
	opts := env.Options{
		Prefix: prefix,
	}

	if err := env.Parse(i, opts); err != nil {
		return err
	}

	return nil
}
