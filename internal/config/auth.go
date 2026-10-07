package config

type AuthConfig struct {
	BaseConfig
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		BaseConfig: loadBase("8001"),
	}
}
