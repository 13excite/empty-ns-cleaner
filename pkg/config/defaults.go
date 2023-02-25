package config

import (
	"fmt"
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v2"
)

const DefaultConfigPath = "/data/namespaces_cleaner.yaml"

// C is the global configuration
var C = Config{}

// Defaults returns config's object with default values
func (c *Config) Defaults() {
	c.NumWorkers = 3
	c.ProtectedNS = []string{
		"default",
		"kube-node-lease",
		"kube-public",
		"kube-system",
		"local-path-storage",
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
	c.Logger.Level = "debug"
	c.Logger.Encoding = "console"
	c.Logger.Color = true
	c.Logger.DisableStacktrace = true
	c.Logger.DisableCaller = false
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
	// Sort slice. Binary search is used on this slice
	sort.Strings(c.ProtectedNS)
	return nil

}
