package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ffmt "gopkg.in/ffmt.v1"
)

// AllRules : internal struct of the rules yaml file
type AllRules struct {
	Rules []struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Solution    string   `yaml:"solution"`
		Environment string   `yaml:"environment"`
		Fatal       bool     `yaml:"fatal"`
		Patterns    []string `yaml:"patterns"`
	} `yaml:"Rules"`
}

var (
	rules *AllRules
)

func getConfStruct() *AllRules {

	err := viper.Unmarshal(&rules)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
	return rules

}

// auditCmd
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audits your code against your defined rules",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		rules = getConfStruct()

		rgbin := "rg"

		path, err := exec.LookPath("rg")

		if err != nil {

			if runtime.GOOS == "windows" {
				rgbin = "rg/rg.exe"
			}
			if runtime.GOOS == "darwin" {
				rgbin = "rg/rgm"
			}
			if runtime.GOOS == "linux" {
				rgbin = "rg/rgl"
			}

			fmt.Println("| rg not available in PATH, using local binary")

		}

		detail, _ := rootCmd.PersistentFlags().GetBool("detail")
		if detail {

			fmt.Printf("| rg is available at %s\n", path)
			fmt.Println("|")
			ffmt.Puts(rules)
			fmt.Println("|")
		}

		for index, value := range rules.Rules {

			fmt.Println("| ------------------------------------------------------------")
			fmt.Println("| Rule #", index)
			fmt.Println("| Rule name : ", value.Name)

			for pindex, pvalue := range value.Patterns {

				fmt.Println("| ----------")
				fmt.Println("| Pattern #", pindex)
				fmt.Println("| Pattern : ", pvalue)

				codePattern := []string{"-p", "-i", "-C2", "-U", pvalue, scanpath}
				xcmd := exec.Command(rgbin, codePattern...)

				xcmd.Stdout = os.Stdout
				xcmd.Stderr = os.Stderr

				errr := xcmd.Run()

				if errr != nil {
					if xcmd.ProcessState.ExitCode() == 2 {
						log.Fatal(errr)
					} else {
						fmt.Println("| Clean")
					}

				} else {
					if value.Fatal && value.Environment == "prod" {
						fmt.Println("|")
						fmt.Println("| This violation is fatal for non-prod environments")
					}
					if value.Fatal && value.Environment == "non-prod" {
						fmt.Println("|")
						fmt.Println("| This violation blocks your code promotion between environments")
					}
					fmt.Println("|")
					fmt.Println("||", value.Name)
					fmt.Println("|| Solution : ", value.Solution)
					fmt.Println("|")
				}

			}
		}

	},
}

func init() {

	rootCmd.AddCommand(auditCmd)

}
