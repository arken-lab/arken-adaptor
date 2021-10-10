package universaladaptor

const (
	EventTypeTradeTx     EventType = 0
	EventTypePoolCreated EventType = 1
	EventTypePoolUpdated EventType = 2
	EventTypePoolDeleted EventType = 3
)

type Event struct {
	Type    EventType
	payload interface{}
}

type EventType int
