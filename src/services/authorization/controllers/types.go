package controllers

import "time"

// general
type (
	UserTypeResponse struct {
		ID string
		Name string
		Email string
		CreatedAT time.Time
	}

	UserType struct {
		ID string
		Name string
		Email string
		Password string
		CreatedAT time.Time
	}
)

// register types
type (
	RequestRegisterService struct {
		Name string
		Email string
		Password string
	}

	ResponseRegisterService struct {
		User UserTypeResponse `json:"user"`
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

// login types
type (
	RequestLoginService struct {
		Name string
		Email string
		Password string
	}

	ResponseLoginService struct {
		User UserTypeResponse `json:"user"`
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

// logout types
type (
	RequestLogoutService struct {
		UserID string
	}

	ResponseLogoutService struct {}
)

// refresh token types
type (
	RequestRefreshTokenService struct {
		RefreshToken string
	}

	ResponseRefreshTokenService struct {
		AccessToken string `json:"access_token"`
	}
)