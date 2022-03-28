package fs

import (
	"strconv"
	"text/template"

	tpl "github.com/bluebrown/blueprint/pkg/template"
	"github.com/bluebrown/blueprint/pkg/types"
)

// use the condition if it exists to determine if the path should be excluded,
// otherwise return true
func shouldExclude(t *template.Template, e types.Exclude, data *types.Data) (bool, error) {
	if e.Condition != "" {
		boolString, err := tpl.RenderString(t, e.Condition, data)
		if err != nil {
			return true, err
		}
		return strconv.ParseBool(boolString)
	}
	return true, nil
}

// check for each item in the exclude list if they should be included
// in the resulting string slice and append them to the slice if they do
func CompileExcludes(t *template.Template, excludes []types.Exclude, data *types.Data) ([]string, error) {
	exSlice := make([]string, 0, len(excludes))
	for _, e := range excludes {
		ok, err := shouldExclude(t, e, data)
		if err != nil {
			return nil, err
		}
		if ok {
			exSlice = append(exSlice, e.Pattern)
		}
	}
	return exSlice, nil
}
