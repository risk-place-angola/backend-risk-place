package port

import "github.com/risk-place-angola/backend-risk-place/internal/domain/event"

type EventDispatcher interface {
	Register(eventName string, handler event.EventHandler)
	Dispatch(event event.Event)
}
