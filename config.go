package csb

// Config represents the structure of a json config file.
type Config struct {
	HTTP   httpConfig   `json:"http"`   // HTTP related configs.
	Sqlite sqliteConfig `json:"sqlite"` // Sqlite related configs.
	Engage engageConfig `json:"engage"` // Engage related configs.
}

// engageConfig holds all the config fields related to engage.
type engageConfig struct {
	// Token used for engage auth.
	Token string `json:"token"`
	// Fallback indicates wether failed queries to the database should fallback to engage.
	Fallback bool `json:"fallback"`
}

// httpConfig holds all the config fields related to http services.
type httpConfig struct {
	// AddrBackend is the http adress of the backend server.
	AddrBackend string `json:"addr_backend"`
	// AddrFrontend is the http adress of the frontend server.
	AddrFrontend string `json:"addr_frontend"`
}

// sqliteConfig holds all the config fields related to the sqlite database.
type sqliteConfig struct {
	// DSN is the data source name of the database.
	//
	// :memory: isnt allowed.
	DSN string `json:"dsn"`
	// MigrationsPath is the path to the migrations folder.
	MigrationsPath string `json:"migrations_path"`
}
