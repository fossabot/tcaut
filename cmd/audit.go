package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	ffmt "gopkg.in/ffmt.v1"
)

// AllRules : internal struct of the rules yaml file
type allRules struct {
	Rules []struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Solution    string   `yaml:"solution"`
		Environment string   `yaml:"environment"`
		Fatal       bool     `yaml:"fatal"`
		Patterns    []string `yaml:"patterns"`
	} `yaml:"Rules"`
}

// OpaOutput : opa output schema
type opaOutput struct {
	Nonprod      int  `json:"Nonprod"`
	Nonprodfatal bool `json:"Nonprodfatal"`
	Prod         int  `json:"Prod"`
	Prodfatal    bool `json:"Prodfatal"`
}

var (
	rules     *allRules
	opaReport *opaOutput
)

func getRuleStruct() *allRules {

	err := viper.Unmarshal(&rules)
	if err != nil {
		fmt.Printf("| Unable to decode into rule config struct, %v", err)
	}
	return rules

}

// auditCmd
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audits your code against your defined rules",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		opaReport := opaOutput{
			Nonprod:      0,
			Nonprodfatal: false,
			Prod:         0,
			Prodfatal:    false,
		}

		rules = getRuleStruct()

		rgbin := "rg"
		path, err := exec.LookPath("rg")

		if err != nil {

			switch runtime.GOOS {
			case "windows":
				rgbin = "rg/rg.exe"
			case "darwin":
				rgbin = "rg/rgm"
			case "linux":
				rgbin = "rg/rgl"
			default:
				log.Fatalln("| OS not supported")
			}

			fmt.Println("| ripgrep not available in PATH, using local binary")

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
				fmt.Printf("| Rule #%d Pattern #%d : %s\n", index, pindex, pvalue)

				codePattern := []string{"-p", "-i", "-C2", "-U", pvalue, scanpath}
				xcmd := exec.Command(rgbin, codePattern...)

				xcmd.Stdout = os.Stdout
				xcmd.Stderr = os.Stderr

				errr := xcmd.Run()

				if errr != nil {
					if xcmd.ProcessState.ExitCode() == 2 {
						fmt.Println("| Check regex pattern syntax")
						log.Fatal(errr)
					} else {
						fmt.Println("| Clean")
					}
				} else {
					if value.Environment == "non-prod" {
						opaReport.Nonprod++
						if value.Fatal {
							fmt.Println("|")
							fmt.Println("| This violation blocks your code promotion between environments")
							opaReport.Nonprodfatal = true

						}
					} else {
						opaReport.Prod++
						if value.Fatal {
							fmt.Println("|")
							fmt.Println("| This violation is fatal for prod environments")
							opaReport.Prodfatal = true
						}
					}
					fmt.Println("|")
					fmt.Println("||", value.Name)
					fmt.Println("|| Solution : ", value.Solution)
					fmt.Println("|")
				}
			}
		}

		file, _ := json.MarshalIndent(opaReport, "", " ")
		_ = ioutil.WriteFile("opa-report.json", file, 0644)

		fmt.Println("|")
		fmt.Println("|")
		fmt.Println("| OPA REPORT")
		ffmt.Puts(opaReport)
		fmt.Println("|")
		fmt.Println("| EXIT")

	},
}

func init() {

	rootCmd.AddCommand(auditCmd)

}
