package config

type Config struct {
	NumWorkers      int                `yaml:"num_workers"`
	ProtectedNS     []string           `yaml:"protected_ns"`
	IgnoredResouces []IgnoredResources `yaml:"ignored_resources"`
	RunEveeryMins   int                `yaml:"run_every_mins"`
	Logger          Logger             `yaml:"logger"`
}

type IgnoredResources struct {
	APIGroup string `yaml:"api_group"`
	Kind     string `yaml:"kind"`
	NameMask string `yaml:"name_mask"` // should convert to regexp string
}

type Logger struct {
	Level             string `yaml:"level"`
	Encoding          string `yaml:"encoding"`
	Color             bool   `yaml:"color"`
	DisableCaller     bool   `yaml:"disable_caller"`
	DisableStacktrace bool   `yaml:"disable_stacktrace"`
}
