package mustargs

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Args         []Arg    `yaml:"args"`
	FilePatterns []string `yaml:"file_patterns,omitempty"`
	FuncPatterns []string `yaml:"func_patterns,omitempty"`
	RecvPatterns []string `yaml:"recv_patterns,omitempty"`
}

type Arg struct {
	Type  string `yaml:"type"`
	Index *int   `yaml:"index,omitempty"`
}
