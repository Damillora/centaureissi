package httpinterface

import (
	"net/http"

	"github.com/Damillora/centaureissi/pkg/config"
	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/Damillora/centaureissi/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (chs *CentaureissiHttpInterface) InitializeUserRoutes() {
	chs.g.POST("/api/user/register", chs.registerUser)

	protected := chs.g.Group("/api/user").Use(chs.AuthMiddleware())
	{
		protected.GET("/profile", chs.userProfile)
		protected.PUT("/update", chs.userUpdate)
		protected.PUT("/update-password", chs.userUpdatePassword)
	}
}

func (chs *CentaureissiHttpInterface) registerUser(c *gin.Context) {
	var disableRegistration = config.CurrentConfig.DisableRegistration
	if disableRegistration {
		c.JSON(http.StatusForbidden, models.ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "registration is disabled",
		})
		c.Abort()
		return
	}

	var model models.UserCreateModel
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

	user, err := chs.services.CreateUser(model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, models.UserProfileResponse{
		Username: user.Username,
	})
}

func (chs *CentaureissiHttpInterface) userProfile(c *gin.Context) {
	result, ok := c.Get("user")
	if ok {
		user := result.(*schema.User)
		c.JSON(http.StatusOK, models.UserProfileResponse{
			Username: user.Username,
		})
	} else {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User does not exist",
		})
	}
}

func (chs *CentaureissiHttpInterface) userUpdate(c *gin.Context) {
	var model models.UserUpdateModel

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

	result, ok := c.Get("user")
	if ok {
		user := result.(*schema.User)
		chs.services.UpdateUserProfile(user.ID, model)
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User does not exist",
		})
	}
}

func (chs *CentaureissiHttpInterface) userUpdatePassword(c *gin.Context) {
	var model models.UserUpdatePasswordModel

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

	result, ok := c.Get("user")
	if ok {
		user := result.(*schema.User)
		chs.services.UpdateUserPassword(user.ID, model)
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User does not exist",
		})
	}
}
