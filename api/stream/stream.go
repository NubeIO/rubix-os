package stream

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/NubeIO/flow-framework/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// The API provides a handler for a WebSocket stream API.
type API struct {
	clients     []*client
	lock        sync.RWMutex
	pingPeriod  time.Duration
	pongTimeout time.Duration
	upgrader    *websocket.Upgrader
}

// New creates a new instance of API.
// pingPeriod: is the interval, in which is server sends the a ping to the client.
// pongTimeout: is the duration after the connection will be terminated, when the client does not respond with the
// pong command.
func New(pingPeriod, pongTimeout time.Duration, allowedWebSocketOrigins []string, prod bool) *API {
	return &API{
		clients:     []*client{},
		pingPeriod:  pingPeriod,
		pongTimeout: pingPeriod + pongTimeout,
		upgrader:    newUpgrader(allowedWebSocketOrigins, prod),
	}
}

// NotifyDeletedUser closes existing connections for the given user.
func (a *API) NotifyDeletedUser(userID uint) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, client := range a.clients {
		client.Close()
	}
	return nil
}

// Notify notifies the clients with the given userID that a new messages was created.
func (a *API) Notify(msg *model.MessageExternal) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	for _, c := range a.clients {
		c.write <- msg
	}
}

func (a *API) remove(remove *client) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for i, client := range a.clients {
		if client == remove {
			a.clients = append(a.clients[:i], a.clients[i+1:]...)
			break
		}
	}
}

func (a *API) register(client *client) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.clients = append(a.clients, client)
}

// Handle handles incoming requests. First it upgrades the protocol to the WebSocket protocol and then starts listening
func (a *API) Handle(ctx *gin.Context) {
	conn, err := a.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Error(err)
		return
	}

	client := newClient(conn, a.remove)
	a.register(client)
	go client.startReading(a.pongTimeout)
	go client.startWriteHandler(a.pingPeriod)
}

// Close closes all client connections and stops answering new connections.
func (a *API) Close() {
	a.lock.Lock()
	defer a.lock.Unlock()
	for _, client := range a.clients {
		client.Close()
	}
	for _, client := range a.clients {
		a.remove(client)
	}
}

func isAllowedOrigin(r *http.Request, allowedOrigins []*regexp.Regexp) bool {
	origin := r.Header.Get("origin")
	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if strings.EqualFold(u.Host, r.Host) {
		return true
	}

	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin.Match([]byte(strings.ToLower(u.Hostname()))) {
			return true
		}
	}

	return false
}

func newUpgrader(allowedWebSocketOrigins []string, prod bool) *websocket.Upgrader {
	compiledAllowedOrigins := compileAllowedWebSocketOrigins(allowedWebSocketOrigins)
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if !prod {
				return true
			}
			return isAllowedOrigin(r, compiledAllowedOrigins)
		},
	}
}

func compileAllowedWebSocketOrigins(allowedOrigins []string) []*regexp.Regexp {
	var compiledAllowedOrigins []*regexp.Regexp
	for _, origin := range allowedOrigins {
		compiledAllowedOrigins = append(compiledAllowedOrigins, regexp.MustCompile(origin))
	}

	return compiledAllowedOrigins
}
