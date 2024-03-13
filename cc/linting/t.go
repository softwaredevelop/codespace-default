//revive:disable:package-comments,exported
package linting

import (
	"cc/util"
	"context"
	"path/filepath"

	"dagger.io/dagger"
)

func Trivy(dir string, c *dagger.Client, id dagger.ContainerID, mountedDir string) error {
	ctx := context.Background()
	p := filepath.Join(dir, "..")
	id, err := util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	if err != nil {
		return err
	}

	_, err = t(c, id).
		WithNewFile(".trivyignore", dagger.ContainerWithNewFileOpts{
			Contents: "AVD-DS-0002",
		}).
		WithExec([]string{"sh", "-c", "/trivy config --severity=MEDIUM,HIGH,CRITICAL --exit-code=1 $(find . -type f \\( -name Dockerfile -o -name Dockerfile.* \\))"}).
		Stdout(ctx)
	if err != nil {
		return err
	}
	return nil
}

func T(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return t(c, id)
}

func t(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("alpine").
		WithExec([]string{
			"apk",
			"add",
			"--no-cache",
			"curl",
		}).
		WithExec([]string{"sh",
			"-c",
			"curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin"})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", install.File("/usr/local/bin/trivy"))
}
