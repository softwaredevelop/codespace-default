package linting_test

import (
	"context"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestSh(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_shellcheck_container_shellcheck_command_fail", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_shellcheck_container_shellcheck_command_fail")
		require.NotNil(t, c)

		_, err := c.
			Container().
			From("koalaman/shellcheck-alpine:latest").
			WithWorkdir("/mountedtmp").
			WithNewFile("test_fail.sh", dagger.ContainerWithNewFileOpts{
				Contents: "",
			}).
			WithExec([]string{"sh", "-c", "shellcheck $(find . -type f -name '*.sh')"}).
			Stdout(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "SC2148")
	})
	t.Run("Test_shellcheck_container_shellcheck_command", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_shellcheck_container_shellcheck_command")
		require.NotNil(t, c)

		out, err := c.
			Container().
			From("koalaman/shellcheck-alpine:latest").
			WithWorkdir("/mountedtmp").
			WithNewFile("test.sh", dagger.ContainerWithNewFileOpts{
				Contents: "#!/usr/bin/env bash",
			}).
			WithExec([]string{"sh", "-c", "shellcheck $(find . -type f -name '*.sh')"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotNil(t, out)
	})
	t.Run("Test_shellcheck_container_find_command", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_shellcheck_container_find_command")
		require.NotNil(t, c)

		out, err := c.
			Container().
			From("koalaman/shellcheck-alpine:latest").
			WithExec([]string{"find",
				".",
				"-type",
				"f",
				"-name",
				"*.sh",
				"-exec",
				"shellcheck",
				"{}",
				";"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotNil(t, out)
	})
	t.Run("Test_shellcheck_container_find", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_shellcheck_container_find")
		require.NotNil(t, c)

		out, err := c.
			Container().
			From("koalaman/shellcheck-alpine:latest").
			WithExec([]string{"find", "--help"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotNil(t, out)
	})
	t.Run("Test_shellcheck_container_shellcheck", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_shellcheck_container_shellcheck")
		require.NotNil(t, c)

		out, err := c.
			Container().
			From("koalaman/shellcheck-alpine:latest").
			WithExec([]string{"shellcheck", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "version")
	})
}
