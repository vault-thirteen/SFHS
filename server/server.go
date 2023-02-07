package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	ss "github.com/vault-thirteen/SFHS/server/settings"
	"github.com/vault-thirteen/SFRODB/client"
	cs "github.com/vault-thirteen/SFRODB/client/settings"
)

type Server struct {
	settings  *ss.Settings
	listenDsn string
	dbDsnA    string
	dbDsnB    string

	// HTTP server.
	httpServer *http.Server

	// DB client.
	dbClient     *client.Client
	dbClientLock *sync.Mutex

	// Channel for an external controller. When a message comes from this
	// channel, a controller must stop this server. The server does not stop
	// itself.
	mustBeStopped chan bool

	// Internal control structures.
	mustStop                   *atomic.Bool
	subRoutines                *sync.WaitGroup
	httpErrors                 chan error
	dbErrors                   chan error
	dbReconnectionIsInProgress *atomic.Bool
}

const (
	HttpStatusCodeOnError       = 0
	DbClientRestartSleepTimeSec = 15
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
		dbClientLock:  new(sync.Mutex),
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

	srv.mustStop = new(atomic.Bool)
	srv.mustStop.Store(false)
	srv.subRoutines = new(sync.WaitGroup)
	srv.httpErrors = make(chan error, 8)
	srv.dbErrors = make(chan error, 8)
	srv.dbReconnectionIsInProgress = new(atomic.Bool)
	srv.dbReconnectionIsInProgress.Store(false)

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

func (srv *Server) Stop(forcibly bool) (err error) {
	srv.mustStop.Store(true)

	ctx, cf := context.WithTimeout(context.Background(), time.Minute)
	defer cf()
	err = srv.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	if forcibly {
		_ = srv.dbClient.Stop()
	} else {
		err = srv.dbClient.Stop()
		if err != nil {
			return err
		}
	}

	close(srv.httpErrors)
	close(srv.dbErrors)

	srv.subRoutines.Wait()

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
		log.Println("DB error: " + err.Error()) //TODO:Debug.

		if !srv.dbReconnectionIsInProgress.Load() {
			srv.dbReconnectionIsInProgress.Store(true)
			srv.subRoutines.Add(1)
			go srv.reconnectDb()
		}
	}

	log.Println("DB error listener has stopped.")
}
