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
		case
			errors.Is(err, errs.InvalidToken),
			errors.Is(err, errs.IncorrectCredentials),
			errors.Is(err, errs.Unauthorized):
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		case errors.Is(err, errs.Forbidden):
			c.AbortWithStatusJSON(http.StatusForbidden, err.Error())
			return
		case errors.Is(err, errs.ResourceAlreadyExists):
			c.AbortWithStatusJSON(http.StatusConflict, err.Error())
			return
		case errors.Is(err, errs.InternalServerError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		case
			errors.Is(err, errs.BadRequest),
			errors.Is(err, errs.ResourceNotFound),
			errors.Is(err, errs.DbError):
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
	}
}
