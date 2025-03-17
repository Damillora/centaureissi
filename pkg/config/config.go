package config

import (
	"os"
	"strings"
)

type Config struct {
	AuthSecret          string
	DisableRegistration bool
	DataDirectory       string
	BlobDirectory       string
	ImapTlsCertFile     string
	ImapTlsKeyFile      string
	InsecureAuth        bool
	ListenImap          string
	Debug               bool
}

var CurrentConfig Config

func (config *Config) InitializeConfig() {
	dataDir := os.Getenv("DATA_DIR")
	if len(dataDir) == 0 {
		dataDir = "/data"
	}
	var blobDirBuilder strings.Builder
	blobDirBuilder.WriteString(dataDir)
	blobDirBuilder.WriteString("/blobs")

	blobDir := os.Getenv("BLOB_DIR")
	if len(blobDir) == 0 {
		blobDir = blobDirBuilder.String()
	}

	listenImap := os.Getenv("IMAP_LISTEN")
	if len(listenImap) == 0 {
		listenImap = "localhost:143"
	}

	config.DataDirectory = dataDir
	config.BlobDirectory = blobDir
	config.AuthSecret = os.Getenv("AUTH_SECRET")
	config.DisableRegistration = strings.ToLower(os.Getenv("DISABLE_REGISTRATION")) == "true"
	config.ImapTlsCertFile = os.Getenv("IMAP_CERT_FILE")
	config.ImapTlsKeyFile = os.Getenv("IMAP_KEY_FILE")
	config.InsecureAuth = strings.ToLower(os.Getenv("INSECURE_AUTH")) == "true"
	config.Debug = strings.ToLower(os.Getenv("DEBUG")) == "true"
	config.ListenImap = listenImap
}
