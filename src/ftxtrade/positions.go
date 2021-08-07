package ftxtrade

import (
	"errors"

	"github.com/cqtrade/infobot/src/ftx"
	"github.com/cqtrade/infobot/src/ftx/structs"
)

func (ft *FtxTrade) CheckFuturePosition(client *ftx.FtxClient, future string) (structs.Position, error) {
	var position structs.Position
	positions, err := client.GetPositions(true)
	if err != nil {
		return position, err
	}

	if !positions.Success {
		return position, errors.New("No Success CheckFuturePosition")
	}

	for _, currPosition := range positions.Result {
		if currPosition.Future == future {
			position = currPosition
			return position, err
		}
	}
	return position, err
}

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
