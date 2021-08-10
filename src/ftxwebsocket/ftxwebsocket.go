package ftxwebsocket

import (
	"context"
	"fmt"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/state"
	"github.com/cqtrade/infobot/src/types"
	"github.com/wzbear/go-ftx/realtime"
)

type FtxWebSocket struct {
	cfg   config.Config
	notif notification.Notification
	st    state.State
}

func New(cfg config.Config, notif notification.Notification, st state.State) *FtxWebSocket {
	return &FtxWebSocket{
		cfg:   cfg,
		notif: notif,
		st:    st,
	}
}

func (ftws *FtxWebSocket) RunWebSocket() {
	time.Sleep(5 * time.Second)
	ftws.notif.Log("INFO", "ws START")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{"ticker"}, []string{"BTC/USD", "ETH/USD", ftws.cfg.FutureBTC, ftws.cfg.FutureETH}, nil)
	timeAfter := time.After(30 * time.Second)
	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.TICKER:
				timeAfter = time.After(30 * time.Second)
				write := types.WritePriceOp{
					Key: v.Symbol,
					Val: types.ValAt{
						Price: v.Ticker.Last,
						At:    time.Now().Unix(),
					},
					Resp: make(chan bool)}
				ftws.st.PriceWrites <- write
				<-write.Resp

			case realtime.UNDEFINED:
				fmt.Printf("UNDEFINED %s	%s\n", v.Symbol, v.Results.Error())
			default:
				fmt.Println(fmt.Sprintf("Default %v", v))
			}
		case <-timeAfter:
			ftws.notif.Log("ERROR", "ws timeout")
			go ftws.RunWebSocket()
			return

		case <-ctx.Done():
			ftws.notif.Log("ERROR", "ws ctx done")
			go ftws.RunWebSocket()
			return
		}
	}
}
