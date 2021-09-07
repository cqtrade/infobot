package ftx

import (
	"log"

	"github.com/cqtrade/infobot/src/ftx/structs"
)

type Positions structs.Positions

func (client *FtxClient) GetPositions(showAvgPrice bool) (Positions, error) {
	var positions Positions
	var url string
	if showAvgPrice {
		url = "positions?showAvgPrice=true"
	} else {
		url = "positions"
	}
	resp, err := client._getRetry(2, url, []byte(""))
	if err != nil {
		log.Printf("Error GetPositions", err)
		return positions, err
	}
	err = _processResponse(resp, &positions)
	positions.HTTPCode = resp.StatusCode
	return positions, err
}
