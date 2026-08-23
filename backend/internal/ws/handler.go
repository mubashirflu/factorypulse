package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	GlobalHub.Register(conn)
	defer GlobalHub.Unregister(conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
