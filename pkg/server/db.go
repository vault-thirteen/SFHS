package server

import (
	"time"

	"github.com/vault-thirteen/SFRODB/pkg/client"
	ce "github.com/vault-thirteen/SFRODB/pkg/common/error"
	"github.com/vault-thirteen/errorz"
)

const (
	DbClientPoolTakeRetryDelayMs = 100

	// DbClientPoolTakeAttemptsCountLimit is the limit for number of attempts
	// to try to get a client from the pool.
	// (100ms * 10) * 60 => 1 Minute.
	DbClientPoolTakeAttemptsCountLimit = 10 * 60
)

func (srv *Server) takeClient() (cli *client.Client, err error) {
	var attemptsCount = 1
	cli, err = srv.poolOfClients.GiveIdleClient()
	for (err != nil) && (attemptsCount <= DbClientPoolTakeAttemptsCountLimit) {
		time.Sleep(time.Millisecond * DbClientPoolTakeRetryDelayMs)
		attemptsCount++
		cli, err = srv.poolOfClients.GiveIdleClient()
	}

	return cli, err
}

func (srv *Server) returnClient(cli *client.Client, clientError *ce.CommonError) (err error) {
	var isClientBoken = (clientError != nil) && (clientError.IsServerError())
	err = srv.poolOfClients.TakeIdleClient(cli.GetId(), isClientBoken)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) getData(uid string) (data []byte, cerr *ce.CommonError) {
	srv.dbClientLock.Lock()
	defer srv.dbClientLock.Unlock()

	var cli *client.Client
	var err error
	cli, err = srv.takeClient()
	if err != nil {
		cerr = ce.NewClientError(err.Error(), 0, client.ClientIdNone)
		return nil, cerr
	}

	defer func() {
		err = srv.returnClient(cli, cerr)
		if err != nil {
			combinedError := errorz.Combine(cerr, err)
			cerr = ce.NewClientError(combinedError.Error(), 0, cli.GetId())
		}
		return
	}()

	data, cerr = cli.ShowData(uid)
	if cerr != nil {
		return nil, cerr
	}

	return data, nil
}
