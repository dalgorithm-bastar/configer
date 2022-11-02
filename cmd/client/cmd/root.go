package cmd

import (
	"fmt"
	"os"

	"github.com/configcenter/internal/log"
	"github.com/configcenter/pkg/define"
	manage "github.com/configcenter/pkg/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

const (
	_target    = "target"
	_type      = "type"
	_pathin    = "pathin"
	_version   = "version"
	_env       = "env"
	_scheme    = "scheme"
	_platform  = "platform"
	_nodeType  = "nodetype"
	_cluster   = "cluster"
	_pathout   = "pathout"
	_mode      = "mode"
	_topicIp   = "topicIp"
	_topicPort = "topicPort"
	_tcpPort   = "tcpPort"
	_ezeiInner = "ezeiInner"
	_ezeiEnv   = "ezeiEnv"
	_envCover  = "envCover"
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
	EzeiInner      string
	EzeiEnv        string
	EnvCover       bool
	mode           string
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
		get configdata or infrastructure from remote`,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ../config/client.yaml)")
	rootCmd.PersistentFlags().StringVarP(&object.UserName, "user", "u", "", "current userName(required)")

	rootCmd.Flags().String(define.GrpcSocket, "", "set grpc socket")
	_ = viper.BindPFlag(define.GrpcSocket, rootCmd.Flag(define.GrpcSocket))

	rootCmd.Flags().Int(define.GrpcLockTimeout, 30, "set etcd lock timeout by second when post config")
	_ = viper.BindPFlag(define.GrpcLockTimeout, rootCmd.Flag(define.GrpcLockTimeout))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//rootCmd.MarkPersistentFlagRequired("user")
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
		viper.SetConfigType("yaml")
		//viper.SetConfigName("configcenter")
		viper.SetConfigFile("../config/client.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		//fmt.Println(viper.AllSettings())
	} else {
		log.Logger.Warn("load cfg file failed, online function limited", zap.Error(err))
	}
}

// GetGrpcClient 供多个命令调用，新建grpc连接
func GetGrpcClient() error {
	//读取grpc配置文件
	GrpcInfo = manage.GrpcInfoStruct{
		Socket:      viper.GetString(define.GrpcSocket),
		LockTimeout: viper.GetInt(define.GrpcLockTimeout),
	}
	return nil
}

func transfermInput() {
	object.Env = viper.GetString(_env)
	object.Target = viper.GetString(_target)
	object.Type = viper.GetString(_type)
	object.Platform = viper.GetString(_platform)
	object.NodeType = viper.GetString(_nodeType)
	object.Version = viper.GetString(_version)
	object.Scheme = viper.GetString(_scheme)
	object.Set = viper.GetString(_cluster)
	object.PathIn = viper.GetString(_pathin)
	object.PathOut = viper.GetString(_pathout)
	object.TopicIpRange = viper.GetString(_topicIp)
	object.TopicPortRange = viper.GetString(_topicPort)
	object.TcpPortRange = viper.GetString(_tcpPort)
	object.EzeiEnv = viper.GetString(_ezeiEnv)
	object.EzeiInner = viper.GetString(_ezeiInner)
	object.EnvCover = viper.GetBool(_envCover)
	object.mode = viper.GetString(_mode)
}
