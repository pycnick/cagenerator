package main

import (
	"fmt"
	"os"

	"github.com/pycnick/cagenerator/internal/config"
	"github.com/pycnick/cagenerator/internal/generator"

	"github.com/spf13/cobra"
)

func main() {
	var configFile string

	rootCmd := &cobra.Command{
		Use:   "generate [flags] [project_path]",
		Short: "Generate project layers from YAML config",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runGenerate(&configFile),
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "entity.yaml", "Path to the config YAML file")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runGenerate(configFile *string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadConfig(*configFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		projectPath := "."
		if len(args) > 0 {
			projectPath = args[0]
		}

		gen, err := generator.New(cfg, projectPath)
		if err != nil {
			return fmt.Errorf("creating generator: %w", err)
		}

		if err := gen.Generate(); err != nil {
			return fmt.Errorf("generating project: %w", err)
		}

		fmt.Println("Project generated successfully!")
		return nil
	}
}
