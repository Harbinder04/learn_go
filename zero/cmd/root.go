package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "zero",
	Aliases: []string{"zero", "Zero"},
	Short:   "Zero is a CLI tool for doing simple calculations",
	Long:    `Zero is a CLI tool that allows you to perform simple calculations like addition, subtraction, multiplication, and division.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Zero CLI!")
	},
}

func Execution() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops something wend wrong: %s", err)
		os.Exit(1)
	}
}
