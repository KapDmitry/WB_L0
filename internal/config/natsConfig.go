package config

import (
	"os"

	"github.com/joho/godotenv"
)

type NATSConfig struct {
	ClusterID string `json:"ClusterID"`
	Subject   string `json:"Subject"`
	URL       string `json:"URL"`
}

func (n *NATSConfig) Load(path string) error {
	err := godotenv.Load(path + ".env")
	if err != nil {
		return err
	}

	n.ClusterID = os.Getenv("CLUSTER_ID")
	n.Subject = os.Getenv("SUBJECT")
	n.URL = os.Getenv("URL")

	return nil
}
