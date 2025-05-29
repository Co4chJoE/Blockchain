// middleware/recovery.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"`
}

func Recovery(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic: %v\n", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
					Error: "Internal Server Error",
					Code:  http.StatusInternalServerError,
				})
			}
		}()
		c.Next()
	}
}
