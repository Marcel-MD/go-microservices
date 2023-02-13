package http

import (
	"file/config"
	"file/http/handlers"
	"file/http/middleware"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
	once sync.Once
	srv  *http.Server
)

func GetServer() *http.Server {
	once.Do(func() {

		log.Info().Msg("Initializing server")

		cfg := config.GetConfig()

		e := gin.Default()

		e.Use(middleware.CORS(cfg.AllowOrigin))

		r := e.Group("/api")

		routeUserHandler(r, cfg)

		s := &http.Server{
			Addr:    cfg.Port,
			Handler: e,
		}

		srv = s
	})

	return srv
}

func routeUserHandler(router *gin.RouterGroup, cfg config.Config) {
	h := handlers.GetFileHandler()

	r := router.Group("/files")
	r.GET("/", h.GetAll)
	r.GET("/owner/:id", h.GetByOwnerId)
	r.GET("/:name", h.GetByName)
	r.GET("/read/:name", h.Read)

	p := r.Use(middleware.JwtAuth(cfg.ApiSecret))
	p.POST("/", h.Upload)
	p.DELETE("/:name", h.Delete)
}
