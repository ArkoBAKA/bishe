package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func RespondOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, APIResponse{Code: 0, Message: "ok", Data: data})
}

func RespondFail(c *gin.Context, httpStatus, code int, msg string) {
	c.JSON(httpStatus, APIResponse{Code: code, Message: msg, Data: map[string]any{}})
}
