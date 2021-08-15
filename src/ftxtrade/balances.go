package ftxtrade

import (
	"errors"

	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/ftx/structs"
)

func (ft *FtxTrade) CheckSpotBalance(client *ftx.FtxClient, subAcc string, spot string) (structs.SubaccountBalance, error) {
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
}
