/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const (
	GO_GIT_DIR = ".go-git/"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Args:  cobra.MinimumNArgs(1),
	Short: "Stage files specified",
	Long: `Stage the given flies for commit
	only supports simle staging, and files must be specified by name
	'.' or '*' is not supported
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		// prerequisites
		err := os.Chdir(GO_GIT_DIR)
		if os.IsNotExist(err) {
			log.Fatalf("could not find the .go-git directory; call `go-git init`. error: %v", err)
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
