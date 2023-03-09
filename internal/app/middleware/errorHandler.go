package middleware

import (
	"errors"
	"net/http"

	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		switch {
		case errors.Is(err, errs.InvalidToken):
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
	}
}
