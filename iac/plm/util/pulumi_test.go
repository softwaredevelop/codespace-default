package util_test

import (
	"context"
	"flag"
	"os"
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

	c = c.Pipeline("pulumi_util_test")

	id, _ = c.
		Container().
		From("busybox:uclibc").
		ID(ctx)

	code := m.Run()
	defer c.Close()
	os.Exit(code)
}

func TestPulumi(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_pulumi_login_local", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_pulumi_login_local")
		require.NotNil(t, c)

		container := util.PulumiInstall(c, id)
		require.NotNil(t, container)
		out, err := container.
			WithExec([]string{"pulumi", "login", "--local"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "Logged in")
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
}
