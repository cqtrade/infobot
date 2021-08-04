package state

import (
	"fmt"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/notification"
)

// https://gobyexample.com/stateful-goroutines

type ValAt struct {
	Price float64
	At    int64
}
type ReadPriceOp struct {
	Key  string
	Resp chan ValAt
}

type WritePriceOp struct {
	Key  string
	Val  ValAt
	Resp chan bool
}

type State struct {
	cfg         config.Config
	notif       notification.Notification
	PriceReads  chan ReadPriceOp
	PriceWrites chan WritePriceOp
}

func New(cfg config.Config, notif notification.Notification) *State {
	return &State{
		cfg:         cfg,
		notif:       notif,
		PriceReads:  make(chan ReadPriceOp),
		PriceWrites: make(chan WritePriceOp),
	}
}

func (s *State) StateLatestPrices() {
	latestPrices := make(map[string]ValAt)
	for {
		select {
		case read := <-s.PriceReads:
			read.Resp <- latestPrices[read.Key]
		case write := <-s.PriceWrites:
			latestPrices[write.Key] = ValAt{Price: write.Val.Price, At: write.Val.At}
			write.Resp <- true
		}
	}
}

func (s *State) ReadLatestPriceForMarket(market string) (float64, error) {
	var err error
	var latestPrice float64
	read := ReadPriceOp{
		Key:  market,
		Resp: make(chan ValAt)}

	s.PriceReads <- read

	priceAt := <-read.Resp

	shouldRequestFromRest := false
	if priceAt.Price == 0 {
		shouldRequestFromRest = true
		s.notif.Log("ERROR", "ReadLatestPriceForMarket - price is 0 get new reading via rest")
	}

	if priceAt.At == 0 {
		shouldRequestFromRest = true
		s.notif.Log("ERROR", "ReadLatestPriceForMarket - price is 0 get new reading via rest")
	}

	elapsed := time.Now().Unix() - priceAt.At

	if elapsed > 5 {
		shouldRequestFromRest = true
		s.notif.Log("ERROR", "ReadLatestPriceForMarket - Price older than 5 secs", fmt.Sprintf("%d", time.Now().Unix()-priceAt.At))
	}

	latestPrice = priceAt.Price
	if shouldRequestFromRest {
		s.notif.Log("INFO", "ReadLatestPriceForMarket - request from REST API")
		key := s.cfg.FTXKey
		secret := s.cfg.FTXSecret
		client := ftx.New(key, secret, "")
		defer client.Client.CloseIdleConnections()

		candles, err := client.GetHistoricalPriceLatest(market, 15, 1)
		if err != nil {
			return latestPrice, err
		}

		if len(candles.Result) > 0 {
			candle := candles.Result[0]
			s.notif.Log("", market, latestPrice, candle.Close)
			latestPrice = candle.Close
		}
	}

	return latestPrice, err
}

func (s *State) ReadPriceState() {
	time.Sleep(time.Second * 3)
	for {

		latestBTCf, _ := s.ReadLatestPriceForMarket(s.cfg.FutureBTC)
		latestBTCs, _ := s.ReadLatestPriceForMarket("BTC/USD")
		latestETHf, _ := s.ReadLatestPriceForMarket(s.cfg.FutureETH)
		latestETHs, _ := s.ReadLatestPriceForMarket("ETH/USD")

		fmt.Println("BTC/USD\t\t", fmt.Sprintf("%.2f", latestBTCs), "\tETH/USD\t\t", fmt.Sprintf("%.2f", latestETHs))
		// fmt.Println(s.cfg.FutureBTC, "\t", fmt.Sprintf("%.2f", latestBTCf), "\t", s.cfg.FutureETH, "\t", fmt.Sprintf("%.2f", latestETHf))
		fmt.Println("BTC premium\t\t", fmt.Sprintf("%.2f%%", latestBTCf*100/latestBTCs-100), "\tETH premium\t\t", fmt.Sprintf("%.2f%%", latestETHf*100/latestETHs-100))
		fmt.Println(time.Now().Unix())
		time.Sleep(time.Second * 5)
	}
}
