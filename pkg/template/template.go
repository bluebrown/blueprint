package template

import (
	"text/template"
)

// the base template with the funcmap
func BaseTemplate() *template.Template {
	t := template.New("")
	t = t.Funcs(MakeFuncMap(t))
	return t
}
