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
			_ = ctx.AbortWithError(http.StatusUnauthorized, errs.Unauthorized)
			return
		}

		authCtx := context.Background()
		tokenDetails, err := tokenService.ValidateToken(authCtx, accessToken, config.AccessTokenPublicKey)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errs.Unauthorized)
			return
		}

		user, err := userService.FindById(authCtx, tokenDetails.UserId)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errs.Unauthorized)
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
