package universaladaptor

type Pool struct{}

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
