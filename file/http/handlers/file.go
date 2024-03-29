package handlers

import (
	"file/services"
	"mime"
	"net/http"
	"path"
	"sync"

	"github.com/gin-gonic/gin"
)

type IFileHandler interface {
	GetAll(c *gin.Context)
	GetByOwnerId(c *gin.Context)
	GetByName(c *gin.Context)
	Read(c *gin.Context)
	Upload(c *gin.Context)
	Delete(c *gin.Context)
}

type fileHandler struct {
	fileService services.IFileService
}

var (
	fileOnce sync.Once
	fileHnd  IFileHandler
)

func GetFileHandler() IFileHandler {
	fileOnce.Do(func() {
		fileHnd = &fileHandler{
			fileService: services.GetFileService(),
		}
	})

	return fileHnd
}

func (h *fileHandler) GetAll(c *gin.Context) {
	files, err := h.fileService.FindAll(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

func (h *fileHandler) GetByOwnerId(c *gin.Context) {
	ownerId := c.Param("owner-id")
	files, err := h.fileService.FindByOwnerId(c, ownerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, files)
}

func (h *fileHandler) GetByName(c *gin.Context) {
	fileName := c.Param("file-name")
	file, err := h.fileService.FindByName(c, fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, file)
}

func (h *fileHandler) Read(c *gin.Context) {
	fileName := c.Param("file-name")
	reader, err := h.fileService.Read(c, fileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mimeType := mime.TypeByExtension(path.Ext(fileName))
	c.DataFromReader(http.StatusOK, -1, mimeType, reader, nil)
}

func (h *fileHandler) Upload(c *gin.Context) {
	userId := c.GetString("user_id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	name := header.Filename
	extension := path.Ext(name)

	newFile, err := h.fileService.Upload(c, file, extension, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newFile)
}

func (h *fileHandler) Delete(c *gin.Context) {
	userId := c.GetString("user_id")
	fileName := c.Param("file-name")

	err := h.fileService.Delete(c, fileName, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
