package remotes

var System = SystemRemote{}

const (
	Rpc_System_System           = "ServiceSystem"
	Rpc_System_Method_Ping    = Rpc_System_System+".Ping"
)

type SystemRemote struct {
	BaseRemote
}

// ping 检测是否联通
func (this *SystemRemote) Ping(ip string, port string, args map[string]interface{}) error {
	_, err := this.Call(ip, port, Rpc_System_Method_Ping, args)
	return err
}