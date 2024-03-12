//revive:disable:package-comments,exported
package linting

import (
	"cc/util"
	"context"
	"path/filepath"

	"dagger.io/dagger"
)

func Hadolint(dir string, c *dagger.Client, id dagger.ContainerID, mountedDir string) error {
	ctx := context.Background()

	p := filepath.Join(dir, "..")
	file := c.Host().
		Directory(p).
		File(".hadolint.yaml")

	id, err := util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	if err != nil {
		return err
	}

	_, err = H(c, id).
		WithWorkdir(mountedDir).
		WithMountedFile("/.config/.hadolint.yaml", file).
		WithExec([]string{"sh", "-c", "/hadolint --config /.config/.hadolint.yaml $(find . -type f \\( -name Dockerfile -o -name Dockerfile.* \\))"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}

func H(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return h(c, id)
}

func h(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	hadolint := c.
		Container().
		From("hadolint/hadolint:latest-alpine")

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", hadolint.File("/bin/hadolint"))
}
