package server

import (
	"log"
	"net/http"
	"path/filepath"

	ss "github.com/vault-thirteen/SFHS/server/settings"
	ce "github.com/vault-thirteen/SFRODB/common/error"
	hdr "github.com/vault-thirteen/header"
)

func (srv *Server) httpRouter(rw http.ResponseWriter, req *http.Request) {
	uid := req.URL.Path[1:]

	data, cerr := srv.getData(uid)
	if cerr != nil {
		srv.processError(rw, cerr)
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

func (srv *Server) processError(rw http.ResponseWriter, cerr *ce.CommonError) {
	if cerr.IsClientError() {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if cerr.IsServerError() {
		rw.WriteHeader(http.StatusInternalServerError)
		srv.dbErrors <- cerr
		return
	}

	log.Println("Anomaly: " + cerr.Error())
	rw.WriteHeader(http.StatusInternalServerError)
}

func (srv *Server) getContentDisposition(uid string) string {
	return ss.ContentDispositionInline +
		`; filename="` + filepath.Base(uid) + srv.settings.FileExtension + `""`
}
