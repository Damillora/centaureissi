package httpinterface

import (
	"io/fs"
	"net/http"

	"github.com/Damillora/centaureissi/pkg/services"
	"github.com/Damillora/centaureissi/pkg/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CentaureissiHttpInterface struct {
	g        *gin.Engine
	services *services.CentaureissiService
}

func New(cs *services.CentaureissiService) *CentaureissiHttpInterface {
	httpServer := &CentaureissiHttpInterface{
		services: cs,
	}
	httpServer.Initialize()
	return httpServer
}

func (chs *CentaureissiHttpInterface) Initialize() {
	webFS := web.WebAssets()
	webAssets, _ := fs.Sub(webFS, "_app")

	chs.g = gin.Default()

	chs.g.NoRoute(func(c *gin.Context) {
		c.FileFromFS("./app.html", http.FS(webFS))
	})
	chs.g.StaticFS("/_app", http.FS(webAssets))

	chs.g.Use(cors.Default())
	chs.InitializeRoutes()
}

func (chs *CentaureissiHttpInterface) Start() {
	chs.g.Run()
}
