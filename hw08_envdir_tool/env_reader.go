package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	result := make(Environment)
	envDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entity := range envDir {
		if !entity.IsDir() {
			path := filepath.Join(dir, entity.Name())
			envValue, err := getEnvFromFile(path)
			if err != nil {
				return nil, err
			}
			result[entity.Name()] = envValue
		}
	}

	return result, nil
}

func getEnvFromFile(path string) (EnvValue, error) {
	resultValue := EnvValue{
		Value:      "",
		NeedRemove: true,
	}

	src, err := os.Open(path)
	if err != nil {
		return resultValue, err
	}
	defer func() {
		err = src.Close()
		if err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(src)
	if scanner.Scan() {
		text := strings.TrimRight(scanner.Text(), " ")
		text = string(bytes.ReplaceAll([]byte(text), []byte{0x00}, []byte("\n")))
		if text != "" {
			resultValue.Value = text
			resultValue.NeedRemove = false
		}
		return resultValue, nil
	}

	return resultValue, scanner.Err()
}
