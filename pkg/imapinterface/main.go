package imapinterface

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"

	"github.com/Damillora/centaureissi/pkg/config"
	"github.com/Damillora/centaureissi/pkg/services"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	_ "github.com/emersion/go-message/charset"
)

type CentaureissiImapServer struct {
	services *services.CentaureissiService
	tracker  *CentaureissiMailboxTracker

	imapServer *imapserver.Server
}

func New(s *services.CentaureissiService) *CentaureissiImapServer {
	imapServer := &CentaureissiImapServer{
		services: s,
		tracker:  NewCentaureissiMailboxTracker(s),
	}
	imapServer.Initialize()
	return imapServer
}

func (cis *CentaureissiImapServer) Initialize() {
	tlsCert := config.CurrentConfig.ImapTlsCertFile
	tlsKey := config.CurrentConfig.ImapTlsKeyFile
	insecureAuth := config.CurrentConfig.InsecureAuth

	var tlsConfig *tls.Config
	if tlsCert != "" || tlsKey != "" {
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			log.Fatalf("Failed to load TLS key pair: %v", err)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	var debugWriter io.Writer
	if config.CurrentConfig.Debug {
		debugWriter = os.Stdout
	}

	cis.imapServer = imapserver.New(&imapserver.Options{
		NewSession: func(conn *imapserver.Conn) (imapserver.Session, *imapserver.GreetingData, error) {
			return NewCentaureissiImapSession(cis.services, cis.tracker), nil, nil
		},
		Caps: imap.CapSet{
			imap.CapIMAP4rev1: {},
			imap.CapIMAP4rev2: {},
		},
		TLSConfig:    tlsConfig,
		InsecureAuth: insecureAuth,
		DebugWriter:  debugWriter,
	})
}

func (cis *CentaureissiImapServer) Start() {
	listen := config.CurrentConfig.ListenImap
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("IMAP server listening on %v", ln.Addr())

	if err := cis.imapServer.Serve(ln); err != nil {
		log.Fatalf("Serve() = %v", err)
	}
}
