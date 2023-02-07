package server

import (
	"log"
	"net/http"
	"path/filepath"

	ss "github.com/vault-thirteen/SFHS/server/settings"
	"github.com/vault-thirteen/SFRODB/common"
	hdr "github.com/vault-thirteen/header"
)

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	uid := req.URL.Path[1:]

	data, err, de := srv.getData(uid)
	if err != nil {
		srv.processError(rw, err, de)
		return
	}

	srv.respondWithData(rw, data)
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

func (srv *Server) getContentDisposition(uid string) string {
	return ss.ContentDispositionInline +
		`; filename="` + filepath.Base(uid) + srv.settings.FileExtension + `""`
}
