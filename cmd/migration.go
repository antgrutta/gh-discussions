/*
Copyright © 2022 Anthony Grutta antgrutta@github.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/antgrutta/gh-discussions/pkg/migration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrationCmd represents the migration command
var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Import discussions from a csv file.",
	Long:  `Import discussions from a csv file.`,
	Run: func(cmd *cobra.Command, args []string) {
		dataFile := cmd.Flag("data-file").Value.String()
		repoName := cmd.Flag("repo").Value.String()
		importType := cmd.Flag("type").Value.String()
		token := cmd.Flag("access-token").Value.String()
		logfile := cmd.Flag("logfile").Value.String()

		// Set the access token
		os.Setenv("GHD_IMPORT-TOKEN", token)
		os.Setenv("GHD_LOGFILE", logfile)
		viper.BindEnv("IMPORT-TOKEN")
		viper.BindEnv("LOGFILE")

		switch importType {
		case "stackoverflow":
			migration.ImportData(dataFile, repoName, importType)
		default:
			fmt.Println("Invalid import type, valid types are: stackoverflow")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	migrationCmd.Flags().StringP("type", "t", "", "Type of import, valid types are: stackoverflow")
	migrationCmd.Flags().StringP("data-file", "d", "", "CSV file to import")
	migrationCmd.Flags().StringP("repo", "r", "", "Repository to import discussions into (owner/repo)")
	migrationCmd.MarkFlagRequired("type")
	migrationCmd.MarkFlagRequired("data-file")
	migrationCmd.MarkFlagRequired("repo")
}
