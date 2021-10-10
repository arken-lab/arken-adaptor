package universaladaptor

import "errors"

const (
	TradeTxBuy  TradeTxSide = "BUY"
	TradeTxSell TradeTxSide = "SELL"
)

type TradeTxToken struct {
	Chain   string
	Address string
}

type TradeTxSide string

type TradeTx struct {
	Token        TradeTxToken
	TxHash       string
	LogIndex     int64
	Block        int64
	LpAddress    string
	AmountBase   float64
	AmountUsd    float64
	PriceBaseUsd float64
	Side         TradeTxSide
	Timestamp    int64
	IsRemoved    bool
	IsHistorical bool
}

func NewTradeTxEvent(tradeTx *TradeTx) Event {
	return Event{
		Type:    EventTypeTradeTx,
		payload: tradeTx,
	}
}

func (e *Event) GetTradeTx() (*TradeTx, error) {
	if e.Type != EventTypeTradeTx {
		return nil, errors.New("only EventTypeTradeTx is allowed to use this method")
	}
	data, ok := e.payload.(*TradeTx)
	if !ok {
		return nil, errors.New("cannot convert payload to TradeTx")
	}
	return data, nil
}
