package lit

import (
	trader "github.com/mit-dci/lit-rpc-client-go-samples/dlcexchange/trader"
  orderbook "github.com/mit-dci/lit-rpc-client-go-samples/dlcexchange/orderbook"
)

const (
	coinType    uint32 = 1
	mHost       string = "127.0.0.1"
	mPort       int32  = 8001
	mListenPort uint32 = 2448
)

func handleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func GetBids() (bids []orderbook.Order) {
	m, err := trader.NewTrader("Market Maker", mHost, mPort, nil)
  handleError(err)
	return m.GetBids()
}

func GetAsks() (asks []orderbook.Order) {
  m, err := trader.NewTrader("Market Maker", mHost, mPort, nil)
  handleError(err)
	return m.GetAsks()
}
