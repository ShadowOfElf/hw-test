package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"BAR":   EnvValue{"bar", false},
		"EMPTY": EnvValue{"", true},
		"FOO":   EnvValue{"   foo\nwith new line", false},
		"HELLO": EnvValue{"\"hello\"", false},
		"UNSET": EnvValue{"", true},
	}
	t.Run("correct_run", func(t *testing.T) {
		args := []string{"/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2"}
		code := RunCmd(args, env)
		require.Equal(t, 0, code)
	})
}
