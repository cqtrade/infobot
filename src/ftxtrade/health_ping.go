package ftxtrade

import (
	"fmt"
	"time"

	"github.com/cqtrade/infobot/src/types"
)

func (ft *FtxTrade) RunHealthPing() {
	if ft.cfg.HealthLogEnabled {
		for t := range time.Tick(time.Second) {
			if t.Second() == 55 {
				message := t.Format("Jan 02 15:04") + " " + ft.GetOverview("test1")

				read := types.ReadPositionsInfo{
					Resp: make(chan map[string]types.PositionInfo)}

				ft.appState.PositionsInfoReads <- read

				positionsInfo := <-read.Resp
				message += "\n"
				for key, positionInfo := range positionsInfo {
					message += key + " " + fmt.Sprintf("%+v", positionInfo) + "\n"
				}

				write := types.WriteLogMessage{
					Val:  types.LogMessage{Message: message, Channel: ft.cfg.DiscordChHealth},
					Resp: make(chan bool)}
				ft.notif.ChLogMessageWrites <- write
				<-write.Resp
			}
		}
	}
}
