package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	cmd := []string{"echo", "SOME TEXT"}

	customEnvName := "CUSTOM_ENV"

	t.Run("change env value", func(t *testing.T) {
		os.Setenv(customEnvName, "old_value")

		returnCode := RunCmd(cmd, Environment{
			customEnvName: EnvValue{
				Value:      "new_value",
				NeedRemove: false,
			},
		})

		require.Equal(t, "new_value", os.Getenv(customEnvName))
		require.Equal(t, 0, returnCode)
	})

	t.Run("create env", func(t *testing.T) {
		returnCode := RunCmd(cmd, Environment{
			customEnvName: EnvValue{
				Value:      "new_value",
				NeedRemove: false,
			},
		})

		value := os.Getenv(customEnvName)

		require.Equal(t, "new_value", value)
		require.Equal(t, 0, returnCode)
	})

	t.Run("delete env", func(t *testing.T) {
		os.Setenv(customEnvName, "old_value")

		returnCode := RunCmd(cmd, Environment{
			customEnvName: EnvValue{
				Value:      "",
				NeedRemove: true,
			},
		})

		_, isPresent := os.LookupEnv(customEnvName)

		require.False(t, isPresent)
		require.Equal(t, 0, returnCode)
	})
}
