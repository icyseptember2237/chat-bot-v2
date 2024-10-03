package cmd

import (
	"chatbot/config"
	"chatbot/function"
	"chatbot/job"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
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

		manager := job.GetManager()
		for _, jobConf := range conf.Jobs {
			job := job.NewJob(jobConf)
			if job != nil {
				manager.Add(job)
			}
		}
		manager.StartAll()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		if functionServer != nil {
			functionServer.Stop()
		}

		job.GetManager().StopAll()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
