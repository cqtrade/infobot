package server

import (
	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/controller"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Run()
}

type serv struct {
	cfg    config.Config
	server *gin.Engine
	tvCtrl controller.TradingviewController
}

func New(cfg config.Config, tvCtrl controller.TradingviewController) Server {
	server := gin.New()
	return &serv{
		cfg:    cfg,
		server: server,
		tvCtrl: tvCtrl,
	}
}

func (s *serv) Run() {
	s.server.Use(gin.Recovery(), gin.Logger())

	s.server.GET("/", func(c *gin.Context) {
		c.Abort()
	})

	s.server.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ok": true,
		})
	})

	apiRoutes := s.server.Group("/api")
	{
		apiRoutes.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"ok": true})
		})

		apiRoutes.POST("/signal-text", s.tvCtrl.PostText)
		apiRoutes.POST("/signal-json", s.tvCtrl.PostJson)
		apiRoutes.POST("/flash", s.tvCtrl.PostFlash)
	}

	s.server.Run(s.cfg.GetServerUrl())
}
