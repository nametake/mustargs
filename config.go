package mustargs

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	ArgType []struct {
		TypeName string `yaml:"type_name"`
		Index    int    `yaml:"index,omitempty"`
	} `yaml:"arg_type,omitempty"`
	FilePatterns []string `yaml:"file_patterns,omitempty"`
	FuncPatterns []string `yaml:"func_patterns,omitempty"`
}
