package handler

import (
	"encoding/base64"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// 定义一个切片也存放传入的用户名密码信息
type AuthSecrets map[string]string

// auth 函数传参 http.Handler，因为 promhttp.Handler() 是一个 Handler 类型
// 并且返回值也是一个 Handler，因为在 main 包中 http.Handle 接收 Handler 类型
// 传入 secrets AuthSecrets 实现认证
func Auth(handler http.Handler, secrets AuthSecrets) http.Handler {

	// 这里返回 HandlerFunc ，因为 HandlerFunc 也是 http.Handler 接口的一个实现
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		// 验证，现在当用户要访问 promhttp 就得在这里进行身份验证
		// 获取请求体中 authorization
		// 拿到 authorization 之后需要进行验证
		secret := request.Header.Get("Authorization")

		// 如果 isAuth 没有通过为 false 那就在响应体中返回状态码 401
		if !isAuth(secret, secrets) {
			// 实现弹出登录框
			response.Header().Set("WWW-Authenticate", `Basic realm=""`)

			// 返回响应码 401
			response.WriteHeader(401)
			// 最后退出程序
			return
		}

		// 如果通过调用原生 handler
		// func (http.Handler).ServeHTTP(http.ResponseWriter, *http.Request)
		handler.ServeHTTP(response, request)
	})
}

func isAuth(secret string, secrets AuthSecrets) bool {
	// 这就是不适用密码验证的方式直接登录
	if secrets == nil {
		return true
	}

	// 由于输入的是 Basic base64 编码的字符串格式如：Basic name:password ，所以需要先去除掉 Basic 字符串
	// strings.Fields ，通过空格分割返回一个 []string
	node := strings.Fields(secret)

	// 判断 node 长度是否不等于 2 ，因为 strings.Fields 会将 Basic name:password 切割为长度为 2 的切片
	// 如果不等于 2 就说明用户没有提交账户密码就 return false
	if len(node) != 2 {
		return false
	}

	// 现在 node[1] 就是一个用户 base64 的编码字符串，所以需要通过解码
	// 通过 base64 解码, DecodeString 返回的是一个 []byte
	plaintext, err := base64.StdEncoding.DecodeString(node[1])
	if err != nil {
		return false
	}

	// 这里通过 SplitN 来取 : 分割之后的第二个元素就是用户在 web 页面的登录数据
	node = strings.SplitN(string(plaintext), ":", 2)
	// 这里通过判断 node 长度是否不等于 2
	if len(node) != 2 {
		return false
	}

	// secrets 这个 map 中将获取 node[0] 的 key 并拿到 node[0] 的 value 赋值给 password
	password, ok := secrets[node[0]]

	// 如果能拿到则为 ok 并且 password 和 node[1] 匹配成功就 return true
	// bcrypt 进行对 node[1] 用户提交的密码和 password 在程序中定义的密码进行验证，如果 err == nil 就为真
	return ok && bcrypt.CompareHashAndPassword([]byte(password), []byte(node[1])) == nil

}
