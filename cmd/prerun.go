package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gotemplate/config"
	"gotemplate/logger"
	"gotemplate/storage/gorm"
	"gotemplate/utils/exit"
	"time"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/config/")
	viper.AddConfigPath("./config")
}

func preRun(cmd *cobra.Command, args []string) {
	configFile, _ := cmd.Flags().GetString("config")
	logger.Infof(context.Background(), "cmdline config file: %v", configFile)
	if _, err := config.Load(configFile); err != nil {
		panic(fmt.Errorf("config load failed: %s", err.Error()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conf := config.Get()
	// TODO Mongo Redis初始化
	if conf.Resources != nil && len(conf.Resources.Storage.Mysql) > 0 {
		if err := gorm.Init(ctx, conf.Resources.Storage.Mysql, gorm.DBTypeMysql); err != nil {
			panic(fmt.Errorf("MySQL init error: %s", err.Error()))
		}
	}
	if conf.Resources != nil && len(conf.Resources.Storage.Postgresql) > 0 {
		if err := gorm.Init(ctx, conf.Resources.Storage.Postgresql, gorm.DBTypePostgresql); err != nil {
			panic(fmt.Errorf("PostgreSQL init error: %s", err.Error()))
		}
	}
	go exit.HouseKeeping()
}
