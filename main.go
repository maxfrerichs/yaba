package main

import (
	"fmt"
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Build struct {
		Input struct {
			Directory   string   `yaml:"directory"`
			Entrypoints []string `yaml:"entrypoints"`
		} `yaml:"input"`
		Output struct {
			Directory string `yaml:"directory"`
			File      string `yaml:"file"`
		} `yaml:"output"`
		Options struct {
			Minify      bool `yaml:"minify"`
			Sourcemap   bool `yaml:"sourcemap"`
			TreeShaking bool `yaml:"treeShaking"`
		} `yaml:"options"`
	} `yaml:"build"`
}

type Command struct {
	name        string
	description string
	callback    func() error
}

func getCommands() map[string]Command {
	return map[string]Command{
		"build": {
			name:        "build",
			description: "Executes the build",
			callback:    commandBuild,
		},
		"help": {
			name:        "help",
			description: "See available commands",
			callback:    commandHelp,
		},
	}
}

func commandHelp() error {
	fmt.Println("Usage:")
	commands := getCommands()
	for _, command := range commands {
		fmt.Println(command.name + " - " + command.description)
	}
	return nil
}

func commandBuild() error {
	file, err := os.ReadFile("build.yaml")
	if err != nil {
		fmt.Println("Error reading config file")
		return err
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error reading config file")
		os.Exit(1)
	}

	esbuildOptions := api.BuildOptions{
		EntryPoints:      getEntrypointPaths(config),
		Bundle:           true,
		Sourcemap:        api.SourceMapNone,
		AssetNames:       "[name]-[hash]",
		MinifyWhitespace: config.Build.Options.Minify,
		Write:            true,
	}

	if config.Build.Options.TreeShaking {
		esbuildOptions.TreeShaking = api.TreeShakingTrue
	}

	if len(config.Build.Input.Entrypoints) < 2 {
		esbuildOptions.Outfile = config.Build.Output.Directory + "/" + config.Build.Output.File
	} else {
		esbuildOptions.Outdir = config.Build.Output.Directory
		esbuildOptions.Outfile = ""
	}

	if config.Build.Options.Sourcemap {
		esbuildOptions.Sourcemap = api.SourceMapLinked
	}

	if _, err := os.Stat(config.Build.Output.Directory); os.IsNotExist(err) {
		os.MkdirAll(config.Build.Output.Directory, os.ModePerm)
	}
	os.RemoveAll(config.Build.Output.Directory)

	fmt.Println("Starting build...")
	result := api.Build(esbuildOptions)

	if len(result.Errors) > 0 {
		for _, err := range result.Errors {
			fmt.Println(err.Text)
		}
		os.Exit(1)
	}

	fmt.Println("Build completed successfully!")
	return nil
}

func main() {
	if len(os.Args) < 2 {
		cmd := getCommands()["help"]
		cmd.callback()
		os.Exit(0)
	}
	cmd, exists := getCommands()[os.Args[1]]
	if !exists {
		fmt.Println("Command not found")
		os.Exit(0)
	}
	cmd.callback()
}

func getEntrypointPaths(config Config) []string {
	var entrypointPaths []string
	for _, entrypoint := range config.Build.Input.Entrypoints {
		entrypointPath := config.Build.Input.Directory + "/" + entrypoint
		entrypointPaths = append(entrypointPaths, entrypointPath)
	}
	return entrypointPaths
}
