package db

import "time"

type StorageConfig struct {
	EnableMock            bool          `yaml:"enable_mock"`
	Hosts                 []string      `yaml:"hosts"`
	Port                  int           `yaml:"port"`
	Database              string        `yaml:"database"`
	Username              string        `yaml:"username"`
	Password              string        `yaml:"password" env:"DB_PASSWORD"`
	SSLMode               string        `yaml:"ssl_mode"`
	ConnectionAttempts    int           `yaml:"connection_attempts"`
	InitializationTimeout time.Duration `yaml:"initialization_timeout"`
}

func (c *StorageConfig) Validate() error {
	return nil // todo
}
