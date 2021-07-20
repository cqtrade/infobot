package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/types"
)

type Notification struct {
	cfg        config.Config
	httpClient *http.Client
}

func New(cfg config.Config) *Notification {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	return &Notification{
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

func (ds *Notification) sendNotification(ch string, message string) {
	if !ds.cfg.GetDiscordEnabled() {
		return
	}

	message = time.Now().Format("Jan 02 15:04:01") + " " + message
	go reqToDiscord(ch, message, ds.httpClient)
}

func (ds *Notification) SendTextMessage(msg string) {
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

func (ds *Notification) SendJSONMessageToAltSignals(msgJSON types.JSONMessageBody) {
	m := "**" + msgJSON.Ticker + "**"

	switch s := msgJSON.Signal; s {
	case 1001:
		m += " Flash BUY"
	case -2002:
		m += " take profit sell"
	case 1:
		m += " Long"
	case 2:
		m += " Exit Long"
	case -1:
		m += " Short"
	case -2:
		m += " Exit Short"

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
func (ds *Notification) SendFlashMessage(msgJSON types.JSONMessageBody) {
	// TODO format message
	m := "**" + msgJSON.Ticker + "**"

	switch s := msgJSON.Signal; s {
	case 1001:
		m += " Flash BUY"
	case -2002:
		m += " take profit sell"
	case 1:
		m += " Long"
	case 2:
		m += " Exit Long"
	case -1:
		m += " Short"
	case -2:
		m += " Exit Short"

	default:
		m += " unknown signal: " + fmt.Sprintf("%g", s)
	}

	ch := ds.cfg.GetDiscordChByChName("flash")
	if ch != "" {
		ds.sendNotification(ch, m)
		return
	}
}
