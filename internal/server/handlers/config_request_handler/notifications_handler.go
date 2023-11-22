package config_request_handler

import (
	"config_center/api/types"
	"config_center/internal/config"
	"net/url"
)

type notificationsHandler struct {
	cfg *config.Config
}

func (h *notificationsHandler) buildRequest(httpRequest url.Values, req *Request) error {
	return nil
}

func (h *notificationsHandler) buildResponse(req *Request, res *Response) error {
	res.Notifications = types.JsonRpcService{
		JsonRpcUrl: h.cfg.AppNotificationsEntrypoint,
	}

	return nil
}
