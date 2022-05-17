package cmd

import (
	"fmt"

	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

type LinuxCmd struct {
	Client *ssh.Client
}

func (c *LinuxCmd) Mem(command string) float64 {

	session, _ := c.Client.NewSession()

	// 执行监控项传递过来的数据
	cmd, err := session.CombinedOutput(command)
	if err != nil {
		fmt.Print("远程执行cmd失败", err)
	}

	// 因为数据转换过来有 \n 的换行符，所以这里通过 strings.Split 函数来分割 \n 换行符
	Memory := strings.Split(strings.TrimSpace(string(cmd)), `\n`)

	// 由于得到的 Memory 是一个切片所以这里取他的第一个元素
	newMemory := Memory[0]

	// 转换为 float64 方便 Prometheus 监控获取的值
	availableMemory, _ := strconv.ParseFloat(newMemory, 64)

	// 将 KB 转为 GB
	newavailableMemory := availableMemory / 1048576

	// 保留3位小数
	newavailableMemory, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", newavailableMemory), 64)

	// 返回监控指标项
	return newavailableMemory
}
