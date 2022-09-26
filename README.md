# Blueprint

Template based project generation.

## Motivation

Creating a new project is a time consuming process. It would be nice to have a way to generate new projects form templates. This removes the need to set up the same boilerplate code over and over again.

If you know [python cookiecutter](https://cookiecutter.readthedocs.io/), you are familiar with the concept. However, there are some things that are not supported by cookiecutter. For example, raw and exclude files.

Furthermore, the kubernetes project [helm](https://helm.sh/) is a great example of handling templates. But it only works in the domain of kubernetes and cannot be used fo generic project generation.

Blueprint is a template based project generation tool, like cookiecutter, written in go with a helm flavoured API. The template engine is go-template.

## Synopsis

```console
Usage:
  blueprint [source] [destination] [flags]

Flags:
      --set stringArray      set one or more values separated by comma. i.e. a=1,b=2. Flag can be specified multiple times
  -f, --values stringArray   set values from a file or url. Flag can be specified multiple times
      --no-hooks             disable the pre and post hooks
  -h, --help                 show the help text
  -v, --version              show the version
```

The source is a git repository that contains the templates. The destination is the directory where the project will be generated. The source can also a be a local directory.

The `--set` and `--values` flag are optional. Usually required values are received via the configured inputs from the blueprint.yaml, when generating the project. In some advanced cases it can be useful to set values on the command line.

`--set` takes precedence over `--values`.

## Blueprint Layout

```console
.
├── blueprint.yaml
├── templates
│   └── _helpers.tpl
└── values.yaml
```

## The blueprint.yaml

```yaml
# A List of inputs to prompt the user for.
# The provided responses is merged into the values.
# If a value from the inputs was set with the
# --set flag or --values flag, the input prompt is skipped.
input:
  - "some.value"

# List of exclude objects containing a pattern to match
# relative to the project root, and an optional condition.
# If no condition is specified, the path is excluded,
# otherwise the condition is used to determine if the path
# should be excluded.
# The pattern follow .gitignore rules.
exclude:
  - pattern: "foo/**/*.txt"
    condition: "{{ not .Values.some.value }}"


# List of pattern to match relative to the project root.
# Matched files are not executed as template but directly
# copied to the output directory.
# The pattern follow .gitignore rules.
raw:
  - "foo/**/*.txt"

# List of pre- and post hook objects containing a script to
# run before and after the project has been created.
# The hook scripts are rendered as template before execution.
preHooks:
  - name: my-pre-hook
    script: |
      echo "Just before the project is created"
      echo "my value is: {{ .Values.some.value }}"
postHooks:
  - name: my-post-hook
    script: |
      echo "Just after the project is created"
```

Review the [json schema](./assets/schema/blueprint.json) for more details.

## The values file

The values file contains arbitrary values which can be referenced in the templates and path segments.

## Template Data

The below struct is passed to each template as context.

```go
type Data struct {
 Project     Project          // project related data, i.e project name
 Values      map[string]any   // The values from the values file
}
```

## Templates Dir

Each blueprint repo has a templates directory which contains the templates. This directory is traversed recursively, and found templates are rendered and added to the destination directory. Paths are also rendered as templates

## Inputs

When provided inputs from the user are parsed with helms [strval.Parser](https://github.com/helm/helm/blob/ee3f270e1eff0d462312635ad91cecd6f1fce620/pkg/strvals/parser.go) which is also used for the --set flag. That means all the rules for helms `--set` flag apply. Certain characters need to be escaped with a backslash `\`. Specifically, commas and curly braces.

## JSON Schema

You can use the [json schema](./assets/schema/blueprint.json) to validate your blueprint. For example if you are vscode and the redhat yaml extension is installed you can add a setting to your settings.json file to validate your blueprint.

```json
{
  "yaml.schemas": {
    "https://raw.githubusercontent.com/bluebrown/blueprint/main/assets/schema/blueprint.json": [
      "blueprint.yaml"
    ]
  }
}
```

## Installation

### Binary

Download the binary from the [release page](https://github.com/bluebrown/blueprint/releases). For example

```bash

curl -fsSL https://github.com/bluebrown/blueprint/releases/latest/download/blueprint-linux-amd64 >blueprint
chmod 755 blueprint
```

### Go

If you have go installed, you can use the `go install` command to install the binary.

```bash
go install github.com/bluebrown/blueprint/cmd/blueprint@latest
```

### Docker

The binary is also available as [docker image](https://hub.docker.com/repository/docker/bluebrown/blueprint). If you are using hooks, this might fail as some of the hooks dependencies are not available in the image. You can disable the hooks with the `--no-hooks` flag if you still want to use the image.

### From source

Clone the repo and use the makefile to build the binary. The make install command will move the binary to /usr/local/bin.

```bash
git clone https://github.com/bluebrown/blueprint
cd blueprint && make
```

## Example

The below examples fetches the blueprint repo and generates a project in the my-project directory. It will prompt for some inputs before generating the project.The provided inputs are used in the templates.

```bash
blueprint https://github.com/bluebrown/blueprint-example my-project
```

Some values may not be part of the inputs. It is still possible to set them. For example with the `--set flag`.

```bash
blueprint https://github.com/bluebrown/blueprint-example my-project --set service.enabled=false
```

If you are using docker, mount your workdir to /work and use the --user flag to get the correct permissions. Potentially set the --no-hooks flag to disable the hooks, if they depend on some programming language that isnt available in the image.

```bash
docker run --user "$(id -u):$(id -g)" --volume "$PWD:/work" --tty --interactive \
    bluebrown/blueprint https://github.com/bluebrown/blueprint-example my-project --no-hooks
```

## License

Some files in this repository contain embedded license notes. These files have been placed in the [lib directory](./lib/) with the given vendor name as package name.

The other files in this repository are licensed under the BSD 0-Clause License. See the [LICENSE file](./LICENSE) for more information.
