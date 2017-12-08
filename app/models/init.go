package models

import (
	"fmt"
	"os"
	"github.com/astaxie/beego"
	"github.com/snail007/go-activerecord/mysql"
)

var G *mysql.DBGroup

func init() {
	//init config file
	env := os.Getenv("SQENV")
	if env == "" {
		env = "dev"
	}
	if env == "dev" {
		beego.LoadAppConfig("ini", "conf/dev.conf")
	}
	if env == "test" {
		beego.LoadAppConfig("ini", "conf/test.conf")
	}
	if env == "prod" {
		beego.LoadAppConfig("ini", "conf/prod.conf")
	}
	var err error

	//init db
	host := beego.AppConfig.String("db::host")
	port, _ := beego.AppConfig.Int("db::port")
	user := beego.AppConfig.String("db::user")
	pass := beego.AppConfig.String("db::pass")
	dbname := beego.AppConfig.String("db::name")
	dbTablePrefix := beego.AppConfig.String("db::table_prefix")
	maxIdle, _ := beego.AppConfig.Int("db::conn_max_idle")
	maxConn, _ := beego.AppConfig.Int("db::conn_max_connection")
	cfg := mysql.NewDBConfigWith(host, port, dbname, user, pass)
	cfg.SetMaxIdleConns = maxIdle
	cfg.SetMaxOpenConns = maxConn
	cfg.TablePrefix = dbTablePrefix
	cfg.TablePrefixSqlIdentifier = "__PREFIX__"
	err = G.Regist("default", cfg)
	if err != nil {
		beego.Error(fmt.Errorf("regist db error:%s,with config : %v", err, cfg))
		os.Exit(100)
	}
}
