package config_request_handler

import (
	"config_center/api/types"
	"config_center/internal/config"
	"net/url"
)

type appVersionHandler struct {
	cfg *config.Config
}

func (h *appVersionHandler) buildRequest(httpRequest url.Values, req *Request) error {
	return nil
}

func (h *appVersionHandler) buildResponse(req *Request, res *Response) error {
	res.Version = types.AppVersion{
		Required: h.cfg.AppVersionRequired,
		Store:    h.cfg.AppVersionStore,
	}

	return nil
}
