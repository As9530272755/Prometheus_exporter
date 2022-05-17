package sshlink

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

var (
	Host string
)

func Link() *ssh.Client {
	Host = "10.0.0.30"
	port := "22"
	user := "root"
	password := "666666"

	// 创建ssh登录配置
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// 登录类型 password
	config.Auth = []ssh.AuthMethod{ssh.Password(password)}

	// 定义链接地址
	addr := fmt.Sprintf("%s:%s", Host, port)

	// 登录远程主机
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println(err)
	}

	return client
}
