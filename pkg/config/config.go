package config

type Config struct {
	// TODO: add sort after the config init
	ProtectedNS     []string           `yaml:"protected_ns"`
	IgnoredResouces []IgnoredResources `yaml:"ignored_resources"`
	RunEveeryMins   int                `yaml:"run_every_mins"`
	DebugMode       bool               `yaml:"debug_mode"`
}

type IgnoredResources struct {
	APIGroup string `yaml:"api_group"`
	Kind     string `yaml:"kind"`
	NameMask string `yaml:"name_mask"` // should convert to regexp string
}
