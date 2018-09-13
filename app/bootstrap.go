package app

import (
	"os"
	"log"
	"fmt"
	"flag"
	"net/http"
	"bzppx-codepub/app/utils"
	"bzppx-codepub/app/models"
	"github.com/astaxie/beego"
	"bzppx-codepub/app/container"
	"bzppx-codepub/app/controllers"
	"github.com/snail007/go-activerecord/mysql"
)

var (
	confPath = flag.String("conf", "conf/default.conf", "please set codepub conf path")
)

var (
	Version = "v0.8.4"
)

func init()  {
	poster()
	initConfig()
	initDB()
	initRouter()
	container.InitWorker()
	controllers.Version = Version
}

// poster logo
func poster() {
	logo := `
                    __                     __
  _____ ____   ____/ /___   ____   __  __ / /_
 / ___// __ \ / __  // _ \ / __ \ / / / // __ \
/ /__ / /_/ // /_/ //  __// /_/ // /_/ // /_/ /
\___/ \____/ \__,_/ \___ / .___/ \__,_//_.___/
                        /_/
`+
"Author: bzppx\r\n"+
"Version: "+Version+"\r\n"+
"Link: https://github.com/bzppx/bzppx-codepub"
fmt.Println(logo)
}

// init beego config
func initConfig() {

	flag.Parse()

	if *confPath == "" {
		log.Println("conf file not empty!")
		os.Exit(1)
	}
	ok, _ := utils.NewFile().PathIsExists(*confPath)
	if ok == false {
		log.Println("conf file " + *confPath + " not exists!")
		os.Exit(1)
	}
	//init config file
	beego.LoadAppConfig("ini", *confPath)

	// init name
	beego.AppConfig.Set("sys.name", "codepub")
	beego.BConfig.AppName = beego.AppConfig.String("sys.name")
	beego.BConfig.ServerName = beego.AppConfig.String("sys.name")

	// set static path
	beego.SetStaticPath("/static/", "static")

	// views path
	beego.BConfig.WebConfig.ViewsPath = "views/"

	// session
	beego.BConfig.WebConfig.Session.SessionName = "ssid"
	beego.BConfig.WebConfig.Session.SessionOn = true

	// log
	logConfigs, err := beego.AppConfig.GetSection("log")
	if err != nil {
		log.Println(err.Error())
	}
	for adapter, config := range logConfigs {
		beego.SetLogger(adapter, config)
	}
	beego.SetLogFuncCall(true)
}

func initRouter() {
	// router
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.RouterCaseSensitive = false
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime, _ = beego.AppConfig.Int("sessiongcmaxlifetime")

	// todo add router..
	beego.AutoRouter(&controllers.MainController{})
	beego.AutoRouter(&controllers.UserController{})
	beego.AutoRouter(&controllers.LoginController{})
	beego.AutoRouter(&controllers.ProfileController{})
	beego.AutoRouter(&controllers.PublishController{})
	beego.AutoRouter(&controllers.GroupController{})
	beego.AutoRouter(&controllers.ProjectController{})
	beego.AutoRouter(&controllers.LogController{})
	beego.AutoRouter(&controllers.NodesController{})
	beego.AutoRouter(&controllers.NodeController{})
	beego.AutoRouter(&controllers.ConfigureController{})
	beego.AutoRouter(&controllers.TaskController{})
	beego.AutoRouter(&controllers.TaskLogController{})
	beego.AutoRouter(&controllers.NoticeController{})
	beego.AutoRouter(&controllers.StatisticsController{})
	beego.AutoRouter(&controllers.ApiAuthController{})
	beego.AutoRouter(&controllers.ContactController{})
	beego.Router("/", &controllers.LoginController{}, "*:Index")
	beego.ErrorHandler("404", http_404)
	beego.ErrorHandler("500", http_500)

	// add template func
	beego.AddFuncMap("dateFormat", utils.NewDate().Format)
}

//init db
func initDB() {
	host := beego.AppConfig.String("db::host")
	port, _ := beego.AppConfig.Int("db::port")
	user := beego.AppConfig.String("db::user")
	pass := beego.AppConfig.String("db::pass")
	dbname := beego.AppConfig.String("db::name")
	dbTablePrefix := beego.AppConfig.String("db::table_prefix")
	maxIdle, _ := beego.AppConfig.Int("db::conn_max_idle")
	maxConn, _ := beego.AppConfig.Int("db::conn_max_connection")
	models.G = mysql.NewDBGroup("default")
	cfg := mysql.NewDBConfigWith(host, port, dbname, user, pass)
	cfg.SetMaxIdleConns = maxIdle
	cfg.SetMaxOpenConns = maxConn
	cfg.TablePrefix = dbTablePrefix
	cfg.TablePrefixSqlIdentifier = "__PREFIX__"
	err := models.G.Regist("default", cfg)
	if err != nil {
		beego.Error(fmt.Errorf("database error:%s,with config : %v", err, cfg))
		os.Exit(100)
	}
}

func http_404(rs http.ResponseWriter, req *http.Request) {
	rs.Write([]byte("404 not found!"))
}

func http_500(rs http.ResponseWriter, req *http.Request) {
	rs.Write([]byte("500 server error!"))
}
