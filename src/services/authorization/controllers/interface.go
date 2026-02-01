package controllers

import "context"

type (
	IController interface {
		RegisterUser(ctx context.Context, request *RequestRegisterService) (response *ResponseRegisterService, err error)
		LoginUser(ctx context.Context, request *RequestLoginService) (response *ResponseLoginService, err error)
		LogoutUser(ctx context.Context, request *RequestLogoutService) (response *ResponseLogoutService, err error)
		RefreshToken(ctx context.Context, request *RequestRefreshTokenService) (response *ResponseRefreshTokenService, err error)
	}
)