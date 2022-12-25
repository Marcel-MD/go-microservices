package handlers

import (
	"gateway/dto"
	"gateway/infrastructure/services"
	"gateway/presentation/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s server) routeUserHandler() {
	h := &userHandler{
		service: s.userService,
	}

	r := s.router.Group("/users")
	r.POST("/register", h.register)
	r.POST("/login", h.login)
	r.GET("/", h.listAll)
	r.GET("/:id", h.getById)

	p := r.Use(middleware.JwtAuth(s.cfg.ApiSecret))
	p.GET("/current", h.current)
}

type userHandler struct {
	service services.UserService
}

func (h *userHandler) register(c *gin.Context) {

	var dto dto.RegisterUser
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(c, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) login(c *gin.Context) {

	var dto dto.LoginUser
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *userHandler) getById(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.Get(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) listAll(c *gin.Context) {
	users, err := h.service.List(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *userHandler) current(c *gin.Context) {
	id := c.GetString("user_id")

	user, err := h.service.Get(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
