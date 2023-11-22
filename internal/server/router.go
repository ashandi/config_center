package server

import (
	"config_center/internal/utils"
	"go.uber.org/zap"
	"net/http"
)

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func (s *Server) SetupHandlers(
	configRequestHandler Handler,
) {
	s.setHandler([]string{http.MethodGet}, "/config", configRequestHandler)
}

func (s *Server) setHandler(allowedMethods []string, pattern string, handler Handler) {
	s.router.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Request "+r.RequestURI, zap.String("method", r.Method))

		recorder := &statusRecorder{
			ResponseWriter: w,
		}

		if !utils.InArray(r.Method, allowedMethods) {
			http.Error(recorder, "Method Not Allowed", 405)
			return
		}

		handler.Handle(recorder, r)

		s.logger.Info("Response "+r.RequestURI, zap.Int("status", recorder.Status))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
