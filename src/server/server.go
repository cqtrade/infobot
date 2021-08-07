package server

import (
	"github.com/cqtrade/infobot/src/config"
	tvcontroller "github.com/cqtrade/infobot/src/controller"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    config.Config
	server *gin.Engine
	tvCtrl tvcontroller.TvController
}

func New(cfg config.Config, tvCtrl tvcontroller.TvController) *Server {
	server := gin.New()
	return &Server{
		cfg:    cfg,
		server: server,
		tvCtrl: tvCtrl,
	}
}

func (s *Server) Run() {
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

		apiRoutes.POST("/text", s.tvCtrl.PostText)
		apiRoutes.POST("/flash", s.tvCtrl.PostFlash)
	}

	s.server.Run(s.cfg.ServerUrl)
}
