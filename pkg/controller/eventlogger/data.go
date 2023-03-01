package eventlogger

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type EventData struct {
	UID   types.UID `json:"uid,omitempty"`
	Verb  string    `json:"verb"`
	Event *v1.Event `json:"event"`
	//	OldEvent *v1.Event `json:"old_event,omitempty"`
}

func NewEventData(event *v1.Event) EventData {
	if event.Count == 1 {
		return EventData{
			UID:   event.UID,
			Verb:  "ADDED",
			Event: event,
		}
	} else {
		return EventData{
			UID:   event.UID,
			Verb:  "UPDATED",
			Event: event,
		}
	}
}
