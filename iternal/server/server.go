// Copyright 2022 eightlay (github.com/eightlay). All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package server

import (
	"github.com/gin-gonic/gin"
)

// Start game server
func StartServer(addr string) {
	m := newManager()

	r := gin.New()
	r.GET("/ws", func(m *Manager) gin.HandlerFunc {
		return gin.HandlerFunc(func(c *gin.Context) {
			serveWs(m, c.Writer, c.Request)
		})
	}(m))

	r.Run(addr)
}
