package cmd

import (
	"chatbot/config"
	"chatbot/function"
	"chatbot/function/hook"
	"chatbot/job"
	"chatbot/worker"
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
			if conf.Server.OnlyWhiteGroup {
				functionServer.AddBeforeHook(hook.OnlyWhiteList)
				fmt.Println("enable only white list")
			}
			if conf.Server.SaveMessage {
				functionServer.AddBeforeHook(hook.SaveMessage)
				fmt.Println("enable save message")
			}
			if conf.Server.SaveImage {
				functionServer.AddBeforeHook(hook.GetImage)
				fmt.Println("enable save image")
			}
			go functionServer.Start()
		}

		jobManager := job.GetManager()
		for _, jobConf := range conf.Jobs {
			job := job.NewJob(jobConf)
			if job != nil {
				jobManager.Add(job)
			}
		}
		jobManager.StartAll()

		workerManager := worker.GetManager()
		for _, workerConf := range conf.Workers {
			worker := worker.New(workerConf)
			if worker != nil {
				workerManager.Add(worker)
			}
		}
		workerManager.StartAll()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		if functionServer != nil {
			functionServer.Stop()
		}

		jobManager.StopAll()
		workerManager.StopAll()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
