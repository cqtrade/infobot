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

func reqToDiscord(webhookUrl string, msg string, client *http.Client) {

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

func (ds *Notification) RunReadLogMessages() {
	for {
		read := types.ReadLogMessage{Resp: make(chan types.LogMessage)}
		ds.ChLogMessageReads <- read
		res := <-read.Resp
		if res.Channel == "" || res.Message == "" {
			continue
		}
		fmt.Println(res.Message)
		if !ds.cfg.DiscordEnabled {
			return
		}
		reqToDiscord(res.Channel, res.Message, ds.httpClient)
		time.Sleep(time.Second)
	}
}

func (ds *Notification) SendTextMessage(msg string) {
	write := types.WriteLogMessage{
		Val:  types.LogMessage{Message: msg, Channel: ds.cfg.DiscordChRandomIdeas},
		Resp: make(chan bool)}
	ds.ChLogMessageWrites <- write
	<-write.Resp
}

func (ds *Notification) SendFlashMessage(msgJSON types.JSONMessageBody) {
	m := "**" + msgJSON.Ticker + "**"

	switch s := msgJSON.Signal; s {
	case 1001:
		m += " BUY CRYPTO, partial"
	case -2002:
		m += " Take partial profit, sell crypto"
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

	write := types.WriteLogMessage{
		Val:  types.LogMessage{Message: m, Channel: ds.cfg.DiscordChFlash},
		Resp: make(chan bool)}
	ds.ChLogMessageWrites <- write
	<-write.Resp
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
		fmt.Println(message)
		write := types.WriteLogMessage{
			Val:  types.LogMessage{Message: message, Channel: ds.cfg.DiscordChLogs},
			Resp: make(chan bool)}
		ds.ChLogMessageWrites <- write
		<-write.Resp
	} else {
		fmt.Println(message)
	}
}
