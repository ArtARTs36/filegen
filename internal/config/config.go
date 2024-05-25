package config

type Config struct {
	Vars  map[string]interface{} `yaml:"vars"`
	Files []File                 `yaml:"files"`
}

type File struct {
	TemplatePath string                 `yaml:"template_path"`
	OutputPath   string                 `yaml:"output_path"`
	Vars         map[string]interface{} `yaml:"vars"`
}
