package mustargs

import (
	"fmt"
	"go/ast"
	"os"

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

func loadConfig(filepath string) (*Config, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var config *Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return config, nil
}

type AstArg struct {
	Type *ast.Ident
	Name *ast.Ident
}

func ParseFuncDecl(funcDecl *ast.FuncDecl) []*ast.Ident {
	var args []*ast.Ident
	for _, li := range funcDecl.Type.Params.List {
		for range li.Names {
			args = append(args, li.Type.(*ast.Ident))
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
	fmt.Printf("%+v\n", config)

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
				for _, arg := range rule.Args {
					for _, a := range args {
						if a.Name == arg.Type {
							return
						}
					}
					pass.Reportf(n.Pos(), "func %s not found arg %s", n.Name.Name, arg.Type)
				}
			}
		}

	})

	return nil, nil
}
