package template

import (
	"bytes"
	"text/template"
)

// the base template with the funcmap
func BaseTemplate() *template.Template {
	t := template.New("")

	funcMap := FuncMap()

	funcMap["include"] = func(name string, v any) string {
		b := new(bytes.Buffer)
		t.ExecuteTemplate(b, name, v)
		return b.String()
	}

	funcMap["tpl"] = func(snippet string, v any) string {
		t, err := t.Parse(snippet)
		if err != nil {
			return ""
		}
		b := new(bytes.Buffer)
		t.Execute(b, v)
		return b.String()
	}

	t = t.Funcs(funcMap)

	return t
}
