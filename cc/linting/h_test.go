package linting_test

import (
	"cc/linting"
	"cc/util"
	"context"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestH(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_haskell_dockerfile_linter_run", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_hadolint_run")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		require.NoError(t, err)
		p := filepath.Join(dir, "..", "..")
		file := c.Host().
			Directory(p).
			File(".hadolint.yaml")
		require.NotNil(t, file)
		cont, err := file.Contents(ctx)
		require.NoError(t, err)
		require.Contains(t, cont, "docker.io")

		id, err := util.MountedHostDirectory(c, id, p, "/mountedtmp").
			ID(ctx)
		require.NoError(t, err)
		require.NotNil(t, id)

		_, err = linting.H(c, id).
			WithWorkdir("/mountedtmp").
			WithMountedFile("/.config/.hadolint.yaml", file).
			WithExec([]string{"sh", "-c", "/hadolint --config /.config/.hadolint.yaml $(find . -type f \\( -name Dockerfile -o -name Dockerfile.* \\))"}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_haskell_dockerfile_linter_config", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_hadolint_config")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		require.NoError(t, err)
		p := filepath.Join(dir, "..", "..")
		file := c.Host().
			Directory(p).
			File(".hadolint.yaml")
		require.NotNil(t, file)
		cont, err := file.Contents(ctx)
		require.NoError(t, err)
		require.Contains(t, cont, "docker.io")
	})
	t.Run("Test_haskell_dockerfile_linter_fail", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_haskell_dockerfile_linter_fail")
		require.NotNil(t, c)

		_, err := linting.H(c, id).
			WithNewFile("Dockerfile", dagger.ContainerWithNewFileOpts{
				Contents: "FROM docker.io/library/alpine:$VARIANT",
			}).
			WithNewFile("Dockerfile.test", dagger.ContainerWithNewFileOpts{
				Contents: "FROM docker.io/library/alpine:latest",
			}).
			WithExec([]string{"sh", "-c", "/hadolint $(find . -type f \\( -name Dockerfile -o -name Dockerfile.* \\))"}).
			Stdout(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "DL3007")
	})
	t.Run("Test_haskell_dockerfile_linter", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_haskell_dockerfile_linter")
		require.NotNil(t, c)

		out, err := linting.H(c, id).
			WithExec([]string{"/hadolint", "-v"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "Haskell Dockerfile Linter")
	})
}
