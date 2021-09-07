package ftxtrade

import (
	"errors"
	"fmt"

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
		return position, errors.New("No Success CheckFuturePosition " + positions.ErrorMessage + fmt.Sprintf(" %d", positions.HTTPCode))
	}

	for _, currPosition := range positions.Result {
		if currPosition.Future == future {
			position = currPosition
			return position, err
		}
	}
	return position, err
}
