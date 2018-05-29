package remotes

import (
	"encoding/json"
)

var System = SystemRemote{}

const (
	Rpc_System_System           = "ServiceSystem"
	Rpc_System_Method_Ping    = Rpc_System_System+".Ping"
)

type SystemRemote struct {
	BaseRemote
}

// ping 检测是否联通
func (this *SystemRemote) Ping(ip string, port string, token string, args map[string]interface{}) (res map[string]string, err error) {
	replay, err := this.Call(ip, port, token, Rpc_System_Method_Ping, args, 300)
	if err != nil {
		return map[string]string{"version:":"null"}, err
	}
	if replay == "ok" {
		return map[string]string{"version:":"null"}, nil
	}
	json.Unmarshal([]byte(replay), &res)
	return res, nil
}