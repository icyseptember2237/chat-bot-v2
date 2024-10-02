package gorm

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DBTypeMysql      string = "mysql"
	DBTypePostgresql string = "postgresql"
)

var (
	cache = map[string]*gorm.DB{}
)

// Init init gorm DB connection.
func Init(ctx context.Context, configs map[string]string, dbtype string) error {
	for name, uri := range configs {
		switch dbtype {
		case DBTypeMysql:
			client, err := gorm.Open(mysql.Open(uri), &gorm.Config{})
			if err != nil {
				return err
			}
			sqlDB, err := client.DB()
			if err != nil {
				fmt.Println("connect db server connection pool failed.")
				return err
			}
			sqlDB.SetConnMaxIdleTime(10)
			sqlDB.SetMaxOpenConns(100)

			cache["mysql_"+name] = client
		case DBTypePostgresql:
			client, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
			if err != nil {
				return err
			}
			cache["postgres_"+name] = client
		}
	}
	return nil
}

// Get get connection.
func Get(name string) *gorm.DB {
	return cache[name]
}
