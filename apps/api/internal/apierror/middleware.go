// Package apierror turns errors that handlers record with c.Error into HTTP responses.
//
// It lives outside internal/app so that handler packages can exercise the real
// middleware in their tests rather than reproducing it.
package apierror

import (
	"net/http"

	"contrib.rocks/apps/api/internal/logger"
	"github.com/gin-gonic/gin"
)

// Middleware responds 500 with the last error a handler recorded via c.Error.
// Handlers that want a different status write the response themselves and return.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		logger.LoggerFromContext(c.Request.Context()).Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.JSON())
	}
}
