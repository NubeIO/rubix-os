package eventbus

import (
	"context"
)

type (

	// NotificationService :nodoc:
	NotificationService interface {
		EmitString(ctx context.Context, topicName string, data string)
		Emit(ctx context.Context, topicName string, data interface{})
		RegisterTopic(topic string)
		RegisterTopicParent(parent string, child string)
	}
)
