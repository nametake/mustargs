package mustargs

import (
	"fmt"
	"go/ast"
	"strings"
)

type Config struct {
	Rules []*Rule `yaml:"rules"`
}

type Rule struct {
	Args         ArgRules `yaml:"args"`
	FilePatterns []string `yaml:"file_patterns,omitempty"`
	FuncPatterns []string `yaml:"func_patterns,omitempty"`
	RecvPatterns []string `yaml:"recv_patterns,omitempty"`
}

type ArgRule struct {
	Type  string  `yaml:"type"`
	Index *int    `yaml:"index,omitempty"`
	Pkg   *string `yaml:"pkg,omitempty"`
	Name  *string `yaml:"name,omitempty"`
	Ptr   *bool   `yaml:"ptr,omitempty"`
}

type ArgRules []*ArgRule

func (argRules ArgRules) ErrorMsg(funcName string) string {
	ruleErrMsgs := make([]string, 0, len(argRules))
	for _, rule := range argRules {
		msg := fmt.Sprintf("no %s type arg", rule.Type)
		if rule.Index != nil {
			msg += fmt.Sprintf(" at index %d", *rule.Index)
		}
		ruleErrMsgs = append(ruleErrMsgs, msg)
	}

	return fmt.Sprintf("%s found for func %s", strings.Join(ruleErrMsgs, ", "), funcName)
}

func (argRules ArgRules) Match(args []*ast.Ident) ArgRules {
	unmatchRuleArgs := make(ArgRules, 0, len(argRules))
	for _, ruleArg := range argRules {
		if !ruleArg.Match(args) {
			unmatchRuleArgs = append(unmatchRuleArgs, ruleArg)
		}
	}
	return unmatchRuleArgs
}

func (rule *ArgRule) Match(args []*ast.Ident) bool {
	for i, arg := range args {
		if rule.Index != nil {
			if i == *rule.Index && arg.Name == rule.Type {
				return true
			}
		} else {
			if arg.Name == rule.Type {
				return true
			}
		}
	}
	return false
}
