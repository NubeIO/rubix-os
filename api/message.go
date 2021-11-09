package api

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/NubeDev/flow-framework/model"
	"github.com/NubeDev/location"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// The MessageDatabase interface for encapsulating database access.
type MessageDatabase interface {
	CreateMessage(message *model.Message) error
	GetMessages() ([]*model.Message, error)
	GetMessageByID(id uint) (*model.Message, error)
	GetMessagesSince(limit int, since uint) ([]*model.Message, error)
	DeleteMessages() error
	DeleteMessageByID(id uint) error
}

var timeNow = time.Now

// Notifier notifies when a new message was created.
type Notifier interface {
	Notify(message *model.MessageExternal)
}

// The MessageAPI provides handlers for managing messages.
type MessageAPI struct {
	DB       MessageDatabase
	Notifier Notifier
}

type pagingParams struct {
	Limit int  `form:"limit" binding:"min=1,max=200"`
	Since uint `form:"since" binding:"min=0"`
}

func (a *MessageAPI) CreateMessage(ctx *gin.Context) {
	message := model.MessageExternal{}
	if err := ctx.Bind(&message); err == nil {
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		message.Date = timeNow()
		message.ID = 0
		msgInternal := toInternalMessage(&message)
		if success := successOrAbort(ctx, 500, a.DB.CreateMessage(msgInternal)); !success {
			return
		}
		a.Notifier.Notify(toExternalMessage(msgInternal))
		ctx.JSON(200, toExternalMessage(msgInternal))
	}
}

func (a *MessageAPI) GetMessages(ctx *gin.Context) {
	withPaging(ctx, func(params *pagingParams) {
		// the +1 is used to check if there are more messages and will be removed on buildWithPaging
		messages, err := a.DB.GetMessagesSince(params.Limit+1, params.Since)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		ctx.JSON(200, buildWithPaging(ctx, params, messages))
	})
}

func buildWithPaging(ctx *gin.Context, paging *pagingParams, messages []*model.Message) *model.PagedMessages {
	next := ""
	since := uint(0)
	useMessages := messages
	if len(messages) > paging.Limit {
		useMessages = messages[:len(messages)-1]
		since = useMessages[len(useMessages)-1].ID
		url := location.Get(ctx)
		url.Path = ctx.Request.URL.Path
		query := url.Query()
		query.Add("limit", strconv.Itoa(paging.Limit))
		query.Add("since", strconv.FormatUint(uint64(since), 10))
		url.RawQuery = query.Encode()
		next = url.String()
	}
	return &model.PagedMessages{
		Paging:   model.Paging{Size: len(useMessages), Limit: paging.Limit, Next: next, Since: since},
		Messages: toExternalMessages(useMessages),
	}
}

func withPaging(ctx *gin.Context, f func(pagingParams *pagingParams)) {
	params := &pagingParams{Limit: 100}
	if err := ctx.MustBindWith(params, binding.Query); err == nil {
		f(params)
	}
}

func (a *MessageAPI) DeleteMessages(ctx *gin.Context) {
	successOrAbort(ctx, 500, a.DB.DeleteMessages())
}

func (a *MessageAPI) DeleteMessage(ctx *gin.Context) {
	withID(ctx, "id", func(id uint) {
		msg, err := a.DB.GetMessageByID(id)
		if success := successOrAbort(ctx, 500, err); !success {
			return
		}
		if msg == nil {
			ctx.AbortWithError(404, errors.New("message does not exist"))
			return
		}
		successOrAbort(ctx, 500, a.DB.DeleteMessageByID(id))
	})
}

func toInternalMessage(msg *model.MessageExternal) *model.Message {
	res := &model.Message{
		ID:       msg.ID,
		Message:  msg.Message,
		Title:    msg.Title,
		Priority: msg.Priority,
		Date:     msg.Date,
	}
	if msg.Extras != nil {
		res.Extras, _ = json.Marshal(msg.Extras)
	}
	return res
}

func toExternalMessage(msg *model.Message) *model.MessageExternal {
	res := &model.MessageExternal{
		ID:       msg.ID,
		Message:  msg.Message,
		Title:    msg.Title,
		Priority: msg.Priority,
		Date:     msg.Date,
	}
	if len(msg.Extras) != 0 {
		res.Extras = make(map[string]interface{})
		json.Unmarshal(msg.Extras, &res.Extras)
	}
	return res
}

func toExternalMessages(msg []*model.Message) []*model.MessageExternal {
	res := make([]*model.MessageExternal, len(msg))
	for i := range msg {
		res[i] = toExternalMessage(msg[i])
	}
	return res
}
