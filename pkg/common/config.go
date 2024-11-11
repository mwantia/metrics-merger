package common

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address        string           `mapstructure:"address"`
	MergeLabel     string           `mapstructure:"merge_label"`
	ScrapeInterval string           `mapstructure:"scrape_interval"`
	Endpoints      []EndpointConfig `mapstructure:"endpoints"`
}

type EndpointConfig struct {
	Name    string `mapstructure:"name"`
	Address string `mapstructure:"address"`
}

func LoadServerConfig(path string) (*ServerConfig, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetEnvPrefix("")
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := &ServerConfig{}

	if err := v.Unmarshal(cfg, viper.DecodeHook(func(src, dst reflect.Type, data interface{}) (interface{}, error) {
		if src.Kind() != reflect.String {
			return data, nil
		}

		str := data.(string)
		if strings.Contains(str, "${") || strings.HasPrefix(str, "$") {

			envVar := strings.Trim(str, "${}")
			envVar = strings.TrimPrefix(envVar, "$")

			if value, exists := os.LookupEnv(envVar); exists {
				return value, nil
			}
		}
		return data, nil
	})); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return cfg, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (c *ServerConfig) Validate() error {
	if strings.TrimSpace(c.Address) == "" {
		c.Address = ":8080"
	}
	if strings.TrimSpace(c.MergeLabel) == "" {
		c.MergeLabel = "name"
	}
	if strings.TrimSpace(c.ScrapeInterval) == "" {
		c.ScrapeInterval = "10s"
	}

	if len(c.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint must be defined")
	}

	uniques := make(map[string]bool)
	for _, endpoint := range c.Endpoints {
		if uniques[endpoint.Name] {
			return fmt.Errorf("duplicate endpoint name found: %s", endpoint.Name)
		}

		if strings.TrimSpace(endpoint.Address) == "" {
			return fmt.Errorf("endpoint '%s' must define an address", endpoint.Name)
		}

		uniques[endpoint.Name] = true
	}

	return nil
}
