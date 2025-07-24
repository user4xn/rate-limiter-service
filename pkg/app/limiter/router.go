package limiter

import (
	"rate-limiter/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Route(g *gin.RouterGroup) {
	withKey := g.Use(middleware.AuthAPI())

	withKey.POST("/fixed-window", h.FixedWindow)
	withKey.PUT("/fixed-window/set", h.SetClientConfigFixedWindow)
}
