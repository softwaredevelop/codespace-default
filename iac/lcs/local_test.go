package main_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/stretchr/testify/require"
)

func TestLocalProject(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("Test_new_stack_local_source_secrets", func(t *testing.T) {
		t.Parallel()

		stackName := "testLocalSourceSecrets"
		workDir := filepath.Join(".", "localproject")
		s, err := auto.NewStackLocalSource(ctx, stackName, workDir)
		require.NoError(t, err)
		require.NotNil(t, s)

		defer func() {
			err := os.Unsetenv("FOO")
			require.NoError(t, err)

			err = s.Workspace().RemoveStack(ctx, s.Name())
			require.NoError(t, err)
		}()

		f := os.Getenv("FOO")
		b := os.Getenv("BAZ")
		s.Workspace().SetEnvVars(map[string]string{
			"FOO": f,
			"BAZ": b,
		})

		envvars := s.Workspace().GetEnvVars()
		require.Equal(t, f, envvars["FOO"])
		require.Equal(t, b, envvars["BAZ"])

		b = os.Getenv("BAR_TOKEN")
		err = s.SetAllConfig(ctx, auto.ConfigMap{
			"bar:token": auto.ConfigValue{
				Value:  b,
				Secret: true,
			},
			"buzz:owner": auto.ConfigValue{
				Value:  "xyz",
				Secret: true,
			},
		})
		require.NoError(t, err)

		values, err := s.GetAllConfig(ctx)
		require.NoError(t, err)

		require.Equal(t, b, values["bar:token"].Value)
		require.True(t, values["bar:token"].Secret)

		require.Equal(t, "xyz", values["buzz:owner"].Value)
		require.True(t, values["buzz:owner"].Secret)
	})
	t.Run("Test_new_stack_local_source_config", func(t *testing.T) {
		t.Parallel()

		stackName := "testLocalSourceConfig"
		workDir := filepath.Join(".", "localproject")
		s, err := auto.NewStackLocalSource(ctx, stackName, workDir)
		require.NoError(t, err)
		require.NotNil(t, s)

		defer func() {
			err = s.Workspace().RemoveStack(ctx, s.Name())
			require.NoError(t, err)
		}()

		require.Equal(t, stackName, s.Name())

		err = s.SetAllConfig(ctx, auto.ConfigMap{
			"bar:token": auto.ConfigValue{
				Value:  "abc",
				Secret: true,
			},
			"buzz:owner": auto.ConfigValue{
				Value:  "xyz",
				Secret: true,
			},
		})
		require.NoError(t, err)

		values, err := s.GetAllConfig(ctx)
		require.NoError(t, err)

		require.Equal(t, "abc", values["bar:token"].Value)
		require.True(t, values["bar:token"].Secret)

		require.Equal(t, "xyz", values["buzz:owner"].Value)
		require.True(t, values["buzz:owner"].Secret)
	})
	t.Run("Test_local_source_workspace_env_vars", func(t *testing.T) {
		t.Parallel()

		stackName := "testLocalSourceWorkspaceEnvVars"
		workDir := filepath.Join(".", "localproject")
		s, err := auto.NewStackLocalSource(ctx, stackName, workDir)
		require.NoError(t, err)
		require.NotNil(t, s)

		defer func() {
			err = s.Workspace().RemoveStack(ctx, s.Name())
			require.NoError(t, err)
		}()

		require.Equal(t, stackName, s.Name())

		err = s.Workspace().SetEnvVars(map[string]string{
			"FOO": "BAR",
			"BAZ": "QUX",
		})
		require.NoError(t, err)

		envvars := s.Workspace().GetEnvVars()
		require.Equal(t, "BAR", envvars["FOO"])
		require.Equal(t, "QUX", envvars["BAZ"])

		s.Workspace().UnsetEnvVar("FOO")
		s.Workspace().UnsetEnvVar("BAZ")

		envvars = s.Workspace().GetEnvVars()
		require.NotContains(t, envvars, "FOO")
		require.NotContains(t, envvars, "BAZ")
	})
}
