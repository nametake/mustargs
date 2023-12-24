package mustargs

import (
	"fmt"
	"os"
	"regexp"
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

func (rule *Rule) TargetFile(filename string) (bool, error) {
	if len(rule.FilePatterns) == 0 {
		return true, nil
	}
	for _, pattern := range rule.FilePatterns {
		matched, err := regexp.MatchString(pattern, filename)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func (rule *Rule) TargetFunc(funcName string) (bool, error) {
	if len(rule.FuncPatterns) == 0 {
		return true, nil
	}
	for _, pattern := range rule.FuncPatterns {
		matched, err := regexp.MatchString(pattern, funcName)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
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
		pkgName := ""
		if rule.PkgName != "" {
			pkgName = rule.PkgName + "."
		}
		msg := fmt.Sprintf("no %s%s%s type arg", ptr, pkgName, rule.Type)
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
		if rule.MatchType(arg) && rule.MatchIndex(arg) && rule.MatchPtr(arg) && rule.MatchPkg(arg) {
			return true
		}
	}
	return false
}

func (rule *ArgRule) MatchType(arg *AstArg) bool {
	return arg.Type == rule.Type
}

func (rule *ArgRule) MatchIndex(arg *AstArg) bool {
	if rule.Index != nil {
		return arg.Index == *rule.Index
	}
	return true
}

func (rule *ArgRule) MatchPtr(arg *AstArg) bool {
	return arg.Ptr == rule.Ptr
}

func (rule *ArgRule) MatchPkg(arg *AstArg) bool {
	if rule.Pkg != "" {
		return arg.Pkg == rule.Pkg && rule.PkgName == arg.PkgName
	}
	return true
}
