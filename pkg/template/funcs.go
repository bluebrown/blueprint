package template

import "github.com/Masterminds/sprig/v3"

func FuncMap() map[string]any {
	return sprig.TxtFuncMap()
}
