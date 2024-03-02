package logger

type LoggerConfig struct {
	Level       Level       `yaml:"level"`
	Encoding    string      `yaml:"encoding"`
	Path        string      `yaml:"path"`
	Environment Environment `yaml:"environment"`
}

func (c *LoggerConfig) Validate() error {
	return nil // todo write validator
}
