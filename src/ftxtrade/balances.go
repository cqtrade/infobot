package ftxtrade

import (
	"fmt"

	"github.com/cqtrade/infobot/src/ftx"
)

func (ft *FtxTrade) CheckBalanceUSD(subAcc string, client *ftx.FtxClient) (float64, error) {
	sBalances, err := client.GetSubaccountBalances(subAcc)
	equity := 0.0
	if err != nil {
		return equity, err
	}

	if sBalances.Success {
		for _, balance := range sBalances.Result {
			if balance.Coin == "USD" {
				equity = balance.Free
			}
			fmt.Println("Coin\t", balance.Coin, "\tfree:\t", balance.Free, "\ttotal:\t", balance.Total)
		}
	}

	return equity, nil
}

func (ft *FtxTrade) CheckFreeBalanceETHBULL(subAcc string, client *ftx.FtxClient) (float64, error) {
	sBalances, err := client.GetSubaccountBalances(subAcc)
	equity := 0.0
	if err != nil {
		return equity, err
	}

	if sBalances.Success {
		for _, balance := range sBalances.Result {
			if balance.Coin == "ETHBULL" {
				equity = balance.Free
			}
			fmt.Println("Coin\t", balance.Coin, "\tfree:\t", balance.Free, "\ttotal:\t", balance.Total)
		}
	}

	return equity, nil
}

func (ft *FtxTrade) CheckFreeBalanceBTC(subAcc string, client *ftx.FtxClient) (float64, error) {
	sBalances, err := client.GetSubaccountBalances(subAcc)
	equity := 0.0
	if err != nil {
		return equity, err
	}

	if sBalances.Success {
		for _, balance := range sBalances.Result {
			if balance.Coin == "BTC" {
				equity = balance.Free
			}
		}
	}

	return equity, nil
}
