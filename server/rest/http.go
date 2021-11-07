package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/arxeiss/sample-terraform-provider/server/database"
)

type HTTPServer struct {
	log    *logrus.Entry
	db     *database.DB
	router *gin.Engine
}

func (h *HTTPServer) Run(addr string) error {
	h.log.Info("Listening on ", addr)
	return h.router.Run(addr)
}

func NewHTTPServer(log *logrus.Entry, db *database.DB, token string) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	h := &HTTPServer{
		log:    log,
		db:     db,
		router: r,
	}
	r.Use(authorization(token))
	r.Use(gin.Recovery())
	r.Use(h.loggerMiddleware)

	vmg := r.Group("vm")
	{
		vmg.GET("/:id", h.getVM)
		vmg.PUT("", h.createVM)
		vmg.POST("/:id", h.editVM)
		vmg.DELETE("/:id", h.deleteVM)
	}
	ng := r.Group("network")
	{
		ng.GET("/:id", h.getNetwork)
		ng.PUT("", h.createNetwork)
		ng.POST("/:id", h.editNetwork)
		ng.DELETE("/:id", h.deleteNetwork)
	}
	sg := r.Group("storage")
	{
		sg.GET("/:id", h.getStorage)
		sg.PUT("", h.createStorage)
		sg.POST("/:id", h.editStorage)
		sg.DELETE("/:id", h.deleteStorage)
	}

	r.NoRoute(h.error404)

	return h
}
