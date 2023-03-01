package eventlogger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ctlcorev1 "github.com/rancher/wrangler/pkg/generated/controllers/core/v1"
	"github.com/w13915984028/harvester-event-logger/pkg/config"
	corev1 "k8s.io/api/core/v1"
)

type Handler struct {
	event ctlcorev1.EventController
}

func Register(ctx context.Context, management *config.Management) error {
	eventController := management.CoreFactory.Core().V1().Event()
	h := &Handler{
		event: eventController,
	}

	// why OnChange instead of Watch ?
	// due to the latency of the deployment of this POD, the watch may not get all events
	// the OnChange may cause some events are recorded multi times
	// more seems better than less
	eventController.OnChange(ctx, "event-logger", h.OnEventChange)
	return nil
}

// log the event, it will be fetched by fluentbit&fluentd
func (h *Handler) OnEventChange(key string, event *corev1.Event) (*corev1.Event, error) {
	if event == nil || event.DeletionTimestamp != nil {
		return nil, nil
	}

	ed := NewEventData(event)

	if dt, err := json.Marshal(ed); err == nil {
		fmt.Println(string(dt))
	} else {
		fmt.Fprintf(os.Stderr, "Failed to marshal json, err: %v, event: %v", err, event)
	}

	return event, nil
}
