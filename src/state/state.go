package state

import (
	"errors"
	"fmt"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/types"
)

// https://gobyexample.com/stateful-goroutines

type State struct {
	cfg                 config.Config
	notif               notification.Notification
	PriceReads          chan types.ReadPriceOp
	PriceWrites         chan types.WritePriceOp
	PositionsInfoReads  chan types.ReadPositionsInfo
	PositionsInfoWrites chan types.WritePositionsInfo
}

func New(cfg config.Config, notif notification.Notification) *State {
	return &State{
		cfg:                 cfg,
		notif:               notif,
		PriceReads:          make(chan types.ReadPriceOp),
		PriceWrites:         make(chan types.WritePriceOp),
		PositionsInfoReads:  make(chan types.ReadPositionsInfo),
		PositionsInfoWrites: make(chan types.WritePositionsInfo),
	}
}

func (s *State) RunStateLatestPrices() {
	latestPrices := make(map[string]types.ValAt)
	for {
		select {
		case read := <-s.PriceReads:
			read.Resp <- latestPrices[read.Key]
		case write := <-s.PriceWrites:
			latestPrices[write.Key] = types.ValAt{Price: write.Val.Price, At: write.Val.At}
			write.Resp <- true
		}
	}
}

func (s *State) ReadLatestPriceForMarket(market string) (float64, error) {
	var err error
	var latestPrice float64
	read := types.ReadPriceOp{
		Key:  market,
		Resp: make(chan types.ValAt)}

	s.PriceReads <- read

	priceAt := <-read.Resp

	shouldRequestFromRest := false
	if priceAt.Price == 0 {
		shouldRequestFromRest = true
		s.notif.Log("WARNING", "ReadLatestPriceForMarket - price is 0 get new reading via rest")
	}

	if priceAt.At == 0 {
		shouldRequestFromRest = true
		s.notif.Log("WARNING", "ReadLatestPriceForMarket - price is 0 get new reading via rest")
	}

	elapsed := time.Now().Unix() - priceAt.At

	if elapsed > 5 {
		shouldRequestFromRest = true
		s.notif.Log("WARNING", "ReadLatestPriceForMarket - Price older than 5 secs", fmt.Sprintf("%d", time.Now().Unix()-priceAt.At))
	}

	latestPrice = priceAt.Price
	if shouldRequestFromRest {
		s.notif.Log("DEBUG", "ReadLatestPriceForMarket - request from REST API")
		key := s.cfg.FTXKey
		secret := s.cfg.FTXSecret
		client := ftx.New(key, secret, "")
		defer client.Client.CloseIdleConnections()

		candles, err := client.GetHistoricalPriceLatest(market, 15, 1)
		if err != nil {
			return latestPrice, err
		}

		if !candles.Success {
			s.notif.Log("ERROR", "ReadLatestPriceForMarket rest unsuccessful", candles.HTTPCode, candles.ErrorMessage)
			return latestPrice, errors.New(fmt.Sprintf("%d ", candles.HTTPCode) + candles.ErrorMessage)
		}

		if len(candles.Result) > 0 {
			candle := candles.Result[0]
			s.notif.Log("", market, latestPrice, candle.Close)
			latestPrice = candle.Close
		}
	}

	return latestPrice, err
}

func (s *State) RunPositionsInfo() {
	positionsInfo := make(map[string]types.PositionInfo)
	for {
		select {
		case r := <-s.PositionsInfoReads:
			r.Resp <- positionsInfo
		case w := <-s.PositionsInfoWrites:
			positionsInfo[w.Key] = types.PositionInfo{
				Side:        w.PositionInfo.Side,
				Stops:       w.PositionInfo.Stops,
				TakeProfits: w.PositionInfo.TakeProfits,
			}
			w.Resp <- true
		}
	}
}
