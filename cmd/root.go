package cmd

import (
	"chatbot/config"
	"chatbot/function"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.Flags().StringP("config", "c", "./config/config.json", "配置文件路径, 默认./config/config.json")
}

var conf *config.Config
var rootCmd = &cobra.Command{
	Use:    "QQBot",
	Short:  "Lua脚本QQBot",
	PreRun: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		functionServer := function.New(conf.Server)
		if functionServer != nil {
			go functionServer.Start()
		}
		quit := make(chan os.Signal, 1)
		<-quit
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
