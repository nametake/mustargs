package mustargs

import "go/ast"

type Config struct {
	Rules []*Rule `yaml:"rules"`
}

type Rule struct {
	Args         RuleArgs `yaml:"args"`
	FilePatterns []string `yaml:"file_patterns,omitempty"`
	FuncPatterns []string `yaml:"func_patterns,omitempty"`
	RecvPatterns []string `yaml:"recv_patterns,omitempty"`
}

type RuleArg struct {
	Type  string  `yaml:"type"`
	Index *int    `yaml:"index,omitempty"`
	Pkg   *string `yaml:"pkg,omitempty"`
	Name  *string `yaml:"name,omitempty"`
	Ptr   *bool   `yaml:"ptr,omitempty"`
}

type RuleArgs []*RuleArg

func (ruleArgs RuleArgs) Match(args []*ast.Ident) RuleArgs {
	unmatchRuleArgs := make(RuleArgs, 0, len(ruleArgs))
	for _, ruleArg := range ruleArgs {
		if !ruleArg.Match(args) {
			unmatchRuleArgs = append(unmatchRuleArgs, ruleArg)
		}
	}
	return unmatchRuleArgs
}

func (r *RuleArg) Match(args []*ast.Ident) bool {
	for i, arg := range args {
		if r.Index != nil {
			if i == *r.Index && arg.Name == r.Type {
				return true
			}
		} else {
			if arg.Name == r.Type {
				return true
			}
		}
	}
	return false
}
