//revive:disable:package-comments
package main

import (
	"context"
	"os"
	"path/filepath"
	"plm/util"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	c = c.Pipeline("infrastructure_as_code")
	id, err := c.
		Container().
		From("golang:alpine").
		WithMountedTemp("/mountedtmp").
		ID(ctx)
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	p := filepath.Join(dir, "..", "lcs")
	if err != nil {
		panic(err)
	}

	mountedDir := "/mountedtmp"
	id, err = util.MountedHostDirectory(c, id, p, mountedDir).
		ID(ctx)
	if err != nil {
		panic(err)
	}

	pat := os.Getenv("PULUMI_ACCESS_TOKEN")
	ght := os.Getenv("GITHUB_TOKEN")
	gho := os.Getenv("GITHUB_OWNER")
	id, err = util.PulumiInstall(c, id).
		Pipeline("setup_pulumi_environment").
		WithWorkdir(mountedDir).
		WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
		WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
		WithEnvVariable("PULUMI_ACCESS_TOKEN", pat).
		WithEnvVariable("GITHUB_TOKEN", ght).
		WithEnvVariable("GITHUB_OWNER", gho).
		WithExec([]string{"pulumi", "login"}).
		ID(ctx)
	if err != nil {
		panic(err)
	}

	_, err = c.Container(dagger.ContainerOpts{ID: id}).
		Pipeline("pulumi_local_source").
		WithWorkdir(mountedDir).
		WithExec([]string{"go", "run", "-v", "local.go"}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}
}
