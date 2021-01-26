package main

import (
	"Learnos/Web/action"
	"Learnos/Web/action/middleware"
	"Learnos/Web/sshServer"
	"Learnos/common/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func init() {
	err := config.ReadConf("config.toml")
	if err != nil {
		panic(err)
	}
}

func main() {
	conf := config.GetConf()
	//log.SetFlags(log.Lshortfile)
	gin.DisableConsoleColor()
	gin.SetMode(conf.Web.RunMode)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("./Web/view/*")
	r.Static("/static", "./Web/static")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/login")
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.GET("/console", func(c *gin.Context) {
		c.HTML(http.StatusOK, "console.html", nil)
	})
	r.GET("/term/:cid", action.Term)
	api := r.Group("/api/v1")
	{
		api.POST("/register", action.Register)
		api.POST("/login", action.Login)
		api.POST("/send", action.Send)
		api.Use(middleware.Auth()).GET("/getImage", action.GetImages)
		api.Use(middleware.Auth()).GET("/refreshToken", action.UpdateToken)
		api.Use(middleware.Auth()).POST("/createContainer", action.CreateContainer)
		api.Use(middleware.Auth()).POST("/deleteContainer", action.Delete)
		api.Use(middleware.Auth()).POST("/getStatus", action.GetStatus)
	}
	go func() {
		if err := sshServer.Run(conf.Web.SSHAddr); err != nil {
			log.Println("ssh服务启动失败：", err.Error())
			return
		}
	}()
	if err := http.ListenAndServe(conf.Web.WebAddr, r); err != nil {
		log.Fatal("Web服务启动失败：", err.Error())
	}
}
