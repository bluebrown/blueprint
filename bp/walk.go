package bp

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

type Matcher interface {
	Match(string, os.FileInfo) bool
}

// makeWalker returns a filepath.WalkFunc that renders all
// template found in the file tree with the given data and
// writes them to the output directory
func MakeFsWalker(t *template.Template, data *Data, outPath string, excludeMatcher, rawMatcher Matcher, helpersFileName string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		// return if its an error
		if err != nil {
			return err
		}

		// skip the base dir
		if d.IsDir() && path == "." {
			return nil
		}

		// skip the _helpers.tpl
		if d.Name() == helpersFileName {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		// skip this file or dir if its in the exclude list
		if excludeMatcher.Match(path, fileInfo) {
			return nil
		}

		// render the path as a template
		renderedPath, err := RenderString(t, path, data)
		if err != nil {
			return fmt.Errorf("failed to render path %s: %w", path, err)
		}

		// if its a directory, create it in the output dir
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(outPath, renderedPath), fileInfo.Mode())
		}

		// create the output file
		f, err := os.OpenFile(filepath.Join(outPath, renderedPath), os.O_CREATE|os.O_WRONLY, fileInfo.Mode())
		if err != nil {
			return fmt.Errorf("error creating output file: %w", err)
		}
		defer f.Close()

		// check if the path matches any of the raw paths
		// and if so, write the raw contents to the file and return
		if rawMatcher.Match(path, fileInfo) {
			srcFile, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("error opening raw file: %w", err)
			}
			defer srcFile.Close()
			_, err = io.Copy(f, srcFile)
			if err != nil {
				return fmt.Errorf("error copying raw file: %w", err)
			}
			return nil
		}

		// create a new template with the filename as the name
		// and parse the file at path
		t, err := t.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("error parsing template: %w", err)
		}

		// execute the template with the data
		err = t.ExecuteTemplate(f, d.Name(), data)
		if err != nil {
			return fmt.Errorf("error executing template: %w", err)
		}

		return nil
	}

}
