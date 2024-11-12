package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/hyperproperties/sopher/pkg/language"
)

var (
	sourceFlag string
)

func main() {
	flag.StringVar(&sourceFlag, "source", "", "the source file or directory")
	flag.Parse()

	// If the default source flag is used then we use working directory.
	if sourceFlag == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalln("Failed getting working directory as source flag", err)
		}
		sourceFlag = wd
	}

	files := language.NewFiles()
	err := files.Add(sourceFlag)
	if err != nil {
		log.Fatalln("Failed adding", sourceFlag, "to contracts", err)
	}

	injector := language.NewGoInjector()
	injector.Files(files.Iterator())

	// Run tests.
	time.Sleep(15 * time.Second)

	injector.Restore(files.Iterator())
}
