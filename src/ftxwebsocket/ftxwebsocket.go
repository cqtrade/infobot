package ftxwebsocket

import (
	"context"
	"fmt"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/ftxtrade"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/wzbear/go-ftx/realtime"
)

type ReadPriceOp struct {
	key  string
	resp chan float64
}
type WritePriceOp struct {
	key  string
	val  float64
	resp chan bool
}
type FtxWebSocket struct {
	cfg    config.Config
	ft     ftxtrade.FtxTrade
	notif  notification.Notification
	reads  chan ReadPriceOp
	writes chan WritePriceOp
}

func New(cfg config.Config, ft ftxtrade.FtxTrade, notif notification.Notification) *FtxWebSocket {
	return &FtxWebSocket{
		cfg:    cfg,
		ft:     ft,
		notif:  notif,
		reads:  make(chan ReadPriceOp),
		writes: make(chan WritePriceOp),
	}
}

func (ftws *FtxWebSocket) StateLatestPrices() {
	latestPrices := make(map[string]float64)
	for {
		select {
		case read := <-ftws.reads:
			read.resp <- latestPrices[read.key]
		case write := <-ftws.writes:
			latestPrices[write.key] = write.val
			write.resp <- true
		}
	}
}

func (ftws *FtxWebSocket) ReadPriceState() {
	time.Sleep(time.Second * 3)
	for {
		readF := ReadPriceOp{
			key:  "BTC-1231",
			resp: make(chan float64)}
		ftws.reads <- readF
		fPrice := <-readF.resp

		readS := ReadPriceOp{
			key:  "BTC/USD",
			resp: make(chan float64)}
		ftws.reads <- readS
		sPrice := <-readS.resp

		// fmt.Println("BTC-1231\t", fPrice, "\tBTC/USD\t", sPrice)
		arbBtc := fPrice*100/sPrice - 100
		fmt.Println(fmt.Sprintf("Diff %.2f%%", arbBtc))
		// TODO if less than zero log
		time.Sleep(time.Second)
	}
}

// https://gobyexample.com/stateful-goroutines

func (ftws *FtxWebSocket) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{"ticker"}, []string{"BTC/USD", "BTC-1231", "ETH/USD", "ETH-1231"}, nil)
	// go realtime.ConnectForPrivate(ctx, ch, "<key>", "<secret>", []string{"orders", "fills"}, nil)

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.TICKER:
				write := WritePriceOp{
					key:  v.Symbol,
					val:  v.Ticker.Last,
					resp: make(chan bool)}
				ftws.writes <- write
				<-write.resp

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
