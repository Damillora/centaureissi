package httpinterface

import (
	"strings"

	"github.com/Damillora/centaureissi/pkg/models"
	"github.com/gin-gonic/gin"
)

func (chs *CentaureissiHttpInterface) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(403, models.ErrorResponse{
				Code:    403,
				Message: "Authorization required",
			})
			c.Abort()
			return
		}

		extractedToken := strings.Split(clientToken, "Bearer ")

		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			c.JSON(400, models.ErrorResponse{
				Code:    400,
				Message: "Incorrect Format of Authorization Token",
			})
			c.Abort()
			return
		}

		claims, err := chs.services.ValidateToken(clientToken)
		if err != nil {
			c.JSON(401, models.ErrorResponse{
				Code:    401,
				Message: err.Error(),
			})
			c.Abort()
			return
		}
		user, err := chs.services.GetUserById(claims["sub"].(string))
		if err != nil {
			c.JSON(500, models.ErrorResponse{
				Code:    500,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
