package config

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "5432"),
		User:     GetEnv("DB_USER", "postgres"),
		Password: GetEnv("DB_PASSWORD", "postgres"),
		Name:     GetEnv("DB_NAME", "glamping_db"),
		SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		Timezone: GetEnv("DB_TIMEZONE", "Asia/Jakarta"),
	}
}

func (d *DatabaseConfig) DSN() string {
	return "host=" + d.Host +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.Name +
		" port=" + d.Port +
		" sslmode=" + d.SSLMode +
		" TimeZone=" + d.Timezone
}
