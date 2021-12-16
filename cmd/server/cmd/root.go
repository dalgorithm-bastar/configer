package cmd

import (
    "github.com/spf13/cobra"

    "github.com/spf13/viper"
)

var (
    cfgFile      string
    Version      string
    GoVersion    string
    GitBranch    string
    GitCommit    string
    GitLatestTag string
    BuildTime    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "proxima",
    Short: "server of ConfigCenter",
    Long:  `server of ConfigCenter`,
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

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config/configcenter.json)")
    // Cobra also supports local flags, which will only run
    // when this action is called directly.
    //rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        // Find home directory.
        //home, err := os.UserHomeDir()
        //cobra.CheckErr(err)

        // Search config in home directory with name ".cfgsrv" (without extension).
        //viper.AddConfigPath("config")
        //viper.SetConfigType("json")
        //viper.SetConfigName("configcenter")
        viper.SetConfigFile("config/configcenter.json")
    }

    viper.AutomaticEnv() // read in environment variables that match

    // If a config file is found, read it in.
    if err := viper.ReadInConfig(); err == nil {
        //fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
    }
}
