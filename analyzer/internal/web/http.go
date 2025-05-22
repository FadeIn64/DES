package web

import (
	"DAS/internal/controllers"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type httpServer struct {
	Router *gin.Engine
	srv    http.Server
}

func New(port string, metricHandler http.Handler, ctrl ...controllers.Controller) (HttpServer, error) {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowCredentials = true
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowHeaders = append(corsCfg.AllowHeaders, "*")
	r.Use(cors.New(corsCfg))

	for i := 0; i < len(ctrl); i++ {
		ctrl[i].Register(r)
	}
	r.GET("/metrics", func(c *gin.Context) {
		metricHandler.ServeHTTP(c.Writer, c.Request)
	})

	return &httpServer{
		Router: r,
		srv: http.Server{
			Addr:    ":" + port,
			Handler: r,
		},
	}, nil
}

type HttpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func (h *httpServer) ListenAndServe() error {
	return h.srv.ListenAndServe()
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	return h.srv.Shutdown(ctx)
}
