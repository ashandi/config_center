package config_request_handler

import (
	"config_center/internal/utils"
	"github.com/go-playground/validator/v10"
	"net/url"
)

type definitionsHandler struct {
	validate               *validator.Validate
	dependenciesRepository dependenciesRepository
}

func (h *definitionsHandler) buildRequest(httpRequest url.Values, req *Request) error {
	definitionsVersion := httpRequest.Get("definitionsVersion")

	if err := h.validate.Var(definitionsVersion, "omitempty,semver"); err != nil {
		return err
	}

	(*req)["definitionsVersion"] = definitionsVersion

	return nil
}

func (h *definitionsHandler) buildResponse(req *Request, res *Response) error {
	platform := (*req)["platform"]

	concreteVersion := (*req)["definitionsVersion"]
	if len(concreteVersion) > 0 {
		major, minor, patch, err := utils.SplitSemVer(concreteVersion)
		if err != nil {
			return err
		}

		dep, err := h.dependenciesRepository.FindByMajorMinorPatch("definitions", platform, major, minor, patch)
		if err != nil {
			return err
		}

		res.Definitions = dep
		return nil
	}

	major, minor, _, err := utils.SplitSemVer((*req)["appVersion"])
	if err != nil {
		return err
	}

	dep, err := h.dependenciesRepository.FindByMajorMinor("definitions", platform, major, minor)
	if err != nil {
		return err
	}

	res.Definitions = dep
	return nil
}
