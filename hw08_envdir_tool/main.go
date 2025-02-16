package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("Недостаточно аргументов")
		return
	}

	env, err := ReadDir(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	exitCode := RunCmd(args[2:], env)

	os.Exit(exitCode)
}
