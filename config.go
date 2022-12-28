package csb

type Config struct {
	HTTP   httpConfig   `json:"http"`
	Sqlite sqliteConfig `json:"sqlite"`
	Engage engageConfig `json:"engage"`
}

type engageConfig struct {
	Token   string `json:"token"`
	SaveNew bool   `json:"save_new"`
}

type httpConfig struct {
	AddrBackend  string `json:"addr_backend"`
	AddrFrontend string `json:"addr_frontend"`
}

type sqliteConfig struct {
	DSN string `json:"dsn"`
}
