package template

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func FuncMap() map[string]any {
	return sprig.TxtFuncMap()
}

func MakeInclude(t *template.Template) func(string, any) string {
	return func(name string, v any) string {
		b := new(bytes.Buffer)
		if err := t.ExecuteTemplate(b, name, v); err != nil {
			return ""
		}
		return b.String()
	}
}

func MakeTpl(t *template.Template) func(string, any) string {
	return func(snippet string, v any) string {
		t, err := t.Parse(snippet)
		if err != nil {
			return ""
		}
		b := new(bytes.Buffer)
		if err := t.Execute(b, v); err != nil {
			return ""
		}
		return b.String()
	}
}

func MakeFuncMap(t *template.Template) map[string]any {
	funcMap := FuncMap()
	funcMap["include"] = MakeInclude(t)
	funcMap["tpl"] = MakeTpl(t)
	return funcMap
}
