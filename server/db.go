package server

import (
	"log"
	"time"

	ce "github.com/vault-thirteen/SFRODB/common/error"
)

func (srv *Server) reconnectDb() {
	defer srv.subRoutines.Done()

	for {
		if srv.mustStop.Load() {
			log.Println("DB re-connection has been aborted.")
			break
		}

		log.Println("Re-connecting to the DB ...")
		err := srv.dbClient.Restart(true)
		if err == nil {
			log.Println("Connection to DB has been established.")
			srv.dbReconnectionIsInProgress.Store(false)
			break
		}

		for i := 1; i <= DbClientRestartSleepTimeSec; i++ {
			if srv.mustStop.Load() {
				break
			}
			time.Sleep(time.Second)
		}
	}
}

func (srv *Server) getData(uid string) (data []byte, cerr *ce.CommonError) {
	srv.dbClientLock.Lock()
	defer srv.dbClientLock.Unlock()

	data, cerr = srv.dbClient.ShowBinary(uid)
	if cerr == nil {
		return data, nil
	}

	return nil, cerr
}
