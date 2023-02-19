package settings

import (
	"errors"
	"os"
	"strings"

	ce "github.com/vault-thirteen/SFRODB/pkg/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/common/helper"
	"github.com/vault-thirteen/auxie/reader"
	"github.com/vault-thirteen/errorz"
)

const (
	ErrServerModeIsNotSet     = "server mode is not set"
	ErrServerMode             = "server mode error"
	ErrCertFileIsNotSet       = "certificate file is not set"
	ErrKeyFileIsNotSet        = "key file is not set"
	ErrFileExtensionIsNotSet  = "file extension is not set"
	ErrMimeTypeIsNotSet       = "MIME type is not set"
	ErrDbClientPoolSize       = "DB client pool size is not set"
	ErrHttpCacheControlMaxAge = "HTTP cache control max-age error"
)

const (
	ContentDispositionInline = "inline"

	ServerModeHttp    = "HTTP"
	ServerModeIdHttp  = 1
	ServerModeHttps   = "HTTPS"
	ServerModeIdHttps = 2
)

// Settings is Server's settings.
type Settings struct {
	// Path to the File with these Settings.
	File string

	// Server's Host Name.
	ServerHost string

	// Server's Listen Port.
	ServerPort uint16

	// ServerMode is an HTTP mode selector.
	// Possible values are: HTTP and HTTPS.
	ServerModeStr string
	ServerModeId  byte

	// Server's Certificate and Key.
	CertFile string
	KeyFile  string

	// Database Host Name.
	DbHost string

	// Database Ports.
	DbPortA uint16
	DbPortB uint16

	// DbClientPoolSize is the size of a pool of DB clients.
	DbClientPoolSize int

	// File Extension & MIME Type.
	// Extension which is appended to all files served.
	FileExtension string
	MimeType      string

	// HttpCacheControlMaxAge is time in seconds for which this server's
	// response is fresh (valid). After this period clients will be refreshing
	// the stale content by re-requesting it from the server.
	HttpCacheControlMaxAge uint

	// Allowed Origin for cross-origin requests (CORS).
	AllowedOriginForCORS string
}

func NewSettingsFromFile(filePath string) (stn *Settings, err error) {
	stn = &Settings{
		File: filePath,
	}

	var file *os.File
	file, err = os.Open(stn.File)
	if err != nil {
		return stn, err
	}
	defer func() {
		derr := file.Close()
		if derr != nil {
			err = errorz.Combine(err, derr)
		}
	}()

	rdr := reader.NewReader(file)
	var buf = make([][]byte, 13)

	for i := range buf {
		buf[i], err = rdr.ReadLineEndingWithCRLF()
		if err != nil {
			return stn, err
		}
	}

	// Server Host & Port.
	stn.ServerHost = strings.TrimSpace(string(buf[0]))

	stn.ServerPort, err = helper.ParseUint16(strings.TrimSpace(string(buf[1])))
	if err != nil {
		return stn, err
	}

	// Server Work Mode.
	stn.ServerModeStr = strings.ToUpper(strings.TrimSpace(string(buf[2])))
	switch stn.ServerModeStr {
	case ServerModeHttp:
		stn.ServerModeId = ServerModeIdHttp
	case ServerModeHttps:
		stn.ServerModeId = ServerModeIdHttps
	}

	// Certificate and Key for optional TLS.
	stn.CertFile = strings.TrimSpace(string(buf[3]))
	stn.KeyFile = strings.TrimSpace(string(buf[4]))

	// Database.
	stn.DbHost = strings.TrimSpace(string(buf[5]))

	stn.DbPortA, err = helper.ParseUint16(strings.TrimSpace(string(buf[6])))
	if err != nil {
		return stn, err
	}

	stn.DbPortB, err = helper.ParseUint16(strings.TrimSpace(string(buf[7])))
	if err != nil {
		return stn, err
	}

	stn.DbClientPoolSize, err = helper.ParseInt(strings.TrimSpace(string(buf[8])))
	if err != nil {
		return stn, err
	}

	// HTTP.
	stn.FileExtension = strings.TrimSpace(string(buf[9]))
	stn.MimeType = strings.TrimSpace(string(buf[10]))

	stn.HttpCacheControlMaxAge, err = helper.ParseUint(strings.TrimSpace(string(buf[11])))
	if err != nil {
		return stn, err
	}

	stn.AllowedOriginForCORS = strings.TrimSpace(string(buf[12]))

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	if len(stn.File) == 0 {
		return errors.New(ce.ErrFileIsNotSet)
	}

	if len(stn.ServerHost) == 0 {
		return errors.New(ce.ErrServerHostIsNotSet)
	}

	if stn.ServerPort == 0 {
		return errors.New(ce.ErrServerPortIsNotSet)
	}

	if len(stn.ServerModeStr) == 0 {
		return errors.New(ErrServerModeIsNotSet)
	} else {
		if (stn.ServerModeStr != ServerModeHttp) &&
			(stn.ServerModeStr != ServerModeHttps) {
			return errors.New(ErrServerMode)
		}
	}

	if stn.ServerModeId == 0 {
		return errors.New(ErrServerModeIsNotSet)
	} else {
		if (stn.ServerModeId != ServerModeIdHttp) &&
			(stn.ServerModeId != ServerModeIdHttps) {
			return errors.New(ErrServerMode)
		}
	}

	switch stn.ServerModeStr {
	case ServerModeHttp:
		// Keys are not required.
	case ServerModeHttps:
		if len(stn.CertFile) == 0 {
			return errors.New(ErrCertFileIsNotSet)
		}
		if len(stn.KeyFile) == 0 {
			return errors.New(ErrKeyFileIsNotSet)
		}
	default:
		return errors.New(ErrServerMode)
	}

	if len(stn.DbHost) == 0 {
		return errors.New(ce.ErrClientHostIsNotSet)
	}

	if stn.DbPortA == 0 {
		return errors.New(ce.ErrClientPortIsNotSet)
	}

	if stn.DbPortB == 0 {
		return errors.New(ce.ErrClientPortIsNotSet)
	}

	if stn.DbClientPoolSize == 0 {
		return errors.New(ErrDbClientPoolSize)
	}

	if len(stn.FileExtension) == 0 {
		return errors.New(ErrFileExtensionIsNotSet)
	}

	if len(stn.MimeType) == 0 {
		return errors.New(ErrMimeTypeIsNotSet)
	}

	if stn.HttpCacheControlMaxAge == 0 {
		return errors.New(ErrHttpCacheControlMaxAge)
	}

	// AllowedOriginForCORS is not checked as it may be empty.

	return nil
}
