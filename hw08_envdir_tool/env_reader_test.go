package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("correct read", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		expectedBar := EnvValue{Value: "bar", NeedRemove: false}
		expectedFoo := EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		require.NoError(t, err)
		require.Equal(t, expectedBar, env["BAR"])
		require.Equal(t, expectedFoo, env["FOO"])
		require.Len(t, env, 5)
	})

	t.Run("empty_dir", func(t *testing.T) {
		_ = os.Mkdir("testdata/empty_dir", 0o644)
		defer func() {
			_ = os.Remove("testdata/empty_dir")
		}()
		env, err := ReadDir("./testdata/empty_dir")
		require.NoError(t, err)
		require.Len(t, env, 0)
	})

	t.Run("dir_without_rights", func(t *testing.T) {
		err := os.Mkdir("testdata/test_dir", 0)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		defer func() {
			_ = os.Chmod("testdata/test_dir", 0o644)
			_ = os.Remove("testdata/test_dir")
		}()
		_, err = ReadDir("./testdata/test_dir")
		require.Error(t, err)
	})

	t.Run("file_without_rights", func(t *testing.T) {
		err := os.WriteFile("./testdata/env/test_file.txt", []byte{}, 0)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		defer func() {
			_ = os.Chmod("testdata/env/test_file.txt", 0o644)
			_ = os.Remove("testdata/env/test_file.txt")
		}()
		_, err = ReadDir("./testdata/env")

		require.Error(t, err)
	})
}
