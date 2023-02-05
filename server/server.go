package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	ss "github.com/vault-thirteen/SFHS/server/settings"
	"github.com/vault-thirteen/SFRODB/client"
	cs "github.com/vault-thirteen/SFRODB/client/settings"
	"github.com/vault-thirteen/SFRODB/common"
	hdr "github.com/vault-thirteen/header"
)

type Server struct {
	settings  *ss.Settings
	listenDsn string
	dbDsnA    string
	dbDsnB    string

	httpServer    *http.Server
	dbClient      *client.Client
	mustBeStopped chan bool

	httpErrors chan error
	dbErrors   chan error
}

const (
	HttpStatusCodeOnError = 0
)

func NewServer(stn *ss.Settings) (srv *Server, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	srv = &Server{
		settings:      stn,
		listenDsn:     fmt.Sprintf("%s:%d", stn.ServerHost, stn.ServerPort),
		dbDsnA:        fmt.Sprintf("%s:%d", stn.DbHost, stn.DbPortA),
		dbDsnB:        fmt.Sprintf("%s:%d", stn.DbHost, stn.DbPortB),
		mustBeStopped: make(chan bool, 2),
	}

	// HTTP Server.
	srv.httpServer = &http.Server{
		Addr:    srv.listenDsn,
		Handler: http.Handler(http.HandlerFunc(srv.httpRouter)),
	}

	// DB Client.
	var dbClientSettings *cs.Settings
	dbClientSettings, err = cs.NewSettings(
		srv.settings.DbHost,
		srv.settings.DbPortA,
		srv.settings.DbPortB,
		cs.ResponseMessageLengthLimitDefault,
	)
	if err != nil {
		return nil, err
	}

	srv.dbClient, err = client.NewClient(dbClientSettings)
	if err != nil {
		return nil, err
	}

	srv.httpErrors = make(chan error, 8)
	srv.dbErrors = make(chan error, 8)

	return srv, nil
}

func (srv *Server) GetListenDsn() (dsn string) {
	return srv.listenDsn
}

func (srv *Server) GetDbDsnA() (dsn string) {
	return srv.dbDsnA
}

func (srv *Server) GetDbDsnB() (dsn string) {
	return srv.dbDsnB
}

func (srv *Server) GetStopChannel() *chan bool {
	return &srv.mustBeStopped
}

func (srv *Server) Start() (err error) {
	srv.startHttpServer()

	err = srv.dbClient.Start()
	if err != nil {
		return err
	}

	go srv.listenForHttpErrors()
	go srv.listenForDbErrors()

	return nil
}

func (srv *Server) Stop() (err error) {
	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	err = srv.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = srv.dbClient.Stop()
	if err != nil {
		return err
	}

	close(srv.httpErrors)
	close(srv.dbErrors)

	return nil
}

func (srv *Server) startHttpServer() {
	go func() {
		listenError := srv.httpServer.ListenAndServeTLS(srv.settings.CertFile, srv.settings.KeyFile)
		if (listenError != nil) && (listenError != http.ErrServerClosed) {
			srv.httpErrors <- listenError
		}
	}()
}

func (srv *Server) listenForHttpErrors() {
	for err := range srv.httpErrors {
		log.Println("Server error: " + err.Error())
		srv.mustBeStopped <- true
	}

	log.Println("HTTP error listener has stopped.")
}

func (srv *Server) listenForDbErrors() {
	for err := range srv.dbErrors {
		log.Println("DB error: " + err.Error())
		//TODO
	}

	log.Println("DB error listener has stopped.")
}

func (srv *Server) getContentDisposition(uid string) string {
	return ss.ContentDispositionInline +
		`; filename="` + filepath.Base(uid) + srv.settings.FileExtension + `""`
}

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	uid := req.URL.Path[1:]

	data, err, de := srv.getData(uid)
	if err != nil {
		srv.processError(rw, err, de)
		return
	}

	srv.respondWithData(rw, data)
}

func (srv *Server) processError(
	rw http.ResponseWriter,
	err error, // This error is non-null.
	de *common.Error,
) {
	if de == nil {
		log.Println(err)
		rw.WriteHeader(HttpStatusCodeOnError)
		return
	}

	if de.IsClientError() {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if de.IsServerError() {
		rw.WriteHeader(http.StatusInternalServerError)
		srv.dbErrors <- err
		return
	}

	log.Println(err)
	rw.WriteHeader(HttpStatusCodeOnError)
}

func (srv *Server) respondWithData(
	rw http.ResponseWriter,
	//uid string,
	data []byte,
) {
	rw.Header().Set(hdr.HttpHeaderContentType, srv.settings.MimeType)
	//rw.Header().Set(hdr.HttpHeaderContentDisposition, srv.getContentDisposition(uid))
	rw.WriteHeader(http.StatusOK)

	_, err := rw.Write(data)
	if err != nil {
		log.Println(err)
	}
}
