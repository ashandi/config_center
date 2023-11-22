package config_request_handler

import (
	"config_center/internal/utils"
	"github.com/go-playground/validator/v10"
	"net/url"
)

type assetsHandler struct {
	validate               *validator.Validate
	dependenciesRepository dependenciesRepository
}

func (h *assetsHandler) buildRequest(httpRequest url.Values, req *Request) error {
	assetsVersion := httpRequest.Get("assetsVersion")

	if err := h.validate.Var(assetsVersion, "omitempty,semver"); err != nil {
		return err
	}

	(*req)["assetsVersion"] = assetsVersion

	return nil
}

func (h *assetsHandler) buildResponse(req *Request, res *Response) error {
	platform := (*req)["platform"]

	concreteVersion := (*req)["assetsVersion"]
	if len(concreteVersion) > 0 {
		major, minor, patch, err := utils.SplitSemVer(concreteVersion)
		if err != nil {
			return err
		}

		dep, err := h.dependenciesRepository.FindByMajorMinorPatch("assets", platform, major, minor, patch)
		if err != nil {
			return err
		}

		res.Assets = dep
		return nil
	}

	major, _, _, err := utils.SplitSemVer((*req)["appVersion"])
	if err != nil {
		return err
	}

	dep, err := h.dependenciesRepository.FindByMajor("assets", platform, major)
	if err != nil {
		return err
	}

	res.Assets = dep
	return nil
}
