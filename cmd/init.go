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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		var directories [3]string = [3]string{
			".go-git/",
			".go-git/objects/",
			".go-git/refs/",
		}
		var files [6]string = [6]string{
			".go-git/HEAD",
			".go-git/config",
			".go-git/objects/info",
			".go-git/objects/pack",
			".go-git/refs/heads",
			".go-git/refs/tags",
		}
		for _, d := range directories {
			os.Mkdir(d, os.ModePerm)
		}
		for _, fname := range files {
			f, err := os.Create(fname)
			if err != nil {
				log.Fatal(err)
			}
			if fname == ".go-git/HEAD" {
				f.Write([]byte("ref: refs/heads/master"))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
