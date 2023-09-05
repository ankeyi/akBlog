package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

// 命令快捷执行
func Command(shell string) {
	cmd := exec.Command("/bin/bash", "-c", shell)
	cmd.Run()
}

// 错误显示函数
func ErrprDisplay(err interface{}) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 前端页面生成
func WebPage(r *gin.Engine) {
	// 加载mirrors文件
	r.LoadHTMLFiles("web/mirrors.html")

	r.GET("/", func(ctx *gin.Context) {
		ctx.File("./web/index.html")
	})

	r.GET("/categories", func(ctx *gin.Context) {
		ctx.File("./web/categories/index.html")
	})

	r.GET("/tags", func(ctx *gin.Context) {
		ctx.File("./web/tags/index.html")
	})

	r.GET("/archives", func(ctx *gin.Context) {
		ctx.File("./web/archives/index.html")
	})

	r.GET("/about", func(ctx *gin.Context) {
		ctx.File("./web/about/index.html")
	})

	r.GET("/mirrors", func(ctx *gin.Context) {
		success, fail := CheckingHash()
		ctx.HTML(200, "mirrors.html", gin.H{"sucess": success, "fail": fail})
	})

	// 文章页面路径生成
	func() {
		article := r.Group("/posts/")
		file, _ := os.ReadDir("./web/posts/")
		var page []string
		for _, v := range file {
			if v.IsDir() {
				page = append(page, v.Name())
			}
		}

		for _, v := range page {
			article.GET(v, func(ctx *gin.Context) {
				ctx.File("./web/posts/" + v + "/index.html")
			})
		}
	}()
}
