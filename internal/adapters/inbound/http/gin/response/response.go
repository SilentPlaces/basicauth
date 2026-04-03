package response

import (
	generaldto "github.com/SilentPlaces/basicauth.git/internal/dto/general"
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, generaldto.Response{
		Status:  "success",
		Code:    statusCode,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, generaldto.Response{
		Status:  "error",
		Code:    statusCode,
		Message: message,
	})
}
