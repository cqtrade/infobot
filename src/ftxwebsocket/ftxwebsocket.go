package ftxwebsocket

import (
	"context"
	"fmt"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/state"
	"github.com/wzbear/go-ftx/realtime"
)

type FtxWebSocket struct {
	cfg config.Config
	st  state.State
}

func New(cfg config.Config, st state.State) *FtxWebSocket {
	return &FtxWebSocket{
		cfg: cfg,
		st:  st,
	}
}

func (ftws *FtxWebSocket) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan realtime.Response)
	go realtime.Connect(ctx, ch, []string{"ticker"}, []string{"BTC/USD", "ETH/USD", ftws.cfg.FutureBTC, ftws.cfg.FutureETH}, nil)

	for {
		select {
		case v := <-ch:
			switch v.Type {
			case realtime.TICKER:
				write := state.WritePriceOp{
					Key: v.Symbol,
					Val: state.ValAt{
						Price: v.Ticker.Last,
						At:    time.Now().Unix(),
					},
					Resp: make(chan bool)}
				ftws.st.PriceWrites <- write
				<-write.Resp

			case realtime.UNDEFINED:
				fmt.Printf("UNDEFINED %s	%s\n", v.Symbol, v.Results.Error())
			}
		}
	}
}
