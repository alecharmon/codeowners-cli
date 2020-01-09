/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	codeowners "github.com/alecharmon/codeowners"
	"github.com/alecharmon/codeowners-cli/core"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var head string
var base string
var dir string
var file string
var verbose bool
var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "codeowners-cli",
	Short: "Determine coverage of explicit ownership and check ",
	Long: `Codeowners-cli is meant to be a light weight tool used to both validate the CODEOWNER file as well as 
	establish coverage of which files declared to have ownership`,

	Run: func(cmd *cobra.Command, args []string) {
		if file == "" {
			file = "./CODEOWNERS"
		}
		if dir == "" {
			dir = "./"
		}
		co, err := codeowners.BuildFromFile(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Running coverage check on %s \n", file)
		fmt.Printf("initializing repo at `%s` \n", dir)
		if head != "head" {
			fmt.Printf("Between %s...%s \n", base, head)

		}

		files, err := core.Diff(dir, head, base)
		if err != nil {
			fmt.Println("Could not load files from git repo")
		}

		//Providing uniqueness via a map
		without := make(map[string]bool)
		with := make(map[string]bool)

		for _, file := range files {
			owners := co.FindOwners(file)
			if len(owners) == 0 {
				color.Red("No Owner(s) found for %s \n", file)
				without[file] = true
			} else {
				with[file] = true
			}
		}

		total := float64(len(with) + len(without))
		color.New(color.FgGreen).Printf("✅ %d files between commits that have defined code owners (~%%%.2f) \n", len(with), float64(len(with))/total)
		color.New(color.FgRed).Printf("❌ %d files between commits are missing defined code owners (~%%%.2f) \n", len(without), float64(len(without))/total)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&head, "head", "head", "Latest commit to be used for tests")
	rootCmd.PersistentFlags().StringVar(&base, "base", "master", "Base commit to be used for tests")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Show verbose output")
	rootCmd.PersistentFlags().StringVar(&dir, "dir", viper.GetString("CODEOWNER_CI_DIRECTORY"), "Directory of the related project (default is PWD)")
	rootCmd.PersistentFlags().StringVar(&file, "file", viper.GetString("CODEOWNER_CI_FILE"), "CODEOWNER file to use for tests (default is PWD/CODEOWNER)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.codeowners-ci.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".codeowners-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".codeowners-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
