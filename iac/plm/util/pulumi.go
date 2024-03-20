//revive:disable:package-comments,exported
package util

import (
	"dagger.io/dagger"
)

func PulumiInstall(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	return pulumiInstall(c, id)
}

func pulumiInstall(c *dagger.Client, id dagger.ContainerID) *dagger.Container {
	install := c.
		Container().
		From("alpine").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add",
			"--no-cache",
			"curl",
		}).
		WithExec([]string{"/bin/sh", "-c", "curl -fsSL https://get.pulumi.com | sh"})

	return c.Container(dagger.ContainerOpts{ID: id}).
		WithFile("/", install.File("/root/.pulumi/bin/pulumi")).
		WithExec([]string{"mv", "/pulumi", "/usr/bin/pulumi"}).
		WithFile("/", install.File("/root/.pulumi/bin/pulumi-language-go")).
		WithExec([]string{"mv", "/pulumi-language-go", "/usr/bin/pulumi-language-go"})
}
