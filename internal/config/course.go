package config

type CourseConfig struct {
	BaseConfig
}

func LoadCourseConfig() CourseConfig {
	return CourseConfig{
		BaseConfig: loadBase("8002"),
	}
}
