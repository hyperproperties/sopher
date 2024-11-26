package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hyperproperties/sopher/pkg/language"
)

func main() {
	if len(os.Args) == 1 {
		initial()
		return
	}

	switch os.Args[1] {
	case "inject":
		inject()
	case "restore":
		restore()
	case "help":
		help()
	default:
		command := strings.Join(os.Args[1:], " ")
		fmt.Println("unknown command:", command)
	}
}

func initial() {
	fmt.Println("sopher is a automatic test generation and runtime verification tool for hypercontracts.")
	fmt.Println()
	help()
}

func help() {
	fmt.Println(`Usage:
	
		sopher <command> [arguments]

Commands:
	inject		injects the contracts into functions with contracts
	restore		restores injected source files to the original ones

Use "go help <command>" for more information about a command.`)
}

func files(path string) language.Files {
	if path == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalln("Failed getting working directory as source flag", err)
		}
		path = wd
	}

	files := language.NewFiles()
	err := files.Add(path)
	if err != nil {
		log.Fatalln("Failed adding", path, "to contracts", err)
	}

	return files
}

func inject() {
	flags := flag.NewFlagSet("inject", flag.ExitOnError)

	var path string
	flags.StringVar(&path, "path", "", "the path to go source files")

	flags.Parse(os.Args[2:])

	files := files(path)
	injector := language.NewGoInjector()
	injector.Files(files.Iterator())
}

func restore() {
	flags := flag.NewFlagSet("inject", flag.ExitOnError)

	var path string
	flags.StringVar(&path, "path", "", "the path to go source files")

	flags.Parse(os.Args[2:])

	files := files(path)
	injector := language.NewGoInjector()
	injector.Restore(files.Iterator())
}
