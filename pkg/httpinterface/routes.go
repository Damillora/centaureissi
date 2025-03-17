package httpinterface

func (chs *CentaureissiHttpInterface) InitializeRoutes() {
	chs.InitializeAuthRoutes()
	chs.InitializeUserRoutes()
	chs.InitializeSearchRoutes()
}
