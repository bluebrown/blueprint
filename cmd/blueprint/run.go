package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bluebrown/blueprint/bp"
	"helm.sh/helm/v3/pkg/strvals"
	"sigs.k8s.io/yaml"
)

func run(ctx context.Context, input, output string, sets, vals []string, noHooks bool) error {
	var inPath string
	var err error

	curdir, err := os.Getwd()
	if err != nil {
		return err
	}

	// get the absolute path to the output directory
	outPath, err := filepath.Abs(output)
	if err != nil {
		return err
	}

	// check if the outpath already exists, if so check if it's empty
	// if it's not empty, return an error
	exists, err := bp.PathExists(outPath)
	if err != nil {
		return fmt.Errorf("error checking if output path exists: %w", err)
	}
	if exists {
		isEmpty, err := bp.IsEmptyDir(outPath)
		if err != nil {
			return fmt.Errorf("error while checking if output path is empty: %w", err)
		}
		if !isEmpty {
			return errors.New("output directory is not empty, refusing to overwrite")
		}
	}

	// clone the repo from upstream if required
	if bp.IsUpstreamRepo(input) {
		inPath, err = bp.CloneRepo(ctx, input)
		if err != nil {
			return err
		}
		// remove the tmp repo after we're done
		defer func() {
			os.Chdir(curdir)
			err := os.RemoveAll(inPath)
			if err != nil {
				fmt.Println("error removing repo:", err)
			}
		}()
	} else {
		inPath, err = filepath.Abs(input)
		if err != nil {
			return err
		}
	}

	// get the blueprint meta
	blueprintMetaFile := filepath.Join(inPath, configFileName)
	b, err := os.ReadFile(blueprintMetaFile)
	if err != nil {
		return err
	}

	var blueprintMeta bp.BlueprintMeta
	err = yaml.Unmarshal(b, &blueprintMeta)
	if err != nil {
		return err
	}

	// the data is passed to the templates as a struct
	data := &bp.Data{
		Project: bp.Project{
			// the project name is the name of the output directory
			Name: filepath.Base(filepath.Base(outPath)),
		},
		// the values are loaded from the values.yaml file
		Values: make(map[string]interface{}),
	}

	// read the values
	valuesFilePath := filepath.Join(inPath, valuesFileName)
	err = bp.ReadFile(valuesFilePath, &data.Values)
	if err != nil {
		return err
	}

	overrides := make(map[string]interface{})

	// read the values from the -f or --values flags
	// and merge them into the values
	for _, f := range vals {
		err = bp.ReadValues(f, &overrides)
		if err != nil {
			return err
		}
	}

	// merge the --set values into the values map
	setVals, err := strvals.Parse(strings.Join(sets, ","))
	if err != nil {
		return err
	}
	overrides = bp.MergeMaps(overrides, setVals)

	// remove the input question for values that have been set
	// via --set or --values
	newInput := make([]string, 0, len(blueprintMeta.Input))
	for _, q := range blueprintMeta.Input {
		if _, ok := bp.LookupValue(q, overrides); !ok {
			newInput = append(newInput, q)
		}
	}
	blueprintMeta.Input = newInput

	// merge the overrides into the values
	data.Values = bp.MergeMaps(data.Values, overrides)

	// get the user input
	inputs := make([]string, 0, len(blueprintMeta.Input))
	for _, key := range blueprintMeta.Input {
		val, err := bp.GetInput(key, data.Values)
		if err != nil {
			return err
		}
		// only use it if it has been provided via stdin
		// otherwise, the default is already in the values
		// so we don't need to parse and merge it
		// parsing it would also lead to complications
		// due to syntax differences
		if val != "" {
			inputs = append(inputs, fmt.Sprintf("%s=%s", key, val))
		}
	}

	// merge user input into the values map
	inputVals, err := strvals.Parse(strings.Join(inputs, ","))
	if err != nil {
		return err
	}
	data.Values = bp.MergeMaps(data.Values, inputVals)

	// make the output dir
	err = os.MkdirAll(outPath, 0755)
	if err != nil {
		return err
	}

	// create the base template
	t := bp.BaseTemplate()

	// run pre hook if it exists and we're not in no-hooks mode
	if !noHooks {
		for _, hook := range blueprintMeta.PreHooks {
			err := bp.RunHook(ctx, t, outPath, hook, data)
			if err != nil {
				return err
			}
		}
	}

	// enter the input dir
	err = os.Chdir(filepath.Join(inPath, templatesDir))
	if err != nil {
		return err
	}

	// parse the _helpers.tpl file first to make the helpers
	// available for all further templates
	if _, err := os.Stat(helpersFileName); err == nil {
		t, err = t.ParseFiles(helpersFileName)
		if err != nil {
			return fmt.Errorf("error the helpers file %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error checking for _helpers.txt: %w", err)
	}

	// make the raw rules
	rawMatcher, err := bp.NewFileMatcher(blueprintMeta.Raw)
	if err != nil {
		return err
	}

	// make the exclude rules
	excludes, err := bp.CompileExcludes(t, blueprintMeta.Exclude, data)
	if err != nil {
		return err
	}

	excludeMatcher, err := bp.NewFileMatcher(excludes)
	if err != nil {
		return err
	}

	// walk the input dir
	err = filepath.WalkDir(".", bp.MakeFsWalker(t, data, outPath, excludeMatcher, rawMatcher, helpersFileName))
	if err != nil {
		return err
	}

	// run post hook if it exists and we're not in no-hooks mode
	if !noHooks {
		for _, hook := range blueprintMeta.PostHooks {
			err := bp.RunHook(ctx, t, outPath, hook, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
