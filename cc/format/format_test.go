package format_test

import (
	"cc/format"
	"context"
	"flag"
	"os"
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

	c = c.Pipeline("format_test")

	id, _ = c.
		Container().
		ID(ctx)

	code := m.Run()
	defer c.Close()
	os.Exit(code)
}

func TestGoFormat(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_gofumpt", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_gofumpt")
		require.NotNil(t, c)

		_, err := format.GoFormat(c, id).
			WithExec([]string{"/gofumpt", "-version"}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_goimports", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_goimports")
		require.NotNil(t, c)

		_, err := format.GoFormat(c, id).
			WithExec([]string{"/goimports", "-h"}).
			Stdout(ctx)
		require.Error(t, err)
	})
	t.Run("Test_container_ID_busybox", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_container_ID_busybox")
		require.NotNil(t, c)

		id, err := c.Container().
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
		c = c.Pipeline("test_container_ID")
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
}
