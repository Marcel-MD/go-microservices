package handlers

import (
	"gateway/dto"
	"gateway/services"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	GetAll(c *gin.Context)
	GetById(c *gin.Context)
	Current(c *gin.Context)
}

type userHandler struct {
	service services.IUserService
}

var userOnce sync.Once
var userHnd IUserHandler

func GetUserHandler() IUserHandler {
	userOnce.Do(func() {
		userHnd = &userHandler{
			service: services.GetUserService(),
		}
	})

	return userHnd
}

func (h *userHandler) Register(c *gin.Context) {

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

func (h *userHandler) Login(c *gin.Context) {

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

func (h *userHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.Get(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetAll(c *gin.Context) {
	users, err := h.service.List(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *userHandler) Current(c *gin.Context) {
	id := c.GetString("user_id")

	user, err := h.service.Get(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
