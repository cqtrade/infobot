package notification

import (
	"bytes"
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

// https://motyar.github.io/golang-pretty-print-struct/
/*
func (p *pp) doPrint(a []interface{}) {
	prevString := false
	for argNum, arg := range a {
		isString := arg != nil && reflect.TypeOf(arg).Kind() == reflect.String
		// Add a space between two non-string arguments.
		if argNum > 0 && !isString && !prevString {
			p.buf.writeByte(' ')
		}
		p.printArg(arg, 'v')
		prevString = isString
	}
}
*/
func (ds *Notification) Log(level string, a ...interface{}) {
	var s []string
	s = append(s, time.Now().Format("Jan 02 15:04:05"))
	l := strings.ToUpper(level)
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
			s = append(s, fmt.Sprintf("%.2f", arg))
		case reflect.Struct:
			s = append(s, fmt.Sprintf("%+v", arg))
		default:
			fmt.Println("Not defined type", reflect.TypeOf(arg).Kind(), arg)
		}
	}

	message := strings.Join(s[:], " ")
	println(message)
	if l == "INFO" || l == "ERROR" {
		fmt.Println("TODO send to discord: ", message)
	} else {
		fmt.Println(message)
	}
}

// fmt.Srintf("%+v\n", p) //With name and value
// fmt.Srintf("%#v", p) //with name, value and type
