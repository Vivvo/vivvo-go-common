package consul

import (
	"github.com/hashicorp/consul/api"
	"os"
)

func GetClient(address ...string) (*api.Client, error) {
	consulInsecure := os.Getenv("CONSUL_INSECURE")
	consulHost := "consul.service.consul"
	consulHostEnv := os.Getenv("CONSUL_HOST")

	if consulHostEnv != "" {
		consulHost = consulHostEnv
	}

	if len(address) == 1 {
		consulHost = address[0]
	}

	consulConfig := api.Config{
		Address:   "https://" + consulHost + ":8501",
		TLSConfig: api.TLSConfig{InsecureSkipVerify: true},
	}
	if consulInsecure == "true" {
		consulConfig = api.Config{
			Address: "http://" + consulHost + ":8500",
		}
	}
	return api.NewClient(&consulConfig)
}
