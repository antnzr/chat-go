package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/antnzr/chat-go/config"
	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/gin-gonic/gin"
)

func Auth(
	tokenService domain.TokenService,
	userService domain.UserService,
	config config.Config,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("accessToken")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authCtx := context.Background()
		tokenDetails, err := tokenService.ValidateToken(authCtx, accessToken, config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		user, err := userService.GetMe(authCtx, tokenDetails.UserId)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, errs.ResourceNotFound)
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
