package handlers

import (
	"gateway/domain"
	"gateway/infrastructure/services"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run() error
}

type server struct {
	cfg         domain.Config
	router      *gin.Engine
	userService services.UserService
}

func NewServer(cfg domain.Config, userService services.UserService) Server {

	r := gin.Default()

	s := server{
		cfg:         cfg,
		router:      r,
		userService: userService,
	}

	s.routeUserHandler()

	return &s
}

func (s *server) Run() error {
	return s.router.Run(s.cfg.Port)
}
