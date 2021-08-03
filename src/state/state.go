package state

import (
	"fmt"
	"time"

	"github.com/cqtrade/infobot/src/config"
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
	read := ReadPriceOp{
		Key:  market,
		Resp: make(chan ValAt)}

	s.PriceReads <- read

	priceAt := <-read.Resp

	if priceAt.Price == 0 {
		fmt.Println("Price is 0 get new reading via rest")
	}

	if priceAt.At == 0 {
		fmt.Println("Price time 0 get new reading via rest")
	}

	elapsed := time.Now().Unix() - priceAt.At

	if elapsed > 5 {
		fmt.Println("Price older than 5 secs get new reading via rest", fmt.Sprintf("%d", time.Now().Unix()-priceAt.At))
	}

	fmt.Println("Elapsed", fmt.Sprintf("%d", time.Now().Unix()-priceAt.At))

	// TODO check if price is too old
	return priceAt.Price, err
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
