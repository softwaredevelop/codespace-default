package main_test

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"cc/format"
	"cc/linting"
	"cc/util"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

var c *dagger.Client
var id dagger.ContainerID

func TestMain(m *testing.M) {
	flag.Parse()

	ctx := context.Background()

	c, _ = dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	defer c.Close()

	c = c.Pipeline("code_quality_test")

	id, _ = c.
		Container().
		WithMountedTemp("/mountedtmp").
		ID(ctx)

	code := m.Run()
	defer c.Close()
	os.Exit(code)
}

func TestCodeQualityFunctions(t *testing.T) {
	t.Parallel()

	t.Run("Test_gofumpt_function", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_gofumpt_function")
		require.NotNil(t, c)

		mountedDir := "/mountedtmp"
		err := format.Gofumpt(c, id, mountedDir)
		require.NoError(t, err)
	})
	t.Run("Test_goimports_function", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_goimports_function")
		require.NotNil(t, c)

		mountedDir := "/mountedtmp"
		err := format.GoImports(c, id, mountedDir)
		require.NoError(t, err)
	})
	t.Run("Test_yamllint_function", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_yamllint_function")
		require.NotNil(t, c)

		dir, _ := os.Getwd()
		p := filepath.Join(dir, "..")
		err := linting.Yamllint(p, c)
		require.NoError(t, err)
	})
	t.Run("Test_revive_function", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_revive_function")
		require.NotNil(t, c)

		mountedDir := "/mountedtmp"
		err := linting.Revive(c, id, mountedDir)
		require.NoError(t, err)
	})
	t.Run("Test_editorconfig_checker_function", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_editorconfig_checker_function")
		require.NotNil(t, c)

		mountedDir := "/mountedtmp"
		err := linting.EditorconfigChecker(c, id, mountedDir)
		require.NoError(t, err)
	})
	t.Run("Test_actionlint_function", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_actionlint_function")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		require.NoError(t, err)
		require.NotEmpty(t, dir)
		require.IsType(t, "", dir)

		mountedDir := "/mountedtmp"
		err = linting.Actionlint(dir, c, id, mountedDir)
		require.NoError(t, err)
	})
}

func TestCodeQuality(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	t.Run("Test_revive_with_flag", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_revive_with_flag")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		mountedDir := "/mountedtmp"
		id, err = util.MountedHostDirectory(c, id, p, mountedDir).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		_, err = linting.ReviveL(c, id).
			WithWorkdir(mountedDir).
			WithExec([]string{
				"/revive",
				"-set_exit_status",
				"*/**",
			}).
			Stdout(ctx)
		require.Error(t, err)
	})
	t.Run("Test_revivel", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_revivel")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		mountedDir := "/mountedtmp"
		id, err = util.MountedHostDirectory(c, id, p, mountedDir).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		_, err = linting.ReviveL(c, id).
			WithWorkdir(mountedDir).
			WithExec([]string{"/revive", "-set_exit_status", "./..."}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_yamllint", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_yamllint")
		require.NotNil(t, c)

		id, err := c.
			Container().
			From("pipelinecomponents/yamllint").
			WithMountedTemp("/code").
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		mountedDir := "/code"
		_, err = util.MountedHostDirectory(c, id, p, mountedDir).
			WithWorkdir(mountedDir).
			WithExec([]string{"yamllint",
				"--config-data",
				"{extends: default, rules: {line-length: {level: warning}}}",
				"--no-warnings",
				"."}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_actionlint", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_actionlint")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..", ".github", "workflows")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		mountedDir := "/mountedtmp"
		id, err = util.MountedHostDirectory(c, id, p, mountedDir).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		_, err = linting.Al(c, id).
			WithWorkdir(mountedDir).
			WithExec([]string{"/actionlint",
				"-debug",
				"-pyflakes",
				"-shellcheck",
				"-verbose",
				"code-quality.yml",
			}).
			WithExec([]string{"/actionlint",
				"-debug",
				"-pyflakes",
				"-shellcheck",
				"-verbose",
				"unit-tests.yml",
			}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_actionlint_error", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_actionlint_error")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..", ".github", "workflows")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		mountedDir := "/mountedtmp"
		id, err = util.MountedHostDirectory(c, id, p, mountedDir).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		_, err = linting.Al(c, id).
			WithWorkdir(mountedDir).
			WithExec([]string{"/actionlint",
				"-debug",
				"-pyflakes",
				"-shellcheck",
				"-verbose",
				"a.yml",
			}).
			Stdout(ctx)
		require.Error(t, err)
	})
	t.Run("Test_client_pipeline", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_client_pipeline")
		require.NotNil(t, c)

		ccTest := c.Pipeline("cc_test")
		require.NotNil(t, ccTest)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		mountedDir := "/mountedtmp"
		id, err = util.MountedHostDirectory(ccTest, id, p, mountedDir).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		ec, err := linting.Ec(ccTest, id).
			Pipeline("ec_test").
			WithWorkdir(mountedDir).
			WithExec([]string{"/editorconfig-checker", "-verbose"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, ec)
	})
	t.Run("Test_go_format", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_go_format")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		mountedDir := "/mountedtmp"
		id, err = util.MountedHostDirectory(c, id, p, mountedDir).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		_, err = format.GoFormat(c, id).
			Pipeline("gofumpt_test").
			WithWorkdir(mountedDir).
			WithExec([]string{"/gofumpt", "-l", "-w", "."}).
			Stdout(ctx)
		require.NoError(t, err)

		_, err = format.GoFormat(c, id).
			Pipeline("goimports_test").
			WithWorkdir(mountedDir).
			WithExec([]string{"/goimports", "-l", "-w", "."}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_mounted_host_root_directory", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_mounted_host_root_directory")
		require.NotNil(t, c)

		container := c.
			Container().
			From("busybox:uclibc")
		require.NotNil(t, container)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		id, err := container.
			WithMountedDirectory("/mountedtmp", c.Host().Directory(p)).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		mntdir, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"ls", "/mountedtmp"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, mntdir)
	})
	t.Run("Test_container_ID", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_container_id")
		require.NotNil(t, c)

		container := c.Container().
			From("alpine")
		require.NotNil(t, container)

		id, err := container.ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		releaseName, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"/bin/sh", "-c", "cat /etc/os-release | awk -F= '/^NAME/ {print $2}' | tr -d '\"'"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Equal(t, "Alpine Linux\n", releaseName)
	})
	t.Run("Test_error_message", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_error_message")
		require.NotNil(t, c)

		_, err := c.
			Container().
			From("fake.invalid").
			ID(ctx)
		require.Error(t, err)
		require.ErrorContains(t, err, "not exist")
	})
	t.Run("Test_connect", func(t *testing.T) {
		t.Parallel()

		c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
		require.NoError(t, err)
		defer c.Close()
	})
}
