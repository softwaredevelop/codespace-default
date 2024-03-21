package linting_test

import (
	"cc/linting"
	"context"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestT(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_trivy_run", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_trivy_run")
		require.NotNil(t, c)

		_, err := linting.T(c, id).
			WithNewFile(".trivyignore", dagger.ContainerWithNewFileOpts{
				Contents: "AVD-DS-0002",
			}).
			WithNewFile("Dockerfile.test", dagger.ContainerWithNewFileOpts{
				Contents: "FROM docker.io/library/alpine:$VARIANT\nHEALTHCHECK none",
			}).
			WithExec([]string{"sh", "-c", "/trivy config --severity=MEDIUM,HIGH,CRITICAL --exit-code=1 $(find . -type f \\( -name Dockerfile -o -name Dockerfile.* \\))"}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_trivy_run_error", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_trivy_run_error")
		require.NotNil(t, c)

		_, err := linting.T(c, id).
			WithNewFile(".trivyignore", dagger.ContainerWithNewFileOpts{
				Contents: "AVD-DS-0002",
			}).
			WithNewFile("Dockerfile.test", dagger.ContainerWithNewFileOpts{
				Contents: "FROM docker.io/library/alpine:latest\nHEALTHCHECK none",
			}).
			WithExec([]string{"sh", "-c", "/trivy config --severity=MEDIUM,HIGH,CRITICAL --exit-code=1 $(find . -type f \\( -name Dockerfile -o -name Dockerfile.* \\))"}).
			Stdout(ctx)
		require.ErrorContains(t, err, "MEDIUM: 1")
	})
	t.Run("Test_trivy_version", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_trivy_version")
		require.NotNil(t, c)

		out, err := linting.T(c, id).
			WithExec([]string{"/trivy", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "Version")
	})
}
