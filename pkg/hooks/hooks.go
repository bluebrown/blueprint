package hooks

import (
	"context"
	"os"
	"os/exec"
	"text/template"

	tpl "github.com/bluebrown/blueprint/pkg/template"
	"github.com/bluebrown/blueprint/pkg/types"
)

// run the given hook as shell command in the given directory
// the hook script is rendered as a template with the data before running it
func Run(ctx context.Context, t *template.Template, chdir string, hook types.Hook, data *types.Data) error {
	os.Chdir(chdir)
	script, err := tpl.RenderString(t, hook.Script, data)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
