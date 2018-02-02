package storage

var Data = NewData()

const License_Disagree = 0 // 协议不同意
const License_Agree = 1 // 协议同意

const Env_NotAccess = 0 // 环境检测不通过
const Env_Access = 1 // 环境检测通过

const Install_Ready = 0 // 安装准备阶段
const Install_Start = 1 // 安装开始
const Install_End = 2 // 安装完成

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
		License: License_Agree,
		Env: Env_Access,
		SystemConf: defaultSystemConf,
		DatabaseConf: defaultDatabaseConf,
		Status: Install_Ready,
	}
}

type data struct {
	License int
	Env int
	SystemConf map[string]interface{}
	DatabaseConf map[string]interface{}
	Status int
}