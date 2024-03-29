package mustargs

import (
	"fmt"
	"regexp"
	"strings"
)

type Rule struct {
	Args               ArgRules `yaml:"args"`
	FilePatterns       []string `yaml:"file_patterns,omitempty"`
	FuncPatterns       []string `yaml:"func_patterns,omitempty"`
	RecvPatterns       []string `yaml:"recv_patterns,omitempty"`
	IgnoreFilePatterns []string `yaml:"ignore_file_patterns,omitempty"`
	IgnoreFuncPatterns []string `yaml:"ignore_func_patterns,omitempty"`
	IgnoreRecvPatterns []string `yaml:"ignore_recv_patterns,omitempty"`
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

func (rule *Rule) IsIgnoreFile(fileName string) (bool, error) {
	if len(rule.IgnoreFilePatterns) == 0 {
		return false, nil
	}
	for _, pattern := range rule.IgnoreFilePatterns {
		matched, err := regexp.MatchString(pattern, fileName)
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

func (rule *Rule) IsIgnoreFunc(funcName string) (bool, error) {
	if len(rule.IgnoreFuncPatterns) == 0 {
		return false, nil
	}
	for _, pattern := range rule.IgnoreFuncPatterns {
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

func (rule *Rule) IsIgnoreRecv(recvName string) (bool, error) {
	if len(rule.IgnoreRecvPatterns) == 0 {
		return false, nil
	}
	for _, pattern := range rule.IgnoreRecvPatterns {
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
	IsPtr   bool   `yaml:"is_ptr,omitempty"`
	IsArray bool   `yaml:"is_array,omitempty"`
}

func (rule *ArgRule) Match(args []*astArg) bool {
	argsLen := len(args)
	for i, arg := range args {
		if rule.Index != nil && !rule.matchIndex(i, argsLen) {
			continue
		}

		if rule.matchType(arg) && rule.matchPtr(arg) && rule.matchPkg(arg) && rule.matchIsArray(arg) {
			return true
		}
	}

	return false
}
func (rule *ArgRule) matchIndex(index, argsLen int) bool {
	if *rule.Index >= 0 && index != *rule.Index {
		return false
	}
	if *rule.Index < 0 && (index-argsLen) != *rule.Index {
		return false
	}
	return true
}

func (rule *ArgRule) matchType(arg *astArg) bool {
	return arg.Type == rule.Type
}

func (rule *ArgRule) matchPtr(arg *astArg) bool {
	return arg.IsPtr == rule.IsPtr
}

func (rule *ArgRule) matchPkg(arg *astArg) bool {
	if rule.Pkg != "" {
		return arg.Pkg == rule.Pkg
	}
	return true
}

func (rule *ArgRule) matchIsArray(arg *astArg) bool {
	return arg.IsArray == rule.IsArray
}

type ArgRules []*ArgRule

func (argRules ArgRules) Match(args []*astArg) ArgRules {
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
		array := ""
		if rule.IsArray {
			array = "[]"
		}
		ptr := ""
		if rule.IsPtr {
			ptr = "*"
		}
		pkgName := ""
		if rule.Pkg != "" {
			pkgName = rule.Pkg + "."
		}
		msg := fmt.Sprintf("no %s%s%s%s type arg", array, ptr, pkgName, rule.Type)
		if rule.Index != nil {
			msg += fmt.Sprintf(" at index %d", *rule.Index)
		}
		ruleErrMsgs = append(ruleErrMsgs, msg)
	}

	return fmt.Sprintf("%s found for func %s", strings.Join(ruleErrMsgs, ", "), funcName)
}
