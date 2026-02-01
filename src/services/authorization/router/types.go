package router

// types register route
type (
	DTORegisterRoute struct {
		Name string `json:"name" binding:"required,min=3"`
		Email string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
)

// types login route
type (
	DTOLoginRoute struct {
		Name string `json:"name" binding:"min=3"`
		Email string `json:"email" binding:"email"`
		Password string `json:"password" binding:"required,min=8"`
	}
)

// type refresh token route
type (
	DTORefreshTokenRoute struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
)