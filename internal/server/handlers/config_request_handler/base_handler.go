package config_request_handler

import (
	"config_center/api/types"
	"config_center/internal/config"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Request map[string]string

type Response struct {
	types.ConfigResponse
}

type RequestHandler interface {
	buildRequest(httpRequest url.Values, req *Request) error
	buildResponse(req *Request, res *Response) error
}

type dependenciesRepository interface {
	FindByMajor(table, platform string, major int) (types.Dependency, error)
	FindByMajorMinor(table, platform string, major, minor int) (types.Dependency, error)
	FindByMajorMinorPatch(table, platform string, major, minor, patch int) (types.Dependency, error)
}

type BaseHandler struct {
	handlers []RequestHandler
	cfg      *config.Config
	logger   *zap.Logger
	rdb      *redis.Client
}

func NewBaseHandler(
	dependenciesRepository dependenciesRepository,
	rdb *redis.Client,
	cfg *config.Config,
	logger *zap.Logger,
) *BaseHandler {
	validate := validator.New()

	handlers := []RequestHandler{
		&commonRequestValidator{validate: validate},
		&appVersionHandler{cfg: cfg},
		&backendEntrypointHandler{cfg: cfg},
		&notificationsHandler{cfg: cfg},
		&assetsHandler{validate: validate, dependenciesRepository: dependenciesRepository},
		&definitionsHandler{validate: validate, dependenciesRepository: dependenciesRepository},
	}

	return &BaseHandler{
		handlers: handlers,
		cfg:      cfg,
		logger:   logger,
		rdb:      rdb,
	}
}

func (h *BaseHandler) Handle(w http.ResponseWriter, r *http.Request) {
	cachedResponse, err := h.rdb.Get(r.Context(), r.URL.RawQuery).Result()
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(cachedResponse))
		return
	} else if err != redis.Nil {
		h.logger.Error("an error occurred during searching cached response in redis", zap.Error(err))
	}

	req := &Request{}
	if err := h.buildRequest(r.URL.Query(), req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %s", err), http.StatusBadRequest)
		return
	}

	res := &Response{}
	if err := h.buildResponse(req, res); err != nil {
		h.logger.Error("an error occurred during building the config response", zap.Error(err))
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		h.logger.Error("an error occurred during response to json serialization", zap.Error(err))
		http.Error(w, fmt.Sprintf("Internal Server Error: %s", err), http.StatusInternalServerError)
		return
	}

	err = h.rdb.Set(r.Context(), r.URL.RawQuery, jsonResponse, h.cfg.RedisTtl).Err()
	if err != nil {
		h.logger.Error("an error occurred during saving the cached response into redis", zap.Error(err))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonResponse)
}

func (h *BaseHandler) buildRequest(httpRequest url.Values, req *Request) error {
	var errs validator.ValidationErrors

	for _, handler := range h.handlers {
		if err := handler.buildRequest(httpRequest, req); err != nil {
			errs = append(errs, err.(validator.ValidationErrors)...)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func (h *BaseHandler) buildResponse(req *Request, res *Response) error {
	for _, handler := range h.handlers {
		if err := handler.buildResponse(req, res); err != nil {
			return err
		}
	}

	return nil
}
