package main

import (
	"fmt"
	"os"
	"os/exec"
)

const BIN_PATH = "./bin"

var app_srcdir_map = map[string]string{
	"mycc": "./cmd/cc",
}

func main() {
	command := os.Args[1]
	rest := os.Args[2:]

	switch command {
	case "test":
		test(rest)
	case "build":
		build(rest)
	default:
		print_help()
	}
}

func test(args []string) {
	app_name := args[0]
	srcdir, ok := app_srcdir_map[app_name]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown app \"%s\"\n", app_name)
		os.Exit(1)
	}

	execute_command("go", "fmt", srcdir+"/...")
	execute_command("go", "vet", srcdir+"/...")
	execute_command("go", "test", srcdir+"/...")
}
func build(args []string) {
	app_name := args[0]
	srcdir, ok := app_srcdir_map[app_name]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown app \"%s\"\n", app_name)
		os.Exit(1)
	}

	execute_command("go", "generate", "./...")
	execute_command("go", "build", "-o", BIN_PATH+"/"+app_name, srcdir)
}

func print_help() {
	fmt.Printf("help message")
}

func execute_command(app string, args ...string) {
	cmd := exec.Command(app, args...)

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Print(string(stdout))
}
