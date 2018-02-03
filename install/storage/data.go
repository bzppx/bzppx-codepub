package storage

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
)

var Data = NewData()

var installChan = make(chan int, 1)

const License_Disagree = 0 // 协议不同意
const License_Agree = 1 // 协议同意

const Env_NotAccess = 0 // 环境检测不通过
const Env_Access = 1 // 环境检测通过

const Sys_NotAccess = 0 // 系统配置不通过
const Sys_Access = 1 // 系统配置通过

const Database_NotAccess = 0 // 数据库配置不通过
const Database_Access = 1 // 数据库配置通过

const Install_Ready = 0 // 安装准备阶段
const Install_Start = 1 // 安装开始
const Install_End = 2 // 安装完成

const Install_Default = 0 // 默认
const Install_Failed = 1 // 安装失败
const Install_Success = 2 // 安装成功

var defaultSystemConf = map[string]interface{}{
	"addr": "0.0.0.0",
	"port": "8080",
}

var defaultDatabaseConf = map[string]interface{}{
	"host": "127.0.0.1",
	"port": "3306",
	"name": "codepub",
	"user": "",
	"pass": "",
	"table_prefix": "cp_",
	"conn_max_idle": 30,
	"conn_max_connection": 200,
	"admin_name": "",
	"admin_pass": "",
}

func NewData() data {
	return data{
		License: License_Disagree,
		Env: Env_NotAccess,
		System: Sys_NotAccess,
		Database: Database_NotAccess,
		SystemConf: defaultSystemConf,
		DatabaseConf: defaultDatabaseConf,
		Status: Install_Ready,
		Result: "",
		IsSuccess: Install_Default,
	}
}

type data struct {
	License int
	Env int
	System int
	Database int
	SystemConf map[string]interface{}
	DatabaseConf map[string]interface{}
	Status int
	Result string
	IsSuccess int
}

// check db
func checkDB() (err error) {

	host := Data.DatabaseConf["host"].(string)
	port := Data.DatabaseConf["port"].(string)
	user := Data.DatabaseConf["user"].(string)
	pass := Data.DatabaseConf["pass"].(string)

	db, err := sql.Open("mysql", user+":"+pass+"@tcp("+host+":"+port+")/")
	if err != nil {
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return
	}
	return
}

// create db
func createDB() (err error) {

	host := Data.DatabaseConf["host"].(string)
	port := Data.DatabaseConf["port"].(string)
	user := Data.DatabaseConf["user"].(string)
	pass := Data.DatabaseConf["pass"].(string)
	name := Data.DatabaseConf["name"].(string)

	db, err := sql.Open("mysql", user+":"+pass+"@tcp("+host+":"+port+")/")
	if err != nil {
		return
	}
	defer db.Close()
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS "+name+" CHARACTER SET utf8")
	if err != nil {
		return
	}
	return nil
}

// create table
func createTable() (err error) {

	host := Data.DatabaseConf["host"].(string)
	port := Data.DatabaseConf["port"].(string)
	user := Data.DatabaseConf["user"].(string)
	pass := Data.DatabaseConf["pass"].(string)
	name := Data.DatabaseConf["name"].(string)

	installDir, _ := os.Getwd()
	installDir = strings.Replace(installDir, "install", "", 1)
	sqlBytes, err := ioutil.ReadFile(installDir+"docs/databases/table.sql");
	if err != nil {
		return err
	}
	sqlTable := string(sqlBytes);
	fmt.Println(sqlTable)
	db, err := sql.Open("mysql", user+":"+pass+"@tcp("+host+":"+port+")/"+name)
	if err != nil {
		return
	}
	defer db.Close()
	_, err = db.Exec(sqlTable)
	if err != nil {
		return
	}
	return nil
}

func installFailed(err string)  {
	Data.Result = err
	Data.Status = Install_End
	Data.IsSuccess = Install_Failed
	log.Println(err)
}

func installSuccess()  {
	Data.Status = Install_End
	Data.IsSuccess = Install_Success
}

func StartInstall()  {
	installChan <- 1
}

func ListenInstall()  {

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("%v", err)
			}
		}()
		for  {
			select {
			case <-installChan:
				Data.Status = Install_Start
				// 开始安装
				log.Println("codepub start install")
				// 检查db
				err := checkDB()
				if err != nil {
					installFailed("连接数据库出错："+err.Error())
					continue
				}
				// 创建数据库
				err = createDB()
				if err != nil {
					installFailed("创建数据库出错："+err.Error())
					continue
				}
				// 创建表
				err = createTable()
				if err != nil {
					installFailed("创建表出错："+err.Error())
					continue
				}

				installSuccess()
				return
			}
		}
	}()
}

func init()  {
	ListenInstall()
}