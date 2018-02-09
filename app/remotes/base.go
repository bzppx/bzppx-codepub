package remotes

import (
	"crypto/tls"
	"errors"
	"net/rpc"
	"time"
	"net"
)

const Conn_Timeout = 300

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
		return reply, errors.New("codepub connect agent error: " + err.Error())
	}
	conn.Write([]byte(token))

	var buf = make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		return reply, errors.New("codepub read conn error: "+err.Error())
	}
	if string(buf[:n]) != "success" {
		return reply, errors.New("codepub connect agent token error!")
	}
	rpcClient := rpc.NewClient(conn)

	err = rpcClient.Call(method, args, &reply)
	if err != nil {
		return reply, errors.New("codepub call agent error: "+err.Error())
	}

	return reply, nil
}