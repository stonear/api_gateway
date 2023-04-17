package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const (
	READ_TIMEOUT  = time.Duration(1) * time.Minute
	WRITE_TIMEOUT = time.Duration(1) * time.Minute
)

type server struct {
	log     *logrus.Logger
	port    string
	handler http.Handler
}

func New(log *logrus.Logger, port string, mux *httprouter.Router) *server {
	handler := logRequest(log, mux)
	return &server{
		log:     log,
		port:    port,
		handler: handler,
	}
}

func logRequest(log *logrus.Logger, mux *httprouter.Router) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		httpLog := &HttpLog{
			ResponseWriter: writer,
			status:         200,
		}
		mux.ServeHTTP(httpLog, request)
		log.Println(httpLog.status, request.RemoteAddr, request.Method, request.URL)
	})
}

func (s *server) RunGracefully() {
	srv := &http.Server{
		Handler:      s.handler,
		Addr:         fmt.Sprintf(`:%s`, s.port),
		ReadTimeout:  READ_TIMEOUT,
		WriteTimeout: WRITE_TIMEOUT,
	}

	s.log.Infoln("running server at port: ", s.port)
	go s.listenAndServe(srv)
	s.waitForShutdown(srv)
}

func (s *server) listenAndServe(apiServer *http.Server) {
	err := apiServer.ListenAndServe()

	if err != nil && errors.Is(err, http.ErrServerClosed) {
		s.log.Warnln("gracefully shutdown server")
	} else {
		s.log.Errorln("err closed server, err: ", err)
	}
}

func (s *server) waitForShutdown(apiServer *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-sig

	s.log.Warn("shutting down")

	if err := apiServer.Shutdown(context.Background()); err != nil {
		s.log.Println(err)
	}

	s.log.Warn("shutdown complete")
}

// HttpLog wraps a http.ResponseWriter and records the status
type HttpLog struct {
	http.ResponseWriter
	status int
}

func (h *HttpLog) Write(p []byte) (int, error) {
	return h.ResponseWriter.Write(p)
}

// WriteHeader overrides ResponseWriter.WriteHeader to keep track of the response code
func (h *HttpLog) WriteHeader(status int) {
	h.status = status
	h.ResponseWriter.WriteHeader(status)
}
