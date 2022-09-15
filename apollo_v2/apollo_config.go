package apollo_v2

type SimpleConfig struct {
	Cluster string `json:"cluster" toml:"cluster"`
	Host    string `json:"host" toml:"host"`
	AppID   string `json:"appId" toml:"app_id"`
}
