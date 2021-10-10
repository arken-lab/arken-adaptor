package universaladaptor

import (
	"errors"

	"github.com/holiman/uint256"
)

type Pool interface {
	GetChain() string
	GetAddress() string
	GetTokens() []Token
	GetExchangeRate(token Token, amount uint256.Int) uint256.Int
	GetMidPrice(quoteToken Token) (float64, error)
	GetReserve(token Token) (uint256.Int, error)
}

func NewPoolCreatedEvent(pool *Pool) Event {
	return Event{
		Type:    EventTypePoolCreated,
		payload: pool,
	}
}

func NewPoolUpdatedEvent(pool *Pool) Event {
	return Event{
		Type:    EventTypePoolUpdated,
		payload: pool,
	}
}

func NewPoolDeletedEvent(pool *Pool) Event {
	return Event{
		Type:    EventTypePoolDeleted,
		payload: pool,
	}
}

func (e *Event) GetPool() (*Pool, error) {
	if e.Type != EventTypePoolCreated && e.Type != EventTypePoolUpdated && e.Type != EventTypePoolDeleted {
		return nil, errors.New("only EventTypePoolCreated or EventTypePoolUpdated or EventTypePoolDeleted is allowed to use this method")
	}
	data, ok := e.payload.(*Pool)
	if !ok {
		return nil, errors.New("cannot convert payload to Pool")
	}
	return data, nil
}
