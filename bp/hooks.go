package bp

import (
	"context"
	"os"
	"os/exec"
	"text/template"
)

// run the given hook as shell command in the given directory
// the hook script is rendered as a template with the data before running it
func RunHook(ctx context.Context, t *template.Template, chdir string, hook Hook, data *Data) error {
	os.Chdir(chdir)
	script, err := RenderString(t, hook.Script, data)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
