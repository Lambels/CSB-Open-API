package csb

type Config struct {
	Addr        string `json:"addr"`
	SqliteDSN   string `json:"sqlite_dsn"`
	EngageToken string `json:"engage_token"`
}
