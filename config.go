package mustargs

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rules []*Rule `yaml:"rules"`
}

func extractPkgName(importPath string) string {
	parts := strings.Split(importPath, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func loadConfig(filepath string) (*Config, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var config *Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	for _, rule := range config.Rules {
		for _, arg := range rule.Args {
			if arg.PkgName == "" {
				arg.PkgName = extractPkgName(arg.Pkg)
			}
		}
	}

	return config, nil
}

type Rule struct {
	Args         ArgRules `yaml:"args"`
	FilePatterns []string `yaml:"file_patterns,omitempty"`
	FuncPatterns []string `yaml:"func_patterns,omitempty"`
	RecvPatterns []string `yaml:"recv_patterns,omitempty"`
}

type ArgRule struct {
	Type    string `yaml:"type"`
	Index   *int   `yaml:"index,omitempty"`
	Pkg     string `yaml:"pkg,omitempty"`
	PkgName string `yaml:"pkg_name,omitempty"`
	Ptr     bool   `yaml:"ptr,omitempty"`
}

type ArgRules []*ArgRule

func (argRules ArgRules) ErrorMsg(funcName string) string {
	ruleErrMsgs := make([]string, 0, len(argRules))
	for _, rule := range argRules {
		ptr := ""
		if rule.Ptr {
			ptr = "*"
		}
		msg := fmt.Sprintf("no %s%s type arg", ptr, rule.Type)
		if rule.Index != nil {
			msg += fmt.Sprintf(" at index %d", *rule.Index)
		}
		ruleErrMsgs = append(ruleErrMsgs, msg)
	}

	return fmt.Sprintf("%s found for func %s", strings.Join(ruleErrMsgs, ", "), funcName)
}

func (argRules ArgRules) Match(args []*AstArg) ArgRules {
	unmatchRuleArgs := make(ArgRules, 0, len(argRules))
	for _, ruleArg := range argRules {
		if !ruleArg.Match(args) {
			unmatchRuleArgs = append(unmatchRuleArgs, ruleArg)
		}
	}
	return unmatchRuleArgs
}

func (rule *ArgRule) Match(args []*AstArg) bool {
	for _, arg := range args {
		if rule.Index != nil {
			if arg.Index == *rule.Index && arg.Type == rule.Type && arg.Ptr == rule.Ptr {
				return true
			}
		} else {
			if arg.Type == rule.Type && arg.Ptr == rule.Ptr {
				return true
			}
		}
	}

	return false
}
