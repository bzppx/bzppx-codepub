package remotes

import (
	"crypto/tls"
	"errors"
	"net/rpc"
)

type BaseRemote struct {

}

func (b *BaseRemote) Call(ip string, port string, method string, args map[string]interface{}) (reply interface{}, err error) {
	address := ip + ":" +port
	if address == "" {
		return nil, errors.New("codepub connect agent error: ip:port is not empty!")
	}
	if method == "" {
		return nil, errors.New("codepub connect agent error: method is not empty!")
	}

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", address, conf)
	if err != nil {
		return nil, errors.New("codepub connect agent error: " + err.Error())
	}
	conn.Write([]byte("agent-code"))

	var buf = make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		return nil, errors.New("codepub read conn error: "+err.Error())
	}
	if string(buf[:n]) != "success" {
		return nil, errors.New("codepub connect agent token error!")
	}
	client := rpc.NewClient(conn)

	err = client.Call(method, args, &reply)
	if err != nil {
		return nil, errors.New("codepub call agent error: "+err.Error())
	}

	return reply, nil
}