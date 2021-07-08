package app

import (
	"github.com/cqtrade/infobot/src/config"
	tvcontroller "github.com/cqtrade/infobot/src/controller"
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/server"
)

var (
	cfg          config.Config             = *config.New()
	ft           ftxtrade.FtxTrade         = *ftxtrade.New(cfg)
	notif        notification.Notification = *notification.New(cfg)
	tvController tvcontroller.TvController = *tvcontroller.New(notif, ft)
	webServer    server.Server             = *server.New(cfg, tvController)
)

func Run() {
	webServer.Run()
}
