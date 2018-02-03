package main

import (
	_ "bzppx-codepub/install/storage"
	"github.com/astaxie/beego"
)

// 安装程序

func main() {
	beego.Run(":8089")
}

