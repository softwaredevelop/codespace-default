package linting_test

import (
	"cc/linting"
	"context"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestActionlint(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_actionlint_run", func(t *testing.T) {
		t.Parallel()
		c = c.Pipeline("test_actionlint_run")
		require.NotNil(t, c)

		_, err := linting.Al(c, id).
			WithNewFile("test.yml", dagger.ContainerWithNewFileOpts{
				Contents: `{name: "Test workflow", on: {push: {branches: ["main"]}}, jobs: {build: {runs-on: "ubuntu-latest", steps: [{name: "Checkout code", uses: "actions/checkout@v2"}, {name: "Run a one-line script", run: "echo Hello, world!"}]}}}`,
			}).
			WithNewFile("test2.yml", dagger.ContainerWithNewFileOpts{
				Contents: `{name: "Test2 workflow", on: {push: {branches: ["main"]}}}`,
			}).
			WithExec([]string{"sh", "-c", "/actionlint $(find . -type f -name '*.yml')"}).
			Stdout(ctx)
		require.ErrorContains(t, err, "is missing")
	})
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
