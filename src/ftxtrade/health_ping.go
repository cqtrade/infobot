package ftxtrade

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cqtrade/infobot/src/types"
)

func (ft *FtxTrade) StartHealthPing() {
	for t := range time.Tick(time.Second) {
		if t.Second() == 33 {
			reqBody := types.NotificationBody{Content: t.Format("Jan 02 15:04") + ft.GetOverview("test1")}

			reqBodyBytes, err := json.Marshal(reqBody)
			if err != nil {
				println("ERROR json.Marshal(reqBody)" + err.Error())
				ft.StartHealthPing()
				return
			}

			req, err := http.NewRequest(http.MethodPost, ft.cfg.DiscordChHealth, bytes.NewBuffer(reqBodyBytes))

			if err != nil {
				println("ERROR preparing discord payload" + err.Error())
				ft.StartHealthPing()
				return
			}

			req.Header.Add("Content-Type", "application/json")

			resp, err := ft.httpClient.Do(req)
			defer ft.httpClient.CloseIdleConnections()

			if err != nil {
				println("ERROR logger http " + err.Error())
				ft.StartHealthPing()
				return
			}

			if resp.StatusCode < 200 || resp.StatusCode > 299 {
				println("Discord resp.StatusCode: " + fmt.Sprintf("%d", resp.StatusCode))
				ft.StartHealthPing()
				return
			}
		}
	}
}
