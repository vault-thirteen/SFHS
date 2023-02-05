package server

import "github.com/vault-thirteen/SFRODB/common"

func (srv *Server) getData(uid string) (data []byte, err error, de *common.Error) {
	data, err = srv.dbClient.GetBinary(uid)
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
