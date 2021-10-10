package universaladaptor

import "errors"

type TradeTx struct{}

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
