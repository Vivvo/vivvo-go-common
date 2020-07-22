package consul

import (
	"os"
	"github.com/hashicorp/consul/api"
)

func GetClient() (*api.Client, error) {
	consulInsecure := os.Getenv("CONSUL_INSECURE")
	consulHost := "consul.sd.svc.cluster.local"
	consulHostEnv := os.Getenv("CONSUL_HOST")
	if consulHostEnv != "" {
		consulHost = consulHostEnv
	}
	consulConfig := api.Config{
		Address:   "https://" + consulHost + ":8501",
		TLSConfig: api.TLSConfig{InsecureSkipVerify: true},
	}
	if consulInsecure == "true" {
		consulConfig = api.Config{
			Address:   "http://" + consulHost + ":8500",
		}
	}
	return api.NewClient(&consulConfig)
}
