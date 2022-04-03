package model

// Message is a message wrapper with the channel, sender and recipient.
type Message struct {
	Msg         plugin.Message
	ChannelName string
	IsSend      bool
}
