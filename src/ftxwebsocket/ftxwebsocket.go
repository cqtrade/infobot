package ftxwebsocket

import (
	"context"
	"fmt"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/wzbear/go-ftx/realtime"
)

type FtxWebSocket struct {
	cfg   config.Config
	ft    ftxtrade.FtxTrade
	notif notification.Notification
}

func New(cfg config.Config, ft ftxtrade.FtxTrade, notif notification.Notification) *FtxWebSocket {
	return &FtxWebSocket{
		cfg:   cfg,
		ft:    ft,
		notif: notif,
	}
}

// https://gobyexample.com/stateful-goroutines

func (ftws *FtxWebSocket) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{"ticker"}, []string{"BTC/USD", "BTC-1231"}, nil)
	// go realtime.ConnectForPrivate(ctx, ch, "<key>", "<secret>", []string{"orders", "fills"}, nil)

	// go func() {
	// 	state := make(map[int]int)
	// 	for {
	// 		select {
	// 		case read := <-reads:
	// 			read.resp <- state[read.key]
	// 		case write := <-writes:
	// 			state[write.key] = write.val
	// 			write.resp <- true
	// 		}
	// 	}
	// }()

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.TICKER:
				fmt.Printf("%s	%+v\n", v.Symbol, v.Ticker)

			// case realtime.TRADES:
			// 	fmt.Printf("%s	%+v\n", v.Symbol, v.Trades)
			// 	for i := range v.Trades {
			// 		if v.Trades[i].Liquidation {
			// 			fmt.Printf("-----------------------------%+v\n", v.Trades[i])
			// 		}
			// 	}

			// case realtime.ORDERBOOK:
			// 	fmt.Printf("%s	%+v\n", v.Symbol, v.Orderbook)

			// case realtime.UNDEFINED:
			// 	fmt.Printf("%s	%s\n", v.Symbol, v.Results.Error())

			// case realtime.ORDERS:
			// 	fmt.Printf("%d	%+v\n", v.Type, v.Orders)

			// case realtime.FILLS:
			// 	fmt.Printf("%d	%+v\n", v.Type, v.Fills)

			case realtime.UNDEFINED:
				fmt.Printf("UNDEFINED %s	%s\n", v.Symbol, v.Results.Error())
			}
		}
	}
}
