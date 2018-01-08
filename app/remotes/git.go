package remotes

import "fmt"

var Git = GitRemote{}

const (
	Rpc_Git_Service = "ServiceGit"
	Rpc_Git_Method_GetCommitId = Rpc_Git_Service+".GetCommitId"
	Rpc_Git_Method_Publish     = Rpc_Git_Service+".Publish"
)

type GitRemote struct {
	BaseRemote
}

// 获取 commit id
func (this *GitRemote) GetCommitId(ip string, port string, args map[string]interface{}) (string) {
	reply, err := this.Call(ip, port, Rpc_Git_Method_GetCommitId, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(reply)
	return ""
}

// 发布
func (this *GitRemote) Publish(ip string, port string, args map[string]interface{}) {
	reply, err := this.Call(ip, port, Rpc_Git_Method_Publish, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(reply)
}

// 回滚
func (this *GitRemote) Rollback() {

}

// 获取节点执行结果
func (this *GitRemote) GetResults(ip string, port string, args map[string]interface{}) {

}