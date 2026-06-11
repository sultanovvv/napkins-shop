package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"napkins-shop/public-executor/api"
	"shared/entity"
)

func loginResponseFromEntity(user *entity.User, pair *entity.TokenPair) api.AuthLoginResponse {
	return api.AuthLoginResponse{
		User:             authUserResponse(user),
		AccessToken:      pair.AccessToken,
		AccessExpiresAt:  pair.AccessExpiresAt,
		RefreshExpiresAt: pair.RefreshExpiresAt,
	}
}

func authUserResponse(user *entity.User) api.AuthUser {
	return api.AuthUser{
		Id:              user.PublicID,
		Email:           user.Email,
		EmailVerifiedAt: user.EmailVerifiedAt,
	}
}

func badRequest(ctx echo.Context, msg string) error {
	return errorJSON(ctx, http.StatusBadRequest, "BAD_REQUEST", msg)
}

func internalErr(ctx echo.Context) error {
	return errorJSON(ctx, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}

func unauthorizedJSON(ctx echo.Context, msg string) error {
	return errorJSON(ctx, http.StatusUnauthorized, "UNAUTHORIZED", msg)
}

func errorJSON(ctx echo.Context, status int, code, msg string) error {
	return ctx.JSON(status, api.ErrorModel{
		Error: api.ErrorDetailsModel{
			Code:    code,
			Message: msg,
		},
	})
}