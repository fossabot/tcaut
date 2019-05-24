package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var scanpath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tcaut",
	Short: "TCAUT / Ripgrep for code auditing",
	Long:  ``,
}

// Execute called by main
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "rules.yaml", "config file (rules.yaml)")
	rootCmd.PersistentFlags().StringVarP(&scanpath, "target", "t", ".", "scanning target path")
	rootCmd.PersistentFlags().BoolP("detail", "d", false, "detailed output")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {

		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("yaml")

	} else {
		os.Exit(1)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("|")
		fmt.Println("| Loading policy file :", viper.ConfigFileUsed())
	}
}
