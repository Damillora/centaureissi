package httpinterface

import (
	"net/http"
	"strconv"

	"github.com/Damillora/centaureissi/pkg/database/schema"
	"github.com/Damillora/centaureissi/pkg/models"
	"github.com/gin-gonic/gin"
)

func (chs *CentaureissiHttpInterface) InitializeSearchRoutes() {
	protected := chs.g.Group("/api/search").Use(chs.AuthMiddleware())
	{
		protected.GET("/", chs.searchMail)
	}
}

func (chs *CentaureissiHttpInterface) searchMail(c *gin.Context) {
	result, ok := c.Get("user")
	pageQuery := c.Query("page")
	perPageQuery := c.Query("perPage")
	page, _ := strconv.Atoi(pageQuery)
	perPage, _ := strconv.Atoi(perPageQuery)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	if ok && result != nil {
		user := result.(*schema.User)
		q := c.Query("q")
		result, err := chs.services.SearchMessages(user.ID, q, page, perPage)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Error on searching",
			})
			return
		}
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "User does not exist",
		})
	}
}
