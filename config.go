// 生成全局配置和admin SSl证书
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

const configPath = "config/main.conf"

// 检测是否存在配置文件，不存在则新建
func init() {
	if _, err := os.ReadFile(configPath); err != nil {
		fmt.Println("未检测到配置文件，是否创建(y/n)")
		var is string
		fmt.Scanln(&is)
		if is != "y" {
			fmt.Println("退出成功")
			os.Exit(1)
		}
		writeConfig()
	}
}

type ConfigInfo struct {
	ServerIP  string
	Domain    string
	AdminPort string
	Port      string
}

func (c ConfigInfo) Get(name string) string {

	data, err := os.ReadFile(configPath)
	ErrprDisplay(err)

	err = json.Unmarshal(data, &c)
	ErrprDisplay(err)

	switch name {
	case "port":
		return c.Port

	case "serverIP":
		return c.ServerIP

	case "domain":
		return c.Domain

	case "adminPort":
		return c.AdminPort
	}
	return ""

}
func writeConfig() {
	configInfo := ConfigInfo{}
	os.Mkdir("config", 0700)
	os.Mkdir("config/cert/", 0700)
	os.Mkdir("config/cert/adminCa", 0700)

	//写入配置文件
	fmt.Println("输入你的服务器IP(默认:127.0.0.1)")
	fmt.Scanln(&configInfo.ServerIP)
	if configInfo.ServerIP == "" {
		configInfo.ServerIP = "127.0.0.1"
	}

	fmt.Println("输入你的域名(默认:localhost)")
	fmt.Scanln(&configInfo.Domain)
	if configInfo.Domain == "" {
		configInfo.Domain = "localhost"
	}

	fmt.Println("设置你的web端口(默认:8080)")
	fmt.Scanln(&configInfo.Port)
	if configInfo.Port == "" {
		configInfo.Port = ":8080"
	} else {
		configInfo.Port = ":" + configInfo.Port
	}

	fmt.Println("设置你的网站管理员端口(默认:59812)")
	fmt.Scanln(&configInfo.AdminPort)
	if configInfo.AdminPort == "" {
		configInfo.AdminPort = ":59812"
	} else {
		configInfo.AdminPort = ":" + configInfo.AdminPort
	}

	data, err := json.MarshalIndent(configInfo, "", "	")
	ErrprDisplay(err)
	err = os.WriteFile(configPath, data, 0600)

	createCA(configInfo.Domain)

	ErrprDisplay(err)
	exec.Command("bash", "-c", "chmod -R 400 config")

	fmt.Println("网站管理员端口设置为" + config.Get("adminPort"))
	fmt.Println("web端口为" + configInfo.Get("port"))
	fmt.Println("服务器域名是" + config.Get("domain"))
	fmt.Println("服务器IP是" + config.Get("serverIP"))
}

// 创建管理员https证书
func createCA(doMain string) {
	domain := config.Get("domain")
	shell("openssl genrsa -out ./config/cert/adminCa/ca.key 4096")

	shell(`openssl req -x509 -new -nodes -sha512 -days 3650 \
	-subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=` + domain + `" \
	-key ./config/cert/adminCa/ca.key \
	-out ./config/cert/adminCa/ca.crt`)

	shell(`openssl genrsa -out ./config/cert/adminCa/` + domain + `.key 4096`)

	shell(`openssl req -sha512 -new \
	-subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=` + domain + `" \
	-key ./config/cert/adminCa/` + domain + `.key \
	-out ./config/cert/adminCa/` + domain + `.csr`)

	shell(`cat > ./config/cert/adminCa/v3.ext <<-EOF
	authorityKeyIdentifier=keyid,issuer
	basicConstraints=CA:FALSE
	keyUsage=digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
	extendedKeyUsage=serverAuth
	subjectAltName=@alt_names
	
	[alt_names]
	DNS.1=` + domain + `
	EOF`)

	shell(`openssl x509 -req -sha512 -days 3650 \
	-extfile ./config/cert/adminCa/v3.ext \
	-CA ./config/cert/adminCa/ca.crt -CAkey ./config/cert/adminCa/ca.key -CAcreateserial \
	-in ./config/cert/adminCa/` + domain + `.csr \
	-out ./config/cert/adminCa/` + domain + `.crt`)

	// 删除不需要的证书
	shell("rm -rf ./config/cert/adminCa/*.csr")
	shell("rm -rf ./config/cert/adminCa/v3.ext")
	shell("rm -rf ./config/cert/adminCa/ca.crt")
	shell("rm -rf ./config/cert/adminCa/ca.key")
}

// 执行命令
func shell(cmd string) {

	tmp := exec.Command("bash", "-c", cmd)
	tmp.Stdout = os.Stdout
	tmp.Stderr = os.Stderr
	tmp.Run()
}
