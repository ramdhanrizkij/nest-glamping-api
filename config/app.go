package config

type AppConfig struct {
	Name string
	Port string
	Env  string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		Name: GetEnv("APP_NAME", "Glamping API"),
		Port: GetEnv("APP_PORT", "3000"),
		Env:  GetEnv("APP_ENV", "development"),
	}
}
