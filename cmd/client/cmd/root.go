/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"github.com/spf13/cobra"
	"os"

	"github.com/spf13/viper"
)

//配置中心服务端配置文件存放的路径，应与服务端main.go同步修改
const (
	logConfigLocation  = "config/log.json"
	grpcConfigLocation = "config/grpc.json"
	//etcdConfigLocation = "config/etcdClientv3.json"
)

// Object 用于接收参数的公用结构体，不同指令下初始化不同的变量
type Object struct {
	UserName     string
	Target       string
	Phrase       string
	Version      string
	Env          string
	Cluster      string
	GlobalId     string
	LocalId      string
	TemplateName string
	PathIn       string
	PathOut      string
}

var (
	cfgFile string
	object  Object
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cfgtool",
	Short: "This is the commandline tool for configcenter client",
	Long: `	This is the commandline tool for configcenter client
		The tool enables you to:
		create configurefile locally
		get servicelist or publicinfofile or template from remote
		find particular info from remote by go template function`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cfgtool.yaml)")
	rootCmd.PersistentFlags().StringVarP(&object.UserName, "user", "u", "chqr", "current userName(required)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.MarkPersistentFlagRequired("user")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cfgtool" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cfgtool")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
