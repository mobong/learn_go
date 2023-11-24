package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"learn_go/conf"
	"strconv"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func dbStartUp() (*xorm.Engine, error) {
	conf.ConfInit()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		strconv.Itoa(viper.GetInt("mysql.port")),
		viper.GetString("mysql.dbname"),
		viper.GetString("mysql.charset"))

	db, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))

	db.SetLogLevel(log.LogLevel(viper.GetInt("mysql.log_level")))
	db.ShowSQL(viper.GetBool("mysql.show_sql"))

	return db, nil
}
