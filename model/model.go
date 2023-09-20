package model

type Config struct {
	Server struct {
		Url     string `yaml:"url"`
		Req     string `yaml:"req"`
		Limiter int    `yaml:"limiter"`
	} `yaml:"server"`

	Files struct {
		Path string `yaml:"path"`
	} `yaml:"files"`

	Logger struct {
		Path string `yaml:"path"`
	} `yaml:"logger"`
}
