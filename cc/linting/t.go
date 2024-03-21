//revive:disable:package-comments,exported
package linting

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

func Trivy(c *dagger.Client, id dagger.ContainerID) error {
	ctx := context.Background()

	files, err := t(c, id).
		WithExec([]string{"sh", "-c", "find . -type f \\( -name Dockerfile -o -name Dockerfile.* \\)"}).
		Stdout(ctx)
	if err != nil {
		return err
	}

	for _, file := range strings.Split(files, "\n") {
		if file == "" {
			continue
		}
		_, err = t(c, id).
			WithNewFile(".trivyignore", dagger.ContainerWithNewFileOpts{
				Contents: "AVD-DS-0002",
			}).
			WithExec([]string{"sh", "-c", fmt.Sprintf("/trivy config --severity=MEDIUM,HIGH,CRITICAL --exit-code=1 %s", file)}).
			Stdout(ctx)
		if err != nil {
			return err
		}
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
