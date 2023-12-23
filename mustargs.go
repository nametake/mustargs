package mustargs

import (
	"go/ast"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"gopkg.in/yaml.v3"
)

const doc = "mustargs is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "mustargs",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

var (
	configPath string
)

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", "", "config file path")
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
			if arg.Pkg != nil && arg.Name == nil {
				name := extractPkgName(*arg.Pkg)
				arg.Name = &name
			}
		}
	}

	return config, nil
}

type AstArg struct {
	Index int
	Type  string
	Name  string
}

func ParseFunc(funcDecl *ast.FuncDecl) []*AstArg {
	var args []*AstArg
	for _, li := range funcDecl.Type.Params.List {
		for index, name := range li.Names {
			switch typ := li.Type.(type) {
			case *ast.Ident:
				args = append(args, &AstArg{
					Index: index,
					Type:  typ.Name,
					Name:  name.Name,
				})
			case *ast.SelectorExpr:
				// TODO support pkg
				args = append(args, &AstArg{
					Index: index,
					Type:  typ.Sel.Name,
					Name:  name.Name,
				})
			}
		}
	}
	return args
}

func ParseFuncDecl(funcDecl *ast.FuncDecl) []*ast.Ident {
	var args []*ast.Ident
	for _, li := range funcDecl.Type.Params.List {
		for range li.Names {
			switch t := li.Type.(type) {
			case *ast.Ident:
				args = append(args, t)
			case *ast.SelectorExpr:
				// TODO support pkg
				args = append(args, t.Sel)
			}
		}
	}
	return args
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == "gopher" {
				pass.Reportf(n.Pos(), "identifier is gopher")
			}
		case *ast.FuncDecl:
			args := ParseFuncDecl(n)
			for _, rule := range config.Rules {
				// TODO pattern check
				failedRuleArgs := rule.Args.Match(args)
				for _, failedRuleArg := range failedRuleArgs {
					pass.Reportf(n.Pos(), "func %s not found arg %s", n.Name.Name, failedRuleArg.Type)
				}
			}
		}
	})

	return nil, nil
}
