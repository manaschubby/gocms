package config

type Config struct {
	// Database Creds
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresHost     string
	PostgresPort     string
	PostgresSSLMode  string

	//
}

func New() (*Config, error) {
	s := Config{}
	err := s.LoadEnv()
	if err != nil {
		return nil, err
	}
	return &s, nil
}
