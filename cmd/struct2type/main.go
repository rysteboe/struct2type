package main

import (
	"fmt"
	"os"

	"github.com/rysteboe/struct2type/pkg/converter"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "struct2type",
	Short: "Convert Go structs to TypeScript types",
	Long: `A command-line tool that converts Go structs to TypeScript types.
It supports basic Go types, arrays, maps, and structs with JSON tags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please provide a Go source file path")
		}

		inputFile := args[0]
		outputFile := cmd.Flag("output").Value.String()

		conv := converter.New()
		tsTypes, err := conv.ConvertFile(inputFile)
		if err != nil {
			return fmt.Errorf("conversion failed: %w", err)
		}

		if outputFile == "" {
			fmt.Println(tsTypes)
			return nil
		}

		return os.WriteFile(outputFile, []byte(tsTypes), 0644)
	},
}

func init() {
	rootCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
