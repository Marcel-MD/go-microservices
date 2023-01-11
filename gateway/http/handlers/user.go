package handlers

import (
	"gateway/models"
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

	SendOtp(c *gin.Context)
	RegisterOtp(c *gin.Context)
}

type userHandler struct {
	userService services.IUserService
	mfaService  services.IMfaService
}

var userOnce sync.Once
var userHnd IUserHandler

func GetUserHandler() IUserHandler {
	userOnce.Do(func() {
		userHnd = &userHandler{
			userService: services.GetUserService(),
			mfaService:  services.GetMfaService(),
		}
	})

	return userHnd
}

func (h *userHandler) Register(c *gin.Context) {

	var model models.RegisterUser
	err := c.ShouldBindJSON(&model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c, model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) Login(c *gin.Context) {

	var model models.LoginUser
	err := c.ShouldBindJSON(&model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.Login(c, model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *userHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.Get(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetAll(c *gin.Context) {
	users, err := h.userService.List(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *userHandler) Current(c *gin.Context) {
	id := c.GetString("user_id")

	user, err := h.userService.Get(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "record not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) SendOtp(c *gin.Context) {
	email := c.Param("email")

	_, err := h.mfaService.GenerateOtp(c, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "otp sent"})
}

func (h *userHandler) RegisterOtp(c *gin.Context) {

	var model models.RegisterOtpUser
	err := c.ShouldBindJSON(&model)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isValid, err := h.mfaService.VerifyOtp(c, model.Email, model.Otp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid otp"})
		return
	}

	user, err := h.userService.Register(c, model.RegisterUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
