package router

import (
	"authorization/controllers"
	"authorization/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func New(
	r *gin.RouterGroup,
	controller controllers.IController,
) {
	g := r.Group("/authorization")

	// register route
	g.POST("/v1/register", func(c *gin.Context) {
		log.Printf("hit service register with request %v", c.Request)

		var parsedBody DTORegisterRoute
		if err := c.ShouldBindJSON(&parsedBody); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request body",
				"code": http.StatusBadRequest,
				"error": err.Error(),
			})
			return
		}

		log.Println("running controller")
		result, err := controller.RegisterUser(c.Request.Context(), &controllers.RequestRegisterService{
			Name: parsedBody.Name,
			Email: parsedBody.Email,
			Password: parsedBody.Password,
		})

		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				if pgErr.Code == "23505" {
					log.Printf("user already exist: %v", err)
					c.JSON(http.StatusConflict, gin.H{
						"message": "invalid request body",
						"code": http.StatusConflict,
						"error": err.Error(),
					})
					return
				}
			}

			log.Printf("failed on controller process: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
				"code": http.StatusInternalServerError,
				"error": err.Error(),
			})
			return
		}

		log.Println("success register")
		c.JSON(http.StatusCreated, gin.H{
			"message": "success register",
			"code": http.StatusCreated,
			"result": result,
		})
	})

	// login route
	g.POST("/v1/login", func(c *gin.Context) {
		log.Printf("hit service login with request %v", c.Request)

		var parsedBody DTOLoginRoute
		if err := c.ShouldBindJSON(&parsedBody); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request body",
				"code": http.StatusBadRequest,
				"error": err.Error(),
			})
			return
		}

		if parsedBody.Name == "" && parsedBody.Email == "" {
			log.Printf("missing request")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request body",
				"code": http.StatusBadRequest,
				"error": "missing name and email, at least email or name",
			})
			return
		}

		log.Println("running controller")
		result, err := controller.LoginUser(c.Request.Context(), &controllers.RequestLoginService{
			Name: parsedBody.Name,
			Email: parsedBody.Email,
			Password: parsedBody.Password,
		})

		if err != nil {
			// if pgErr, ok := err.(*pgconn.PgError); ok {
			// 	if pgErr.Code == "23505" {
			// 		log.Printf("user already exist: %v", err)
			// 		c.JSON(http.StatusConflict, gin.H{
			// 			"message": "invalid request body",
			// 			"code": http.StatusConflict,
			// 			"error": err.Error(),
			// 		})
			// 		return
			// 	}
			// }

			log.Printf("failed on controller process: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
				"code": http.StatusInternalServerError,
				"error": err.Error(),
			})
			return
		}

		log.Println("success login")
		c.JSON(http.StatusCreated, gin.H{
			"message": "success login",
			"code": http.StatusCreated,
			"result": result,
		})
	})

	// refresh token route
	g.POST("/v1/refresh-token", func(c *gin.Context) {
		log.Printf("hit service refresh-token with request %v", c.Request)

		var parsedBody DTORefreshTokenRoute
		if err := c.ShouldBindJSON(&parsedBody); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request body",
				"code": http.StatusBadRequest,
				"error": err.Error(),
			})
			return
		}

		log.Println("running controller")
		result, err := controller.RefreshToken(c.Request.Context(), &controllers.RequestRefreshTokenService{
			RefreshToken: parsedBody.RefreshToken,
		})

		if err != nil {
			log.Printf("failed on controller process: %v", err)
			if err.Error() == "invalid token" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "invalid token",
					"code": http.StatusUnauthorized,
					"error": err.Error(),
				})
				return		
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
				"code": http.StatusInternalServerError,
				"error": err.Error(),
			})
			return
		}

		log.Println("success refresh-token")
		c.JSON(http.StatusCreated, gin.H{
			"message": "success refresh-token",
			"code": http.StatusCreated,
			"result": result,
		})
	})

	// logout route
	g.POST("/v1/logout", utils.AuthMiddleware(), func(c *gin.Context) {
		log.Printf("hit service logout with request %v", c.Request)

		userID, exist := c.Get("user_id")
		if !exist {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
				"code": http.StatusInternalServerError,
				"error": "missing user id",
			})
			return
		}

		log.Println("running controller")
		result, err := controller.LogoutUser(c.Request.Context(), &controllers.RequestLogoutService{
			UserID: userID.(string),
		})

		if err != nil {
			log.Printf("failed on controller process: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal server error",
				"code": http.StatusInternalServerError,
				"error": err.Error(),
			})
			return
		}

		log.Println("success refresh-token")
		c.JSON(http.StatusCreated, gin.H{
			"message": "success refresh-token",
			"code": http.StatusCreated,
			"result": result,
		})
	})
}