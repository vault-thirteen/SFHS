package settings

import (
	"errors"
	"os"
	"strings"

	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/helper"
	"github.com/vault-thirteen/auxie/reader"
	"github.com/vault-thirteen/errorz"
)

const (
	ErrCertFileIsNotSet       = "certificate file is not set"
	ErrKeyFileIsNotSet        = "key file is not set"
	ErrFileExtensionIsNotSet  = "file exyension is not set"
	ErrMimeTypeIsNotSet       = "MIME type is not set"
	ErrDbClientPoolSize       = "DB client pool size is not set"
	ErrHttpCacheControlMaxAge = "HTTP cache control max-age error"
)

const (
	ContentDispositionInline = "inline"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the File with these Settings.
	File string

	// Server's Host Name.
	ServerHost string

	// Server's Listen Port.
	ServerPort uint16

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
	var buf = make([][]byte, 11)

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

	// Certificate and Key.
	stn.CertFile = strings.TrimSpace(string(buf[2]))
	stn.KeyFile = strings.TrimSpace(string(buf[3]))

	// Database.
	stn.DbHost = strings.TrimSpace(string(buf[4]))

	stn.DbPortA, err = helper.ParseUint16(strings.TrimSpace(string(buf[5])))
	if err != nil {
		return stn, err
	}

	stn.DbPortB, err = helper.ParseUint16(strings.TrimSpace(string(buf[6])))
	if err != nil {
		return stn, err
	}

	stn.DbClientPoolSize, err = helper.ParseInt(strings.TrimSpace(string(buf[7])))
	if err != nil {
		return stn, err
	}

	stn.FileExtension = strings.TrimSpace(string(buf[8]))
	stn.MimeType = strings.TrimSpace(string(buf[9]))

	stn.HttpCacheControlMaxAge, err = helper.ParseUint(strings.TrimSpace(string(buf[10])))
	if err != nil {
		return stn, err
	}

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

	if len(stn.CertFile) == 0 {
		return errors.New(ErrCertFileIsNotSet)
	}

	if len(stn.KeyFile) == 0 {
		return errors.New(ErrKeyFileIsNotSet)
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

	return nil
}
