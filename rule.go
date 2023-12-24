package mustargs

import (
	"fmt"
	"regexp"
	"strings"
)

type Rule struct {
	Args               ArgRules `yaml:"args"`
	FilePatterns       []string `yaml:"file_patterns,omitempty"`
	IgnoreFilePatterns []string `yaml:"ignore_file_patterns,omitempty"`
	FuncPatterns       []string `yaml:"func_patterns,omitempty"`
	RecvPatterns       []string `yaml:"recv_patterns,omitempty"`
}

func (rule *Rule) IsTargetFile(filename string) (bool, error) {
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

func (rule *Rule) IsIgnoreFile(filename string) (bool, error) {
	if len(rule.IgnoreFilePatterns) == 0 {
		return false, nil
	}
	for _, pattern := range rule.IgnoreFilePatterns {
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

func (rule *Rule) IsTargetFunc(funcName string) (bool, error) {
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

func (rule *Rule) IsTargetRecv(recvName string) (bool, error) {
	if len(rule.RecvPatterns) == 0 {
		return true, nil
	}
	for _, pattern := range rule.RecvPatterns {
		matched, err := regexp.MatchString(pattern, recvName)
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

func (rule *ArgRule) Match(args []*AstArg) bool {
	for _, arg := range args {
		if rule.matchType(arg) && rule.matchIndex(arg) && rule.matchPtr(arg) && rule.matchPkg(arg) {
			return true
		}
	}
	return false
}

func (rule *ArgRule) matchType(arg *AstArg) bool {
	return arg.Type == rule.Type
}

func (rule *ArgRule) matchIndex(arg *AstArg) bool {
	if rule.Index != nil {
		return arg.Index == *rule.Index
	}
	return true
}

func (rule *ArgRule) matchPtr(arg *AstArg) bool {
	return arg.Ptr == rule.Ptr
}

func (rule *ArgRule) matchPkg(arg *AstArg) bool {
	if rule.Pkg != "" {
		return arg.Pkg == rule.Pkg && rule.PkgName == arg.PkgName
	}
	return true
}

type ArgRules []*ArgRule

func (argRules ArgRules) Match(args []*AstArg) ArgRules {
	unmatchedRules := make(ArgRules, 0, len(argRules))
	for _, ruleArg := range argRules {
		if !ruleArg.Match(args) {
			unmatchedRules = append(unmatchedRules, ruleArg)
		}
	}
	return unmatchedRules
}

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
