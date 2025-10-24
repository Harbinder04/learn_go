package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var subCmd = &cobra.Command{
	Use:     "subtract",
	Aliases: []string{"substract"},
	Short:   "Substract 2 numbers",
	Long:    "Carry out substraction operation on 2 numbers",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Substraction of %s and %s =  %s. \n\n", args[0], args[1], Subtract(args[0], args[1]))
	},
}

func init() {
	rootCmd.AddCommand(subCmd)
}
