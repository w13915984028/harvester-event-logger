package eventlogger

import (
	"context"
	"fmt"

  ctlcorev1 "github.com/rancher/wrangler/pkg/generated/controllers/core/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type Handler struct {
	event   ctlcorev1.EventController
}

func Register(ctx context.Context, management *config.Management, opts config.Options) error {
	eventController := management.CoreFactory.Core().V1().Event()
	h := &Handler{
		event:   eventController,
	}

	eventController.OnChange(ctx, "event-logger", h.OnEventChange)
	return nil
}

// log the event, it will be fetched by fluentbit&fluentd
func (h *Handler) OnEventChange(key string, event *corev1.Event) (*corev1.Event, error) {
	if event == nil || event.DeletionTimestamp != nil {
		return nil, nil
	}

  // log the related event
  logrus.Infof("%s", fmt.Sprint(event))

	return event, nil
}
