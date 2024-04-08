package config

const (
	AppSrvName = "core-api"
)

type Config struct {
	DataBaseUrl string `json:"data_base_url"`
	ServiceName string
	Port        int8
}

// New loads the environment variable, parses them to the Config struct
// and returns an instance of Config
func New() *Config {

	//if loadErr := godotenv.Load(".env"); loadErr != nil {
	//	log.Printf("[Env]: unable to load .env file %v", loadErr)
	//}
	//
	var cfg Config
	//if parseErr := env.Parse(&cfg); parseErr != nil {
	//	log.Fatalf("[Env]: failed to parse environment variables: %v", parseErr)
	//}
	cfg.ServiceName = AppSrvName

	return &cfg
}
