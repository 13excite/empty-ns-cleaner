package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const DefaultConfigPath = "/data/namespaces_cleaner.yaml"

// C is the global configuration
var C = Config{}

// Defaults returns config's object with default values
func (c *Config) Defaults() {
	c.ProtectedNS = []string{
		"default",
		"kube-public",
		"kube-system",
		"local-path-storage",
		"kube-node-lease",
	}
	c.RunEveeryMins = 1
	c.IgnoredResouces = []IgnoredResources{
		{
			NameMask: "kube-root-ca.crt",
			Kind:     "ConfigMap",
			APIGroup: "",
		},
		{
			NameMask: `default-token-\w+$`,
			Kind:     "Secret",
			APIGroup: "",
		},
		{
			NameMask: `^default$`,
			Kind:     "ServiceAccount",
			APIGroup: "",
		},
	}
	c.LogLevel = "debug"
	c.LogEncoding = "console"
	c.LoggerColor = true
	c.LoggerDisableStacktrace = true
	c.LoggerDisableCaller = false
}

func (c *Config) ReadConfig(configPath string) error {
	if configPath == "" {
		configPath = DefaultConfigPath
	}
	yamlConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlConfig, &c)
	if err != nil {
		return fmt.Errorf("could not unmarshal config %v", c)
	}
	return nil

}
