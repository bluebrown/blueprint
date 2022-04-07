package values

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bluebrown/blueprint/lib/helm"
	"sigs.k8s.io/yaml"
)

// merge all maps prefering values found later in the list
func Merge(maps ...map[string]any) map[string]any {
	out := make(map[string]any)
	for _, m := range maps {
		out = helm.MergeMaps(out, m)
	}
	return out
}

// lookup nested values in a map by flat key path. i.e. foo.bar.baz
// if the map contains a dot, it must be escaped with a backslash
func Lookup(key string, data map[string]any) (value any, exists bool) {
	// split the keys into parts
	// swapping the escaped dot for a null byte
	key = strings.ReplaceAll(key, "\\.", "\x00")
	keys := strings.Split(key, ".")

	// loop through the keys
	for i, key := range keys {
		// swap the null byte back to the unescaped dot
		key = strings.ReplaceAll(key, "\x00", ".")

		// set val to the current value
		value, ok := data[key]
		if !ok {
			return nil, false
		}
		// if its the last key, return the value
		if i == len(keys)-1 {
			return value, true
		}
		data, ok = value.(map[string]any)
		if !ok {
			return nil, false
		}
	}
	return nil, false
}

// get the user input from the terminal
func GetInput(key string, data map[string]any) (input any, err error) {
	// lookup the default value to show in the prompt
	defaultValue, _ := Lookup(key, data)
	fmt.Printf("%s [%v]: ", key, defaultValue)

	// read the user input
	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		return err, nil
	}

	// return the trimmed input
	return strings.TrimSpace(line), nil

}

// read the values from a file at the given path into dist
func ReadFile(path string, dist any) error {
	// read the values file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	// unmarshal the yaml
	return yaml.Unmarshal(data, &dist)
}

// read the values from the given url into dist
func ReadURL(url string, dist any) error {
	res, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &dist)
}

// read the values from url or file into dist
func ReadValues(urlOrPath string, dist any) error {
	if strings.HasPrefix(urlOrPath, "http") {
		return ReadURL(urlOrPath, dist)
	}
	return ReadFile(urlOrPath, dist)
}
