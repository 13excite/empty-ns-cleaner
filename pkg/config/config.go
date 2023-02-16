package config

type Config struct {
	// TODO: add sort after the config init
	ProtectedNS             []string           `yaml:"protected_ns"`
	IgnoredResouces         []IgnoredResources `yaml:"ignored_resources"`
	RunEveeryMins           int                `yaml:"run_every_mins"`
	LogLevel                string             `yaml:"log_level"`
	LogEncoding             string             `yaml:"log_encoding"`
	LoggerColor             bool               `yaml:"logger_color"`
	LoggerDisableCaller     bool               `yaml:"logger_disable_caller"`
	LoggerDisableStacktrace bool               `yaml:"Logger_disable_stacktrace"`
}

type IgnoredResources struct {
	APIGroup string `yaml:"api_group"`
	Kind     string `yaml:"kind"`
	NameMask string `yaml:"name_mask"` // should convert to regexp string
}
