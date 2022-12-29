package http

import (
	"gateway/config"
	"gateway/http/handlers"
	"gateway/http/middleware"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type IServer interface {
	Run()
}

type server struct {
	cfg    config.Config
	engine *gin.Engine
}

var once sync.Once
var srv IServer

func GetServer() IServer {
	once.Do(func() {

		log.Info().Msg("Initializing server")

		cfg := config.GetConfig()

		e := gin.Default()

		e.Use(middleware.CORS("*"))

		s := &server{
			cfg:    cfg,
			engine: e,
		}

		s.routeUserHandler()

		srv = s
	})

	return srv
}

func (s *server) Run() {
	log.Info().Msg("Running server on port " + s.cfg.Port)

	s.engine.Run(s.cfg.Port)
}

func (s *server) routeUserHandler() {
	h := handlers.GetUserHandler()

	r := s.engine.Group("/users")
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/", h.GetAll)
	r.GET("/:id", h.GetById)

	p := r.Use(middleware.JwtAuth(s.cfg.ApiSecret))
	p.GET("/current", h.Current)
}
