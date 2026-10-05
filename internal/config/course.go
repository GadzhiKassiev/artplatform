package config

type CourseConfig struct {
	BaseConfig
	Port string
}

func LoadCourseConfig() CourseConfig {
	return CourseConfig{
		BaseConfig: loadBase(),
		Port:       getEnv("PORT", "8002"),
	}
}
