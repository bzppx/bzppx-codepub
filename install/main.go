package main

import (
	_ "bzppx-codepub/install/storage"
	"github.com/astaxie/beego"
	"flag"
)

// 安装程序

var (
	port = flag.String("port", "8090", "listen port")
)

func main() {
	flag.Parse()
	beego.Run(":"+*port)
}

