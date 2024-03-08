package linting_test

import (
	"cc/linting"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActionlint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_actionlint_version", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_actionlint_version")
		require.NotNil(t, c)

		out, err := linting.Al(c, id).
			WithExec([]string{"/actionlint", "-version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "v")
	})
	t.Run("Test_shellcheck_version", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_shellcheck_version")
		require.NotNil(t, c)

		out, err := linting.Al(c, id).
			WithExec([]string{"/shellcheck", "--version"}).
			Stdout(ctx)
		require.NoError(t, err)
		require.Contains(t, out, "v")
	})
}
