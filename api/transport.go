package api

import (
	"github.com/gin-gonic/gin"
)

// The TransportDatabase interface for encapsulating database access.
type TransportDatabase interface {
	//GetTransport(uuid string, withPoints bool) (*model.Transport, error)
	//GetTransports(withPoints bool) ([]*model.Transport, error)
	CreateTransport(t string, body interface{}) (interface{}, error)
	//UpdateTransport(uuid string, body *model.Transport) (*model.Transport, error)
	//DeleteTransport(uuid string) (bool, error)
	//DropTransports() (bool, error)

}
type TransportAPI struct {
	DB TransportDatabase
}

//func (a *TransportAPI) GetTransports(ctx *gin.Context) {
//	q, err := a.DB.GetTransports(false)
//	reposeHandler(q, err, ctx)
//
//}

//func (a *TransportAPI) GetTransport(ctx *gin.Context) {
//	uuid := resolveID(ctx)
//	q, err := a.DB.GetTransport(uuid, false)
//	reposeHandler(q, err, ctx)
//
//}
//

func (a *TransportAPI) CreateTransport(ctx *gin.Context) {
	body, _ := getBODYTransport(ctx)
	q, err := a.DB.CreateTransport("serial", body)
	reposeHandler(q, err, ctx)
}
