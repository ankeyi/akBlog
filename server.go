package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 公开的web服务
func Server() {
	// 使用默认中间件(Logger和Recovery)
	r := gin.Default()

	// 主页头像
	r.StaticFile("/header.svg", "web/header.svg")

	//引入网站所需静态资源
	r.Static("/assets", "web/assets")

	// 显示mirrors目录
	r.StaticFS("/static/files/", http.Dir("web/files"))

	// 显示静态页面
	WebPage(r)

	CreateHash()

	r.Run(":8080")
}

// 管理员后台服务器
func AdminServer() http.Handler {
	// 使用默认路由
	adminR := gin.New()

	// 导入mirrors文件
	adminR.LoadHTMLFiles("web/mirrors.html")

	// index页面
	adminR.GET("/", func(ctx *gin.Context) {

		// 验证文件哈希值
		success, fail := CheckingHash()

		// 返回网页验证记录
		ctx.HTML(200, "mirrors.html", gin.H{"success": success, "fail": fail})
	})

	// 文件上传
	adminR.POST("/saveFile", func(ctx *gin.Context) {
		// 单独文件上传
		file, _ := ctx.FormFile("file")

		//保存文件
		ctx.SaveUploadedFile(file, "web/files/"+file.Filename)

		// 写入HASH
		addHash(file.Filename)
	})

	return adminR
}
