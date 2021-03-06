package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/burbokop/design-practice-1/build/modules/gomodule"
	"github.com/burbokop/design-practice-1/build/modules/zip_archive"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	dryRun  = flag.Bool("dry-run", false, "Generate ninja build file but don't start the build")
	verbose = flag.Bool("v", false, "Display debugging logs")
	task    = flag.String("task", "", "Run specific task")
)

func NewContext(allovedTasks []string) *blueprint.Context {
	ctx := bood.PrepareContext()
	ctx.RegisterModuleType("go_binary", gomodule.CreateSimpleBinFactory(allovedTasks))
	ctx.RegisterModuleType("zip_archive", zip_archive.SimpleBinFactory)
	return ctx
}

func main() {
	flag.Parse()

	config := bood.NewConfig()
	if !*verbose {
		config.Debug = log.New(ioutil.Discard, "", 0)
	}

	tasks := []string{}
	if *task != "" {
		tasks = append(tasks, *task)
	}

	ctx := NewContext(tasks)

	ninjaBuildPath := bood.GenerateBuildFile(config, ctx)

	if !*dryRun {
		config.Info.Println("Starting the build now")
		cmd := exec.Command("ninja", append([]string{"-f", ninjaBuildPath}, flag.Args()...)...)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			config.Info.Fatal("Error invoking ninja build. See logs above.")
		}
	}
}
