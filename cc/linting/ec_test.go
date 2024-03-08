package linting_test

import (
	"cc/linting"
	"context"
	"os"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestEditorconfigChecker(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_editorconfig_checker_latest", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_editorconfig_checker_latest")
		require.NotNil(t, c)

		out, err := linting.Ec(c, id).
			WithExec([]string{"/editorconfig-checker", "-help"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "USAGE")
	})
	t.Run("Test_editorconfig_checker_version", func(t *testing.T) {
		t.Parallel()

		version, err := linting.EcVersion()
		require.NoError(t, err)
		require.Contains(t, version, ".")
	})
	t.Run("Test_editorconfig_checker_on_golang_install", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_editorconfig_checker_on_golang_install")
		require.NotNil(t, c)

		container := linting.Ec2(c, id)
		require.NotNil(t, container)

		_, err := container.
			WithWorkdir("/tmp").
			WithExec([]string{"/ec", "-debug"}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_editorconfig_checker_on_golang_alpine", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_editorconfig_checker_on_golang_alpine")
		require.NotNil(t, c)

		container := linting.Ec1(c, id)
		require.NotNil(t, container)

		_, err := container.
			WithWorkdir("/tmp").
			WithExec([]string{"/ec", "-debug"}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_container_ID_busybox", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_container_ID_busybox")
		require.NotNil(t, c)

		id, err := c.
			Container().
			From("busybox:uclibc").
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		out, err := c.Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"busybox"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "BusyBox")
	})
	t.Run("Test_container_ID", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_container_id")
		require.NotNil(t, c)

		id, err := c.
			Container().
			From("alpine").
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		releaseName, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"/bin/sh", "-c", "cat /etc/os-release | awk -F= '/^NAME/ {print $2}' | tr -d '\"'"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Equal(t, "Alpine Linux\n", releaseName)
	})
	t.Run("Test_git_clone_file_content", func(t *testing.T) {
		t.Parallel()
		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
		require.NoError(t, err)
		defer c.Close()
		c = c.Pipeline("test_git_clone_file_content")
		require.NotNil(t, c)

		repoURL := "https://github.com/editorconfig-checker/editorconfig-checker.git"
		fileName := "VERSION"
		content, err := linting.GitCloneFileContent(ctx, c, repoURL, fileName)
		require.NoError(t, err)
		require.Contains(t, content, ".")
	})
	t.Run("Test_git_clone", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_git_clone")
		require.NotNil(t, c)

		repoURL := "https://github.com/editorconfig-checker/editorconfig-checker.git"
		repo := linting.GitClone(c, repoURL)
		require.NotNil(t, repo)
	})
}
