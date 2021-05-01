package service

import "github.com/cqtrade/infobot/src/config"

type DiscordService interface {
	SendTextMessage(msg string)
}

type discordService struct {
	cfg config.Config
}

func NewDiscordService(cfg config.Config) DiscordService {
	return &discordService{
		cfg: cfg,
	}
}

func (ds *discordService) SendTextMessage(msg string) {
	println("Send")
	println(ds.cfg.GetDiscordChRandomIdeas())
	println(msg)
}
