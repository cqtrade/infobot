package notification

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/types"
)

type Notification struct {
	cfg                config.Config
	httpClient         *http.Client
	ChLogMessageReads  chan types.ReadLogMessage
	ChLogMessageWrites chan types.WriteLogMessage
}

func New(cfg config.Config) *Notification {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	return &Notification{
		cfg:                cfg,
		httpClient:         httpClient,
		ChLogMessageReads:  make(chan types.ReadLogMessage),
		ChLogMessageWrites: make(chan types.WriteLogMessage),
	}
}

func (ds *Notification) RunStateLogMessages() {
	queue := list.New()
	for {
		select {
		case read := <-ds.ChLogMessageReads:
			e := queue.Front()
			if e != nil {
				queue.Remove(e)
				read.Resp <- e.Value.(types.LogMessage)
			} else {
				var l types.LogMessage
				read.Resp <- l
			}

		case write := <-ds.ChLogMessageWrites:
			queue.PushBack(write.Val)
			write.Resp <- true
		}
	}
}

func (ds *Notification) RunReadLogMessages() {
	for {
		time.Sleep(3 * time.Second)
		read := types.ReadLogMessage{Resp: make(chan types.LogMessage)}
		ds.ChLogMessageReads <- read
		res := <-read.Resp
		if res.Channel == "" || res.Message == "" {
			fmt.Println("empty")
		} else {
			fmt.Println(fmt.Sprintf("%+v", res))
		}
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

func (ds *Notification) Log(level string, a ...interface{}) {
	var s []string
	s = append(s, time.Now().Format("Jan 02 15:04:05"))
	var l string

	if level == "" {
		l = "DEBUG"
	} else {
		l = strings.ToUpper(level)
	}

	s = append(s, l)
	for _, arg := range a {
		if arg == nil {
			continue
		}

		switch reflect.TypeOf(arg).Kind() {
		case reflect.String:
			s = append(s, arg.(string))
		case reflect.Int:
			s = append(s, fmt.Sprintf("%d", arg))
		case reflect.Float64:
			s = append(s, fmt.Sprintf("%.4f", arg))
		case reflect.Struct:
			s = append(s, fmt.Sprintf("%+v", arg))
		default:
			fmt.Println("Not defined type", reflect.TypeOf(arg).Kind(), arg)
		}
	}

	message := strings.Join(s[:], " ")

	if l == "INFO" || l == "ERROR" {
		write := types.WriteLogMessage{
			Val:  types.LogMessage{Message: message, Channel: ds.cfg.DiscordChLogs},
			Resp: make(chan bool)}
		ds.ChLogMessageWrites <- write
		<-write.Resp
	} else {
		fmt.Println(message)
	}
}
