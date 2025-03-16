package httpinterface

import (
	"net/http"

	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/Damillora/centaureissi/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (chs *CentaureissiHttpInterface) InitializeAuthRoutes() {
	chs.g.POST("/api/auth/login", chs.createToken)

	protected := chs.g.Group("/api/auth").Use(chs.AuthMiddleware())
	{
		protected.POST("/token", chs.createTokenLoggedIn)
	}
}
func (chs *CentaureissiHttpInterface) createToken(c *gin.Context) {
	var model models.LoginFormModel
	err := c.ShouldBindJSON(&model)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	validate := validator.New()
	err = validate.Struct(model)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	user := chs.services.Login(model.Username, model.Password)

	if user != nil {
		token := chs.services.CreateToken(user)
		c.JSON(http.StatusOK, models.TokenResponse{
			Token: token,
		})

	} else {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Wrong username or password",
		})
	}
}

func (chs *CentaureissiHttpInterface) createTokenLoggedIn(c *gin.Context) {
	result, ok := c.Get("user")
	if ok {
		user := result.(*schema.User)
		if user != nil {
			token := chs.services.CreateToken(user)
			c.JSON(http.StatusOK, models.TokenResponse{
				Token: token,
			})
		}
	} else {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "No authorized user",
		})
	}
}
