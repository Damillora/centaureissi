package httpinterface

import (
	"github.com/Damillora/centaureissi/pkg/services"
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
	chs.g = gin.Default()
	chs.g.Use(cors.Default())
	chs.InitializeRoutes()
}

func (chs *CentaureissiHttpInterface) Start() {
	chs.g.Run()
}
