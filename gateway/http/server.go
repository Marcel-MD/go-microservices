package http

import (
	"gateway/config"
	"gateway/http/handlers"
	"gateway/http/middleware"
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
	h := handlers.GetUserHandler()

	r := router.Group("/users")
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/", h.GetAll)
	r.GET("/:id", h.GetById)

	r.POST("/otp/:email", h.SendOtp)
	r.POST("/register-otp", h.RegisterOtp)

	p := r.Use(middleware.JwtAuth(cfg.ApiSecret))
	p.GET("/current", h.Current)
}
