package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/alecharmon/codeowners-cli/core"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var head string
var base string
var newFilesOnly bool
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
		co := core.MustGetCodeOwners(file)
		fmt.Printf("Initializing repo at `%s` \n", dir)
		fmt.Printf("Running coverage check on %s \n", file)
		if head != "head" {
			fmt.Printf("Between %s...%s \n", base, head)

		}

		files, err := core.Diff(dir, head, base, newFilesOnly, core.NewLogger(verbose))
		if err != nil {
			fmt.Println("Could not load files from git repo")
			os.Exit(1)
		}

		//Providing uniqueness via a map
		without := make(map[string]bool)
		with := make(map[string]bool)
		allOwners := make(map[string]bool)

		for _, file := range files {
			owners := co.FindOwners(file)
			if len(owners) == 0 {
				color.Red("No Owner(s) found for %s \n", file)
				without[file] = true
			} else {
				with[file] = true
				for _, owner := range owners {
					allOwners[owner] = true
				}
			}
		}

		if len(without) > 0 {
			total := float64(len(with) + len(without))
			color.New(color.FgGreen).Printf("✅ %d files between commits that have defined maintainers (~%%%.2f) \n", len(with), (float64(len(with))/total)*100)
			color.New(color.FgRed).Printf("❌ %d files between commits are missing maintainers (~%%%.2f) \n", len(without), (float64(len(without))/total)*100)

		} else {
			color.New(color.FgGreen).Printf("🐣 Hell Ya, All files have declared maintainers\n")

		}
		if len(allOwners) > 0 {
			fmt.Printf("These are the declared maintainers for the analyzed files:\n")
			for owner := range allOwners {
				fmt.Println(owner)
			}
		}

		if len(without) > 0 {
			os.Exit(1)
		}
	},
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies that your code owner file is valid",
	Run: func(cmd *cobra.Command, args []string) {
		if file == "" {
			file = ".github/CODEOWNERS"
		}
		if dir == "" {
			dir = "./"
		}

		fmt.Printf("Initializing repo at `%s` \n", dir)
		core.MustGetCodeOwners(file)
		fmt.Println("Codeowner file initialized without any errors")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(verifyCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&head, "head", "head", "Latest commit to be used for tests")
	rootCmd.PersistentFlags().StringVar(&base, "base", "master", "Base commit to be used for tests")
	rootCmd.PersistentFlags().BoolVar(&newFilesOnly, "new_files_only", false, "Default behavior is to check for changed and new files, set to true to only check for introduced files")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Show verbose output")
	rootCmd.PersistentFlags().StringVar(&dir, "dir", viper.GetString("CODEOWNER_CI_DIRECTORY"), "Directory of the related project (default is PWD)")
	rootCmd.PersistentFlags().StringVar(&file, "file", viper.GetString("CODEOWNER_CI_FILE"), "CODEOWNER file to use for tests (default is PWD/CODEOWNER)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.codeowners-ci.yaml)")

	verifyCmd.PersistentFlags().StringVar(&dir, "dir", viper.GetString("CODEOWNER_CI_DIRECTORY"), "Directory of the related project (default is PWD)")
	verifyCmd.PersistentFlags().StringVar(&file, "file", viper.GetString("CODEOWNER_CI_FILE"), "CODEOWNER file to use for tests (default is PWD/CODEOWNER)")
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
