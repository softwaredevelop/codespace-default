package main_test

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"plm/util"
	"testing"

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

	c = c.Pipeline("pulumi_local_test")

	id, _ = c.
		Container().
		From("busybox:uclibc").
		WithMountedTemp("/mountedtmp").
		ID(ctx)

	code := m.Run()
	defer c.Close()
	os.Exit(code)
}

func TestPulumiLocal(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_pulumi_local", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_pulumi_local")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..", "lcs")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		container := c.
			Container().
			From("golang:alpine")
		require.NotNil(t, container)

		id, err := container.
			WithMountedDirectory("/mountedtmp", c.Host().Directory(p)).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		id, err = util.PulumiInstall(c, id).
			WithWorkdir("/mountedtmp").
			WithEnvVariable("PULUMI_SKIP_UPDATE_CHECK", "true").
			WithEnvVariable("PULUMI_CONFIG_PASSPHRASE", "").
			WithExec([]string{"pulumi", "login", "--local"}).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		reMatching := "TestLocalProject"
		_, err = c.Container(dagger.ContainerOpts{ID: id}).
			Pipeline("pulumi_localproject_test").
			WithExec([]string{"go", "test", "-v", "-run", reMatching}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_mounted_host_parent_directory", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_mounted_host_parent_directory")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		container := c.
			Container().
			From("busybox:uclibc")
		require.NotNil(t, container)

		id, err := container.
			WithMountedDirectory("/mountedtmp", c.Host().Directory(p)).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		out, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"ls", "/mountedtmp"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, out)
	})
	t.Run("Test_pulumi_install", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_pulumi_install")
		require.NotNil(t, c)

		container := util.PulumiInstall(c, id)
		require.NotNil(t, container)
		out, err := container.
			WithExec([]string{"ls", "/usr/bin/"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "pulumi")
		require.Contains(t, out, "pulumi-language-go")
	})
	t.Run("Test_container_id", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_container_id")
		require.NotNil(t, c)

		out, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"busybox"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "BusyBox")
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
}
