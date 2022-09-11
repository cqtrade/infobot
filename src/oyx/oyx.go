package oyx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cqtrade/infobot/src/config"
	"github.com/cqtrade/infobot/src/notification"
	"github.com/cqtrade/infobot/src/state"
)

type Oyx struct {
	cfg        config.Config
	notif      notification.Notification
	appState   state.State
	httpClient *http.Client
}

func New(cfg config.Config, notif notification.Notification, appState state.State) *Oyx {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	return &Oyx{
		cfg:        cfg,
		notif:      notif,
		appState:   appState,
		httpClient: httpClient,
	}
}

func (o *Oyx) Run() {
	// keep logic in separate functions
	for {
		// check if buy allowed
		// get subaccount from conf
		// get portfolio from conf
		// check if any bids open for portfolio
		// cancel expired bids
		// get latest prices every minute
		// calculate rocp from close for period
		// check balances
		// divide free balances
		// if available funds and no bids:
		// market buy some and set bids lower

		// rinse and repeat

		t := time.Now()
		fmt.Println("Fine now", t)
		time.Sleep(time.Second)
	}
}
