package remotes

import (
	"crypto/tls"
	"errors"
	"net/rpc"
	"time"
	"net"
	"bzppx-codepub/app/utils"
	"github.com/astaxie/beego"
)

var Conn_Timeout = beego.AppConfig.DefaultInt64("agent.tls_timeout", 1000)

type BaseRemote struct {

}

func (b *BaseRemote) Call(ip string, port string, token string, method string, args map[string]interface{}, timeout int64) (reply string, err error) {
	address := ip + ":" +port
	if address == "" {
		return reply, errors.New("codepub connect agent error: ip:port is not empty!")
	}
	if method == "" {
		return reply, errors.New("codepub connect agent error: method is not empty!")
	}
	if token == "" {
		return reply, errors.New("codepub connect agent error: token is not empty!")
	}
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	dialer := &net.Dialer{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", address, conf)
	if err != nil {
		return reply, errors.New("codepub connect agent tls handle error: " + err.Error())
	}

	defer conn.Close()

	// encode pack write
	tokenEncode, err := utils.NewCodec().EncodePack([]byte(token))
	if err != nil {
		conn.Close()
		return reply, errors.New("codepub encode pack error, " + err.Error())
	}
	_, err = conn.Write([]byte(tokenEncode))
	if err != nil {
		return reply, errors.New("codepub send token error, " + err.Error())
	}

	// read decode pack
	str, err := utils.NewCodec().DecodePack(conn)
	if err != nil {
		return reply, errors.New("codepub read token result error: "+err.Error())
	}
	if str != "success" {
		return reply, errors.New("codepub connect agent token error!")
	}
	rpcClient := rpc.NewClient(conn)

	err = rpcClient.Call(method, args, &reply)
	if err != nil {
		return reply, errors.New("codepub call agent error: "+err.Error())
	}

	return reply, nil
}