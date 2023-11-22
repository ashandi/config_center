package test

import (
	"config_center/api/types"
	"config_center/internal/config"
	"config_center/internal/server/handlers/config_request_handler"
	"config_center/test/mocks"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/go-redis/redismock/v8"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfigRequest_BadRequest(t *testing.T) {
	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "empty_request",
			url:  "/config",
		},
		{
			name: "missing_app_version",
			url:  "/config?platform=android",
		},
		{
			name: "missing_platform",
			url:  "/config?appVersion=1.0.0",
		},
		{
			name: "incorrect_appVersion",
			url:  "/config?appVersion=test&platform=android",
		},
		{
			name: "incorrect_platform",
			url:  "/config?appVersion=1.0.0&platform=test",
		},
		{
			name: "incorrect_assets_version",
			url:  "/config?appVersion=1.0.0&platform=android&assetsVersion=test",
		},
		{
			name: "incorrect_definitions_version",
			url:  "/config?appVersion=1.0.0&platform=android&definitionsVersion=test",
		},
	}

	logger := zaptest.NewLogger(t)
	depsRepo := mocks.MockDependenciesRepository{}
	rdb, _ := redismock.NewClientMock()
	handler := config_request_handler.NewBaseHandler(&depsRepo, rdb, &config.Config{}, logger)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			handler.Handle(recorder, req)
			res := recorder.Result()

			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestConfigRequest_200Response(t *testing.T) {
	cfg := &config.Config{
		AppVersionRequired:         "1.0.0",
		AppVersionStore:            "1.1.1",
		AppBackendEntrypoint:       "backend.entrypoint.com",
		AppNotificationsEntrypoint: "notifications.com",
	}

	expectedResponse := types.ConfigResponse{
		Version: types.AppVersion{
			Required: cfg.AppVersionRequired,
			Store:    cfg.AppVersionStore,
		},
		BackendEntryPoint: types.JsonRpcService{
			JsonRpcUrl: cfg.AppBackendEntrypoint,
		},
		Notifications: types.JsonRpcService{
			JsonRpcUrl: cfg.AppNotificationsEntrypoint,
		},
		Assets: types.Dependency{
			Version: "1.12.234",
			Hash:    "test1234",
			Urls: []string{
				"test1.com",
				"test2.com",
			},
		},
		Definitions: types.Dependency{
			Version: "1.0.789",
			Hash:    "test4321",
			Urls: []string{
				"test3.com",
				"test4.com",
			},
		},
	}

	logger := zaptest.NewLogger(t)
	depsRepo := mocks.MockDependenciesRepository{}
	rdb, _ := redismock.NewClientMock()
	handler := config_request_handler.NewBaseHandler(&depsRepo, rdb, cfg, logger)

	reqParamMajor := 1
	reqParamMinor := 0
	reqParamPatch := 0
	reqParamPlatform := "android"

	depsRepo.On("FindByMajor", "assets", reqParamPlatform, reqParamMajor).
		Return(expectedResponse.Assets, nil)
	depsRepo.On("FindByMajorMinor", "definitions", reqParamPlatform, reqParamMajor, reqParamMinor).
		Return(expectedResponse.Definitions, nil)

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/config?appVersion=%d.%d.%d&platform=%s", reqParamMajor, reqParamMinor, reqParamPatch, reqParamPlatform),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)
	res := recorder.Result()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resObj types.ConfigResponse
	err = json.Unmarshal(resBody, &resObj)
	if err != nil {
		t.Fatal(err)
	}

	depsRepo.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expectedResponse, resObj)
}

func TestConfigRequest_200Response_ConcreteDependencyVersions(t *testing.T) {
	cfg := &config.Config{
		AppVersionRequired:         "1.0.0",
		AppVersionStore:            "1.1.1",
		AppBackendEntrypoint:       "backend.entrypoint.com",
		AppNotificationsEntrypoint: "notifications.com",
	}

	reqParamAppVersionMajor := 1
	reqParamAppVersionMinor := 0
	reqParamAppVersionPatch := 0
	reqParamPlatform := "android"
	reqParamAssetsVersionMajor := 5
	reqParamAssetsVersionMinor := 4
	reqParamAssetsVersionPatch := 23
	reqParamDefinitionsVersionMajor := 42
	reqParamDefinitionsVersionMinor := 14
	reqParamDefinitionsVersionPatch := 756

	expectedResponse := types.ConfigResponse{
		Version: types.AppVersion{
			Required: cfg.AppVersionRequired,
			Store:    cfg.AppVersionStore,
		},
		BackendEntryPoint: types.JsonRpcService{
			JsonRpcUrl: cfg.AppBackendEntrypoint,
		},
		Notifications: types.JsonRpcService{
			JsonRpcUrl: cfg.AppNotificationsEntrypoint,
		},
		Assets: types.Dependency{
			Version: fmt.Sprintf(
				"%d.%d.%d",
				reqParamAssetsVersionMajor,
				reqParamAssetsVersionMinor,
				reqParamAssetsVersionPatch,
			),
			Hash: "test1234",
			Urls: []string{
				"test1.com",
				"test2.com",
			},
		},
		Definitions: types.Dependency{
			Version: fmt.Sprintf(
				"%d.%d.%d",
				reqParamDefinitionsVersionMajor,
				reqParamDefinitionsVersionMinor,
				reqParamDefinitionsVersionPatch,
			),
			Hash: "test4321",
			Urls: []string{
				"test3.com",
				"test4.com",
			},
		},
	}

	logger := zaptest.NewLogger(t)
	depsRepo := mocks.MockDependenciesRepository{}
	rdb, _ := redismock.NewClientMock()
	handler := config_request_handler.NewBaseHandler(&depsRepo, rdb, cfg, logger)

	depsRepo.On(
		"FindByMajorMinorPatch",
		"assets",
		reqParamPlatform,
		reqParamAssetsVersionMajor,
		reqParamAssetsVersionMinor,
		reqParamAssetsVersionPatch,
	).
		Return(expectedResponse.Assets, nil)
	depsRepo.On(
		"FindByMajorMinorPatch",
		"definitions",
		reqParamPlatform,
		reqParamDefinitionsVersionMajor,
		reqParamDefinitionsVersionMinor,
		reqParamDefinitionsVersionPatch,
	).
		Return(expectedResponse.Definitions, nil)

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"/config?appVersion=%d.%d.%d&platform=%s&assetsVersion=%d.%d.%d&definitionsVersion=%d.%d.%d",
			reqParamAppVersionMajor,
			reqParamAppVersionMinor,
			reqParamAppVersionPatch,
			reqParamPlatform,
			reqParamAssetsVersionMajor,
			reqParamAssetsVersionMinor,
			reqParamAssetsVersionPatch,
			reqParamDefinitionsVersionMajor,
			reqParamDefinitionsVersionMinor,
			reqParamDefinitionsVersionPatch,
		),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	handler.Handle(recorder, req)
	res := recorder.Result()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resObj types.ConfigResponse
	err = json.Unmarshal(resBody, &resObj)
	if err != nil {
		t.Fatal(err)
	}

	depsRepo.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, expectedResponse, resObj)
}
