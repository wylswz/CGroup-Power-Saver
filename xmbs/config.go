package xmbs

type Rule struct {
	Process  string `yaml:"process,omitempty"`
	Resource string `yaml:"resource,omitempty"`
	Amount   string `yaml:"amount,omitempty"`
	When     string `yaml:"when,omitempty"`
}

type Config struct {
	Rules []Rule `yaml:"rules"`
}
