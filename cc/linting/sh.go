//revive:disable:package-comments,exported
package linting

import (
	"cc/util"
	"context"
	"path/filepath"

	"dagger.io/dagger"
)

func Sh(dir string, c *dagger.Client, mountedDir string) error {
	ctx := context.Background()

	p := filepath.Join(dir, "..")
	id, err := c.
		Container().
		From("koalaman/shellcheck-alpine:latest").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	if err != nil {
		return err
	}

	_, err = util.MountedHostDirectory(c, id, p, mountedDir).
		WithWorkdir(mountedDir).
		// WithExec([]string{"find",
		// 	".",
		// 	"-type",
		// 	"f",
		// 	"-name",
		// 	"*.sh",
		// 	"-exec",
		// 	"shellcheck",
		// 	"{}",
		// 	";"}).
		WithExec([]string{"sh",
			"-c",
			"shellcheck $(find . -type f -name '*.sh')"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}
