package controller

import (
	"github.com/cqtrade/infobot/src/service"
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
	tvc.discordService.SendTextMessage(string(text))
}

type Signal struct {
	Ticker   string `json:"ticker"`
	Exchange string `json:"exchange"`
	Signal   string `json:"signal"`
	Type     string `json:"type"`
}

func (tvc *tvController) PostJson(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"post": "json"})
}
