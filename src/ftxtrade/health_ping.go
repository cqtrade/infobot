package ftxtrade

import (
	"time"

	"github.com/cqtrade/infobot/src/types"
)

func (ft *FtxTrade) RunHealthPing() {
	if ft.cfg.HealthLogEnabled {
		for t := range time.Tick(time.Second) {
			if t.Second() == 33 {
				message := t.Format("Jan 02 15:04") + ft.GetOverview("test1")
				write := types.WriteLogMessage{
					Val:  types.LogMessage{Message: message, Channel: ft.cfg.DiscordChHealth},
					Resp: make(chan bool)}
				ft.notif.ChLogMessageWrites <- write
				<-write.Resp
			}
		}
	}
}
