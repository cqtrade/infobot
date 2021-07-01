package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/types"
)

type DiscordService interface {
	SendTextMessage(msg string)
	sendNotification(ch string, message string)
	SendJSONMessageToAltSignals(msgJSON types.JSONMessageBody)
	SendFlashMessage(msgJSON types.JSONMessageBody)
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

func reqToDiscord(webhookUrl string, msg string, client *http.Client) { // maybe not a good idea to share same http client between goroutines

	reqBody := types.NotificationBody{Content: msg}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		println("ERROR json.Marshal(reqBody)" + err.Error())
		return
	}

	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(reqBodyBytes))

	if err != nil {
		println("ERROR preparing discord payload" + err.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	defer client.CloseIdleConnections()

	if err != nil {
		println("ERROR logger http " + err.Error())
		return
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		println("resp.StatusCode: " + fmt.Sprintf("%d", resp.StatusCode))
		return
	}
}

func (ds *discordService) sendNotification(ch string, message string) {
	if !ds.cfg.GetDiscordEnabled() {
		return
	}

	go reqToDiscord(ch, message, ds.httpClient)
}

func (ds *discordService) SendTextMessage(msg string) {
	ch := ds.cfg.GetDiscordChByChName("random-ideas")
	if ch != "" {
		ds.sendNotification(ch, msg)
		return
	}
}

/*
1 oversold LTF
2 breakout LTF
11 oversold HTF
21 breakout HTF

-1 overbought LTF
-2 breakdown LTF
-11 overbought HTF
-21 breakdown HTF
*/

func (ds *discordService) SendJSONMessageToAltSignals(msgJSON types.JSONMessageBody) {
	m := "**" + msgJSON.Ticker + "**"

	switch s := msgJSON.Signal; s {
	case 1:
		m += " check **BUY** - oversold LTF"
	case 2:
		m += " check **BUY** - breakout LTF"
	case 11:
		m += " check **BUY** - oversold **HTF**"
	case 21:
		m += " check **BUY** - breakout **HTF**"
	case -1:
		m += " check **SELL** - overbought LTF"
	case -2:
		m += " check **SELL** - breakdown LTF"
	case -11:
		m += " check **SELL** - overbought **HTF**"
	case -21:
		m += " check **SELL** - breakdown **HTF**"

	default:
		m += " unknown signal: " + fmt.Sprintf("%g", s)
	}

	m += " " + msgJSON.Exchange

	ch := ds.cfg.GetDiscordChByChName("alt-signals")
	if ch != "" {
		ds.sendNotification(ch, m)
		return
	}
}

/**
flash signal gets signals from tested strategies
*/
func (ds *discordService) SendFlashMessage(msgJSON types.JSONMessageBody) {
	// TODO format message
	m := "**" + msgJSON.Ticker + "**"

	switch s := msgJSON.Signal; s {
	case 1:
		m += " check **BUY** - oversold LTF"
	case 2:
		m += " check **BUY** - breakout LTF"
	case 11:
		m += " check **BUY** - oversold **HTF**"
	case 21:
		m += " check **BUY** - breakout **HTF**"
	case -1:
		m += " check **SELL** - overbought LTF"
	case -2:
		m += " check **SELL** - breakdown LTF"
	case -11:
		m += " check **SELL** - overbought **HTF**"
	case -21:
		m += " check **SELL** - breakdown **HTF**"

	default:
		m += " unknown signal: " + fmt.Sprintf("%g", s)
	}

	m += " " + msgJSON.Exchange

	ch := ds.cfg.GetDiscordChByChName("flash")
	if ch != "" {
		ds.sendNotification(ch, m)
		return
	}
}
