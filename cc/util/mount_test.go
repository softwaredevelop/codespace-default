package util_test

import (
	"cc/util"
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
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

	c = c.Pipeline("util_test")

	id, _ = c.
		Container().
		From("busybox:uclibc").
		WithMountedTemp("/mountedtmp").
		ID(ctx)

	code := m.Run()
	defer c.Close()
	os.Exit(code)
}

func TestUtil(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_mounted_host_directory_include_files", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_mounted_host_directory_include_files")
		require.NotNil(t, c)

		dir, err := os.Getwd()
		p := filepath.Join(dir, "..", "..")
		require.NoError(t, err)
		require.NotEmpty(t, p)

		mountedDir := "/mountedtmp"
		id, err = util.MountedHostDirectory(c, id, p, mountedDir).
			ID(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, id)

		mntdir, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"ls", "-la", "/mountedtmp"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, mntdir)
		log.Println(mntdir)
	})
	t.Run("Test_host_directory_path", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_host_directory_path")
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
	})
	t.Run("Test_mounted_temp_directory", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_mounted_temp_directory")
		require.NotNil(t, c)

		mntemp, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithExec([]string{"ls", "/"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, mntemp, "mountedtmp")

		mntdir, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithMountedDirectory("/mountedtmp", c.Host().Directory(".")).
			WithExec([]string{"ls", "/mountedtmp"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, mntdir)
	})
	t.Run("Test_mounted_host_directory", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_mounted_host_directory")
		require.NotNil(t, c)

		container := util.MountedHostDirectory(c, id, ".", "/mountedtmp")
		require.NotNil(t, container)

		mntdir, err := c.
			Container(dagger.ContainerOpts{ID: id}).
			WithMountedDirectory("/mountedtmp", c.Host().Directory(".")).
			WithExec([]string{"ls", "-la", "/mountedtmp"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, mntdir)
	})
}
