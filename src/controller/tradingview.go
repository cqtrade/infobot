package controller

import (
	"github.com/cqtrade/infobot/src/service"
	"github.com/cqtrade/infobot/src/types"
	"github.com/gin-gonic/gin"
)

type TradingviewController interface {
	PostText(ctx *gin.Context)
	PostJson(ctx *gin.Context)
}

type tvController struct {
	discordService service.DiscordService
}

func NewTradingviewController(discordService service.DiscordService) TradingviewController {
	return &tvController{
		discordService: discordService,
	}
}

func (tvc *tvController) PostText(ctx *gin.Context) {
	ctx.JSON(200, gin.H{})
	text, err := ctx.GetRawData()
	if err != nil {
		println("error " + err.Error()) // TODO needs logger
		return
	}

	go func(msg string, t *tvController) {
		tvc.discordService.SendTextMessage(msg)
	}(string(text), tvc)
}

func (tvc *tvController) PostJson(ctx *gin.Context) {
	ctx.JSON(200, gin.H{})
	var message types.JSONMessageBody
	err := ctx.ShouldBindJSON(&message)
	if err != nil {
		println("error " + err.Error()) // TODO needs logger
		return
	}
	go func(msg types.JSONMessageBody, t *tvController) {
		t.discordService.SendJSONMessageToAltSignals(msg)
	}(message, tvc)
}
