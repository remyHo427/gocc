package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"cc260717/cmd/cc/backends/amd64"
	"cc260717/cmd/cc/backends/amd64/cg"
	"cc260717/cmd/cc/cfront/lex"
	"cc260717/cmd/cc/cfront/parse"
	"cc260717/cmd/cc/ir/tacky"
	"cc260717/cmd/cc/util"
)

type CmdFlags struct {
	Compile    bool
	Assemble   bool
	Preprocess bool
	Help       bool
	Output     string
}

func main() {
	bin_name := os.Args[0]
	flags := readFlags()
	source_files := []string{}

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			// skip flags
			continue
		} else if strings.HasSuffix(arg, ".c") {
			source_files = append(source_files, arg)
		}
	}

	if flags.Help {
		print_help(bin_name)
		flag.PrintDefaults()
	} else if len(source_files) == 0 {
		print_help(bin_name)
	} else {
		compile_files(flags, source_files)
	}
}

func readFlags() CmdFlags {
	compile := flag.Bool("c", false, "compile and assemble, don't link")
	assemble := flag.Bool("S", false, "compile but don't assemble or link")
	preprocess := flag.Bool("E", false, "stop after preprocessing")
	help := flag.Bool("help", false, "Print help messages")
	output := flag.String("o", "a.out", "file to output to")

	flag.Parse()

	return CmdFlags{
		Compile:    *compile,
		Assemble:   *assemble,
		Preprocess: *preprocess,
		Help:       *help,
		Output:     *output,
	}
}
func print_help(bin_name string) {
	var ccName string

	if strings.HasPrefix(bin_name, "./") {
		ccName = bin_name[2:]
	}

	fmt.Printf("Usage: %s [options]... file.... \n", ccName)
	fmt.Println("Compiles C source files")
}

func compile_files(flags CmdFlags, source_files []string) {
	for _, file := range source_files {
		compile_file(flags, file)
	}
}
func compile_file(flags CmdFlags, file string) {
	fname := file[:len(file)-2]
	tmp_asm := fmt.Sprintf("%s.s", fname)
	tmp_file := fmt.Sprintf("%s.i", fname)

	if _, err := os.Stat(file); err != nil {
		util.Exit_with_error(err)
	}

	// use gcc to preprocess
	cmd := exec.Command("gcc", "-E", "-P", file, "-o", tmp_file)
	if err := cmd.Run(); err != nil {
		util.Exit_with_error(err)
	}

	// stop at preprocessing
	if flags.Preprocess {
		return
	}

	src, err := os.ReadFile(tmp_file)
	if err != nil {
		util.Exit_with_error(err)
	}
	os.Remove(tmp_file)

	result := compile(string(src))
	f, err := os.Create(tmp_asm)
	if err != nil {
		util.Exit_with_error(err)
	}

	f.WriteString(result)
	f.WriteString("\n\n")

	defer f.Close()

	// don't assemble, emit assembly file
	if flags.Assemble {
		return
	}

	// use gcc to assemble and link
	if flags.Compile {
		cmd = exec.Command("gcc", "-c", tmp_asm)
	} else {
		cmd = exec.Command("gcc", tmp_asm, "-o", flags.Output)
	}

	if err := cmd.Run(); err != nil {
		util.Exit_with_error(err)
	}
	os.Remove(tmp_asm)
}

func compile(src string) string {
	l := lex.New(src)
	p := parse.New(l)
	ast := p.Parse()
	t := tacky.New()
	g := cg.New()

	return g.Generate(amd64.ToAsm(t.Generate(ast)))
}
