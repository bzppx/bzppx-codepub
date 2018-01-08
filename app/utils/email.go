package utils

import (
	"crypto/tls"
	"net"
	"net/smtp"
	"strings"
)

func NewEmail() *Email {
	return &Email{}
}

type Email struct {
}

func (email *Email) SendEmail(emailPara map[string]string, content, contentType, subject string) (err error) {
	var from string
	auth := smtp.PlainAuth("", emailPara["email_username"], emailPara["email_password"], emailPara["email_host"])
	if emailPara["email_from"] == "" {
		from = emailPara["email_username"]
	} else {
		from = emailPara["email_from"]
	}
	if contentType != "html" {
		contentType = "plain"
	}
	sendTo := strings.Split(emailPara["email_cc_list"], "\n")
	ccList := strings.Join(sendTo, ";")
	msg := []byte("To: " + ccList + "\r\nFrom: " + from + "\r\nSubject: " + subject + "\r\nContent-Type: text/" + contentType + "; charset=UTF-8 \r\n\r\n" + content)
	if emailPara["email_is_ssl"] == "1" {
		err = smtp.SendMail(emailPara["email_host"]+":"+emailPara["email_port"], auth, emailPara["email_username"], sendTo, msg)
	} else {
		err = SendMailUsingTLS(emailPara["email_host"]+":"+emailPara["email_port"], auth, emailPara["email_username"], sendTo, msg)
	}

	return
}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	c, err := Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
