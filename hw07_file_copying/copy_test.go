package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tmpFile, _ := os.Create("tmpFile.txt")
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	t.Run("correct copy", func(t *testing.T) {
		err := Copy("testdata/input.txt", tmpFile.Name(), 0, 0)
		templateFile, _ := os.ReadFile("testdata/out_offset0_limit0.txt")
		resultFile, _ := os.ReadFile("tmpFile.txt")

		require.NoError(t, err)
		require.Equal(t, string(templateFile), string(resultFile))
	})

	t.Run("offset exceeds", func(t *testing.T) {
		err := Copy("testdata/input.txt", tmpFile.Name(), 100000000000, 0)

		require.Error(t, err)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("wrong limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", tmpFile.Name(), 0, -1)
		templateFile, _ := os.ReadFile("testdata/out_offset0_limit0.txt")
		resultFile, _ := os.ReadFile("tmpFile.txt")

		require.NoError(t, err)
		require.Equal(t, string(templateFile), string(resultFile))
	})

	t.Run("directory copy", func(t *testing.T) {
		err := Copy("testdata/", tmpFile.Name(), 0, 0)

		require.Error(t, err)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("limit and offset copy", func(t *testing.T) {
		err := Copy("testdata/input.txt", tmpFile.Name(), 100, 1000)
		templateFile, _ := os.ReadFile("testdata/out_offset100_limit1000.txt")
		resultFile, _ := os.ReadFile("tmpFile.txt")

		require.NoError(t, err)
		require.Equal(t, string(templateFile), string(resultFile))
	})

	t.Run("src error", func(t *testing.T) {
		err := Copy("testdata/inp.txt", tmpFile.Name(), 0, 0)
		require.Error(t, err)
	})

	t.Run("dst not empty", func(t *testing.T) {
		_ = Copy("testdata/input.txt", tmpFile.Name(), 0, 0)
		err := Copy("testdata/input.txt", tmpFile.Name(), 0, 0)
		templateFile, _ := os.ReadFile("testdata/out_offset0_limit0.txt")
		resultFile, _ := os.ReadFile("tmpFile.txt")

		require.NoError(t, err)
		require.Equal(t, string(templateFile), string(resultFile))
	})

	t.Run("src stat error", func(t *testing.T) {
		fileName := "acc_wrong.txt"
		err := os.WriteFile(fileName, []byte{}, 0)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		defer func() {
			_ = os.Chmod(fileName, 0o644)
		}()
		defer func() {
			_ = os.Remove(fileName)
		}()

		err = Copy(fileName, tmpFile.Name(), 0, 0)
		require.Error(t, err)
	})

	t.Run("copy error", func(t *testing.T) {
		err := os.Mkdir("testdata/testdir", 0)
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		defer func() {
			_ = os.Chmod("testdata/testdir", 0o644)
		}()
		defer func() {
			_ = os.Remove("testdata/testdir")
		}()
		err = Copy("testdata/input.txt", "testdata/testdir/test.txt", 0, 0)
		_, _ = os.ReadFile("testdata/out_offset0_limit0.txt")
		_, _ = os.ReadFile("tmpFile.txt")

		require.Error(t, err)
	})

	t.Run("src not created", func(t *testing.T) {
		err := Copy("nil_file.txt", tmpFile.Name(), 0, 0)
		require.Error(t, err)
	})
}
