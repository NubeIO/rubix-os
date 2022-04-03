package plugin

import (
	"time"

	"github.com/NubeIO/flow-framework/plugin/compat"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type redirectToChannel struct {
	Messages chan model.MessageExternal
}

// SendMessage sends a message to the underlying message channel.
func (c redirectToChannel) SendMessage(msg compat.Message) error {
	c.Messages <- model.MessageExternal{
		Message:  msg.Message,
		Title:    msg.Title,
		Priority: msg.Priority,
		Date:     time.Now(),
		Extras:   msg.Extras,
	}
	return nil
}
