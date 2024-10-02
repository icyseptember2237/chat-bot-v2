package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gotemplate/server"
	"log"
)

func init() {
	rootCmd.Flags().Int32P("port", "p", 8080, "运行端口, 默认8080")
	rootCmd.Flags().StringP("config", "c", "./config/config.json", "配置文件路径, 默认./config/config.json")
}

var rootCmd = &cobra.Command{
	Use:    "模版",
	Short:  "Go基于gin, viper, cobra....封装的模版",
	PreRun: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		go func() {
			httpServer := server.New("http", ":8080").(*server.HttpServer)
			if err := httpServer.Server.Run(ctx); err != nil {
				log.Fatalln(err)
			}
		}()
		select {}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
