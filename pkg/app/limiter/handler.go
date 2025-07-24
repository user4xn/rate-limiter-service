package limiter

import (
	"net/http"
	"rate-limiter/pkg/dto"
	"rate-limiter/pkg/factory"
	"rate-limiter/pkg/util"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) FixedWindow(c *gin.Context) {
	payload := dto.PayloadLimiter{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		util.CreateErrorLog(err)
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "bad request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	val, exists := c.Get("clientID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Client ID not found"})
		return
	}

	payload.ClientID = val.(string)

	data, err := h.service.FixedWindow(c, payload)
	if err != nil {
		util.CreateErrorLog(err)
		response := util.APIResponse(err.Error(), http.StatusInternalServerError, "internal server error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("success", http.StatusOK, "ok", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) SetClientConfigFixedWindow(c *gin.Context) {
	payload := dto.PayloadConfigClient{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		util.CreateErrorLog(err)
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "bad request", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	val, exists := c.Get("clientID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Client ID not found"})
		return
	}

	payload.ClientID = val.(string)

	err := h.service.SetClientConfigFixedWindow(c, payload)
	if err != nil {
		util.CreateErrorLog(err)
		response := util.APIResponse(err.Error(), http.StatusInternalServerError, "internal server error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("success", http.StatusOK, "ok", nil)
	c.JSON(http.StatusOK, response)
}
