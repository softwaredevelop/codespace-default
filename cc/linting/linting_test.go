package linting_test

import (
	"context"
	"flag"
	"os"
	"testing"

	"dagger.io/dagger"
)

var c *dagger.Client
var id dagger.ContainerID

func TestMain(m *testing.M) {
	flag.Parse()

	ctx := context.Background()

	c, _ = dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	defer c.Close()

	c = c.Pipeline("linting_test")

	id, _ = c.
		Container().
		WithMountedTemp("/mountedtmp").
		ID(ctx)

	code := m.Run()
	defer c.Close()
	os.Exit(code)
}
