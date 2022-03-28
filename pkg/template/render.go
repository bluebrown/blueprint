package template

import (
	"bytes"
	"text/template"

	"github.com/bluebrown/blueprint/pkg/types"
)

// use the given template to parse and render the given string
// since this called template.Parse, the given string is associated with the template
// named templates are left intact
func RenderString(t *template.Template, s string, data *types.Data) (string, error) {
	t, err := t.Parse(s)
	if err != nil {
		return "", err
	}
	b := new(bytes.Buffer)
	err = t.Execute(b, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
