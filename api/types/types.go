package types

type ConfigResponse struct {
	Version           AppVersion     `json:"version"`
	BackendEntryPoint JsonRpcService `json:"backend_entry_point"`
	Assets            Dependency     `json:"assets"`
	Definitions       Dependency     `json:"definitions"`
	Notifications     JsonRpcService `json:"notifications"`
}

type AppVersion struct {
	Required string `json:"required"`
	Store    string `json:"store"`
}

type JsonRpcService struct {
	JsonRpcUrl string `json:"jsonrpc_url"`
}

type Dependency struct {
	Version string   `json:"version"`
	Hash    string   `json:"hash"`
	Urls    []string `json:"urls"`
}
