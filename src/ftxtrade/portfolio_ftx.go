package ftxtrade

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cqtrade/infobot/src/ftx"
)

type PorftolioItem struct {
	Coin       string  `json:"coin,omitempty"`
	Allocation float64 `json:"allocation,omitempty"`
}

// signal comes buy with available money
func (ft *FtxTrade) PortfolioFTX(subAcc string) {

	byt := []byte(`[{"coin":"BTC","allocation":10.0},{"coin":"ETH","allocation":20.0}]`)
	var portfolioItems []PorftolioItem
	if err := json.Unmarshal(byt, &portfolioItems); err != nil {
		panic(err)
	}
	// fmt.Println(dat)
	// get usd balance
	key := ft.cfg.FTXKey
	secret := ft.cfg.FTXSecret
	client := ftx.New(key, secret, subAcc)

	sBalance, err := ft.CheckSpotBalance(client, subAcc, "USD")
	if err != nil {
		fmt.Println("Error getting balance", err)
		return
	}
	usdEquity := sBalance.Free
	usdEquity = 1000

	for _, portfolioItem := range portfolioItems {
		time.Sleep(time.Second)
		if portfolioItem.Coin == "USD" {
			continue
		}
		amountUsd := usdEquity * (portfolioItem.Allocation / 100)
		market := portfolioItem.Coin + "/USD"

		price, err := ft.appState.ReadLatestPriceForMarket(market)
		if err != nil {
			// fmt.Println("ERROR latest price", market, err.Error())
			ft.notif.Log("ERROR", "Read latest price. Abort.", market, err.Error())
		}
		positionSize := RoundDown(amountUsd/price, 4)
		fmt.Println(market, "Buy", positionSize)
		if positionSize < 0.0001 {
			ft.notif.Log("ERROR", "position size too low. Abort.", market, positionSize)
			continue
		}

		order, err := client.PlaceMarketOrder(market, "sell", "market", positionSize)
		if err != nil {
			ft.notif.Log("ERROR", "TpCoinBull Market BUY order. Abort.", market, err.Error())
			return
		}
		if !order.Success {
			ft.notif.Log("ERROR", "TpCoinBull  UNSUCCESSFUL", market, order.HTTPCode, order.ErrorMessage)
			return
		}
		time.Sleep(time.Second)
	}
	// num := dat["num"].(float64)
	// fmt.Println(num)

	/*
		var balance structs.SubaccountBalance
		sBalances, err := client.GetSubaccountBalances(subAcc)
		if err != nil {
			return balance, err
		}

		if !sBalances.Success {
			return balance, errors.New("No Success CheckSpotBalance")
		}

		for _, currBalance := range sBalances.Result {
			if currBalance.Coin == spot {
				balance = currBalance
				return balance, nil
			}
		}

		return balance, nil
	*/
}
