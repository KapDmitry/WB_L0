package config

import (
	"os"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	DBName   string
	Username string
	Password string
	PortExt  string
	PortInt  string
}

func (p *PostgresConfig) Load(path string) error {
	err := godotenv.Load(path + ".env")
	if err != nil {
		return err
	}

	p.DBName = os.Getenv("POSTGRES_DB")
	p.Password = os.Getenv("POSTGRES_PASSWORD")
	p.PortExt = os.Getenv("POSTGRES_PORT_EXTERNAL")
	p.PortInt = os.Getenv("POSTGRES_PORT_INTERNAL")
	p.Username = os.Getenv("POSTGRES_USER")

	return nil
}
