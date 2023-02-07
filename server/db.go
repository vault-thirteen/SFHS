package server

import (
	"log"
	"time"

	"github.com/vault-thirteen/SFRODB/common"
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

func (srv *Server) getData(uid string) (data []byte, err error, de *common.Error) {
	data, err = srv.dbClient.ShowBinary(uid)
	if err == nil {
		return data, nil, nil
	}

	var ok bool
	de, ok = err.(*common.Error)
	if !ok {
		return nil, err, nil
	}

	return nil, err, de
}
