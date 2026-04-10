/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/miloskovacevic/mergeflow/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var appInstance *app.App

var companies = map[int]string{
	970:  "data-manager-go",
	700:  "insurance-api",
	827:  "investment-api",
	912:  "oracle-bridge-go",
	585:  "process-manager",
	696:  "result-job",
	1152: "storage-api",
	649:  "notifications",
	581:  "data-manager",
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mergeflow",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	viper.SetEnvPrefix("MERGEFLOW")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config") // config.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath(home + "/.config/mergeflow")

	if err = viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found, using env/flags only")
	}

	appInstance, err = app.NewApp()
	if err != nil {
		panic(err)
	}
}

type Config struct {
	Jira struct {
		Project string
		Route   string
		Token   string
	}
	Gitlab struct {
		Route string
		Repos []int
		Token string
	}
}
