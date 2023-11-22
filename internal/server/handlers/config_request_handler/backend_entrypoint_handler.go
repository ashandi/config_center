package config_request_handler

import (
	"config_center/api/types"
	"config_center/internal/config"
	"net/url"
)

type backendEntrypointHandler struct {
	cfg *config.Config
}

func (h *backendEntrypointHandler) buildRequest(httpRequest url.Values, req *Request) error {
	return nil
}

func (h *backendEntrypointHandler) buildResponse(req *Request, res *Response) error {
	res.BackendEntryPoint = types.JsonRpcService{
		JsonRpcUrl: h.cfg.AppBackendEntrypoint,
	}

	return nil
}
