package tvcontroller

import (
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/types"
	"github.com/gin-gonic/gin"
)

type TvController struct {
	notification notification.Notification
	ftxTrade     ftxtrade.FtxTrade
}

func New(notification notification.Notification, ftxTrade ftxtrade.FtxTrade) *TvController {
	return &TvController{
		notification: notification,
		ftxTrade:     ftxTrade,
	}
}

func (tvc *TvController) PostText(ctx *gin.Context) {
	ctx.JSON(200, gin.H{})
	text, err := ctx.GetRawData()
	if err != nil {
		println("error " + err.Error()) // TODO needs logger
		return
	}
	go tvc.notification.SendTextMessage(string(text))
}

func (tvc *TvController) PostFlash(ctx *gin.Context) {
	ctx.JSON(200, gin.H{})
	var message types.JSONMessageBody
	err := ctx.ShouldBindJSON(&message)
	if err != nil {
		println("error " + err.Error()) // TODO needs logger
		return
	}
	if message.Signal == 1001 {
		go tvc.ftxTrade.BuyEthBull("test1")
	} else if message.Signal == -2002 {
		go tvc.ftxTrade.TpEthBull("test1")
	} else if message.Signal == 888 {
		go tvc.ftxTrade.ArbStart("arb1", "BTC")
	} else if message.Signal == -888 {
		go tvc.ftxTrade.ArbEnd("arb1", "BTC")
	}
	/**
	1
	-1
	2
	-2
		**/

	go tvc.notification.SendFlashMessage(message)
}
