package main

import (
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
)

var config ConfigInfo

func main() {
	// 生产模式
	gin.SetMode(gin.ReleaseMode)

	// 提供公开服务的服务器
	go Server()

	// 管理员服务器
	adminServer := &http.Server{
		Addr:    config.Get("adminPort"),
		Handler: AdminServer(),
		TLSConfig: &tls.Config{

			// 必需有证书才能访问网站
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	// 管理员服务器
	// 这里会阻塞后续代码执行
	adminServer.ListenAndServeTLS("./config/cert/adminCa/"+config.Get("domain")+".crt", "./config/cert/adminCa/"+config.Get("domain")+".key")

}
