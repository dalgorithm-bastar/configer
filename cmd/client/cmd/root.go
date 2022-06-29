package cmd

import (
    "github.com/configcenter/config"
    manage "github.com/configcenter/pkg/service"
    "github.com/spf13/cobra"

    "github.com/spf13/viper"
)

// Object 用于接收参数的公用结构体，不同指令下初始化不同的变量
type Object struct {
    UserName       string
    Env            string
    Target         string
    Type           string
    Platform       string
    NodeType       string
    Version        string
    Scheme         string
    Set            string
    PathIn         string
    PathOut        string
    TopicIpRange   string
    TopicPortRange string
    TcpPortRange   string
}

var (
    GrpcInfo     manage.GrpcInfoStruct
    cfgFile      string
    object       Object
    Version      string
    GoVersion    string
    GitBranch    string
    GitCommit    string
    GitLatestTag string
    BuildTime    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "proxctl",
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
    rootCmd.PersistentFlags().StringVarP(&object.UserName, "user", "u", "", "current userName(required)")

    rootCmd.Flags().String(config.GrpcSocket, "", "set grpc socket")
    _ = viper.BindPFlag(config.GrpcSocket, rootCmd.Flag(config.GrpcSocket))

    rootCmd.Flags().Int(config.GrpcLockTimeout, 30, "set etcd lock timeout by second when post config")
    _ = viper.BindPFlag(config.GrpcLockTimeout, rootCmd.Flag(config.GrpcLockTimeout))
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

// GetGrpcClient 供多个命令调用，新建grpc连接
func GetGrpcClient() error {
    //读取grpc配置文件
    GrpcInfo = manage.GrpcInfoStruct{
        Socket:      viper.GetString(config.GrpcSocket),
        LockTimeout: viper.GetInt(config.GrpcLockTimeout),
    }
    return nil
}
