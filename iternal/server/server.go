package server

import (
	"github.com/gin-gonic/gin"
)

// Start game server
func StartServer() {
	m := newManager()

	r := gin.New()
	r.GET("/ws", func(m *Manager) gin.HandlerFunc {
		return gin.HandlerFunc(func(c *gin.Context) {
			serveWs(m, c.Writer, c.Request)
		})
	}(m))
	r.Run(":8080")
}
