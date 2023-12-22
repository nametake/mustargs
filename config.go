package mustargs

type Config struct {
	Rules []*Rule `yaml:"rules"`
}

type Rule struct {
	Args         []*RuleArg `yaml:"args"`
	FilePatterns []string   `yaml:"file_patterns,omitempty"`
	FuncPatterns []string   `yaml:"func_patterns,omitempty"`
	RecvPatterns []string   `yaml:"recv_patterns,omitempty"`
}

type RuleArg struct {
	Type  string  `yaml:"type"`
	Index *int    `yaml:"index,omitempty"`
	Pkg   *string `yaml:"pkg,omitempty"`
	Name  *string `yaml:"name,omitempty"`
	Ptr   *bool   `yaml:"ptr,omitempty"`
}
