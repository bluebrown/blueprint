package bp

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type fileMatcher struct {
	matcher gitignore.Matcher
}

func NewFileMatcher(paths []string) (*fileMatcher, error) {
	patterns := make([]gitignore.Pattern, len(paths))
	for i, path := range paths {
		patterns[i] = gitignore.ParsePattern(path, nil)
	}
	return &fileMatcher{gitignore.NewMatcher(patterns)}, nil
}

func (f *fileMatcher) Match(path string, info os.FileInfo) bool {
	return f.matcher.Match(strings.Split(path, "/"), info.IsDir())
}

// use the condition if it exists to determine if the path should be excluded,
// otherwise return true
func shouldExcludeFile(t *template.Template, e Exclude, data *Data) (bool, error) {
	if e.Condition != "" {
		boolString, err := RenderString(t, e.Condition, data)
		if err != nil {
			return true, err
		}
		return strconv.ParseBool(boolString)
	}
	return true, nil
}

// check for each item in the exclude list if they should be included
// in the resulting string slice and append them to the slice if they do
func CompileExcludes(t *template.Template, excludes []Exclude, data *Data) ([]string, error) {
	exSlice := make([]string, 0, len(excludes))
	for _, e := range excludes {
		ok, err := shouldExcludeFile(t, e, data)
		if err != nil {
			return nil, err
		}
		if ok {
			exSlice = append(exSlice, e.Pattern)
		}
	}
	return exSlice, nil
}

func IsEmptyDir(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
