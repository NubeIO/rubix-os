package broadcast_model

import (
	"github.com/NubeIO/flow-framework/plugin/plugin-api"
)

// Message is a message wrapper with the channel, sender and recipient.
type Message struct {
	Msg         plugin.Message
	ChannelName string
	IsSend      bool
}
