package config_request_handler

import (
	"github.com/go-playground/validator/v10"
	"net/url"
)

type commonRequestValidator struct {
	validate *validator.Validate
}

type RequestParams struct {
	AppVersion string `validate:"required,semver"`
	Platform   string `validate:"required,oneof=android ios"`
}

func (h *commonRequestValidator) buildRequest(httpRequest url.Values, req *Request) error {
	params := RequestParams{
		AppVersion: httpRequest.Get("appVersion"),
		Platform:   httpRequest.Get("platform"),
	}

	if err := h.validate.Struct(params); err != nil {
		return err
	}

	(*req)["appVersion"] = params.AppVersion
	(*req)["platform"] = params.Platform

	return nil
}

func (h *commonRequestValidator) buildResponse(req *Request, res *Response) error {
	return nil
}
