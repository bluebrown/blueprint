package fs

import (
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type Matcher interface {
	Match(string, os.FileInfo) bool
}

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
