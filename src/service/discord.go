package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/cqtrade/infobot/src/config"
)

type DiscordService interface {
	SendTextMessage(msg string)
	sendNotification(ch string, message string)
}

type discordService struct {
	cfg        config.Config
	httpClient *http.Client
}

func NewDiscordService(cfg config.Config) DiscordService {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	return &discordService{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

func (ds *discordService) sendNotification(ch string, message string) {
	if !ds.cfg.GetDiscordEnabled() {
		return
	}

	go func(webhookUrl string, msg string, client *http.Client) {
		type notificationBody struct {
			Content string `json:"content"`
		}

		reqBody := notificationBody{Content: msg}

		reqBodyBytes, _ := json.Marshal(reqBody)
		req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(reqBodyBytes))
		if err != nil {
			println("ERROR preparing discord payload" + err.Error())
		} else {
			req.Header.Add("Content-Type", "application/json")
			resp, err := client.Do(req)
			if resp.StatusCode == 204 && err == nil {
				return
			}
			println("ERROR logger http " + err.Error())
		}
	}(ch, message, ds.httpClient)
}

func (ds *discordService) SendTextMessage(msg string) {
	ds.sendNotification(ds.cfg.GetDiscordChRandomIdeas(), msg)
}
