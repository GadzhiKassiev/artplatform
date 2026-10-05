package config

type AuthConfig struct {
	BaseConfig
	Port string
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		BaseConfig: loadBase(),
		Port:       getEnv("PORT", "8001"),
	}
}
