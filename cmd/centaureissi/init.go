package main

import (
	"os"

	"github.com/Damillora/centaureissi/pkg/blob"
	"github.com/Damillora/centaureissi/pkg/config"
	"github.com/Damillora/centaureissi/pkg/database"
	"github.com/Damillora/centaureissi/pkg/httpinterface"
	"github.com/Damillora/centaureissi/pkg/imapinterface"
	"github.com/Damillora/centaureissi/pkg/search"
	"github.com/Damillora/centaureissi/pkg/services"
)

var Repository *database.CentaureissiRepository
var Blobs *blob.CentaureissiBlobRepository
var Service *services.CentaureissiService
var HttpInterface *httpinterface.CentaureissiHttpInterface
var ImapInterface *imapinterface.CentaureissiImapServer
var SearchEngine *search.CentaureissiSearchEngine

func Initialize() {
	config.CurrentConfig.InitializeConfig()

	dataDir := config.CurrentConfig.DataDirectory

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, 0755)
	}

	Repository = database.New()
	Blobs = blob.New()
	SearchEngine = search.New()
	Service = services.New(Repository, Blobs, SearchEngine)
	HttpInterface = httpinterface.New(Service)
	ImapInterface = imapinterface.New(Service)
}

func Deinitialize() {
	defer Repository.Deinitialize()
}

func Start() {
	go ImapInterface.Start()
	HttpInterface.Start()

	defer Deinitialize()
}
