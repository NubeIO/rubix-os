package plugin

import (
	"time"

	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/flow-framework/plugin/compat"
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
