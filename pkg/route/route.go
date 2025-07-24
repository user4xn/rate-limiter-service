package route

import (
	"rate-limiter/pkg/app/limiter"
	"rate-limiter/pkg/factory"

	"github.com/gin-gonic/gin"
)

func NewAPIHttp(g *gin.Engine, f *factory.Factory) {
	Index(g)

	// logs incoming requests, recover panic and return 500 status
	g.Use(gin.Logger(), gin.Recovery())

	v1 := g.Group("/api/v1")

	limiter.NewHandler(f).Route(v1.Group("/rate"))
}

func Index(g *gin.Engine) {
	g.GET("/", func(context *gin.Context) {
		context.JSON(200, struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{
			Name:    "rate-limiter-service",
			Version: "1.0.0",
		})
	})
}
