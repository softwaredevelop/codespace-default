package linting_test

import (
	"cc/linting"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoLint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_revive", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_revive")
		require.NotNil(t, c)

		out, err := linting.ReviveL(c, id).
			WithExec([]string{"/revive", "-version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "version")
	})
	t.Run("Test_golint", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_golint")
		require.NotNil(t, c)

		_, err := linting.GoLint(c, id).
			WithExec([]string{"/golint", "-help"}).
			Stdout(ctx)
		require.NoError(t, err)
	})
	t.Run("Test_fgt", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_fgt")
		require.NotNil(t, c)

		_, err := linting.GoLint(c, id).
			WithExec([]string{"/fgt"}).
			Stdout(ctx)
		require.NoError(t, err)
	})
}
