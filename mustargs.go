package mustargs

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
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

type AstArg struct {
	Index   int
	Type    string
	Ptr     bool
	Pkg     string
	PkgName string
}

type ImportPackage struct {
	Pkg     string
	PkgName string
}

func ParseAst(expr ast.Expr, index int, ptr bool) *AstArg {
	switch typ := expr.(type) {
	case *ast.Ident:
		return &AstArg{
			Index: index,
			Type:  typ.Name,
			Ptr:   ptr,
		}
	case *ast.SelectorExpr:
		// TODO support pkg
		return &AstArg{
			Index: index,
			Type:  typ.Sel.Name,
			Ptr:   ptr,
		}
	}
	return nil
}

func ParseFunc(funcDecl *ast.FuncDecl) []*AstArg {
	var args []*AstArg
	for i, list := range funcDecl.Type.Params.List {
		for j := range list.Names {
			switch typ := list.Type.(type) {
			case *ast.StarExpr:
				args = append(args, ParseAst(typ.X, i+j, true))
			default:
				args = append(args, ParseAst(typ, i+j, false))
			}
		}
	}
	return args
}

func trimQuotes(str string) string {
	replacer := strings.NewReplacer("\"", "", "'", "")
	return replacer.Replace(str)
}

func ParseImport(specs []ast.Spec) map[string]*ImportPackage {
	packages := make(map[string]*ImportPackage)
	for _, spec := range specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			pkg := trimQuotes(s.Path.Value)
			name := ""
			if s.Name != nil {
				name = s.Name.Name
			} else {
				name = extractPkgName(pkg)
			}
			packages[name] = &ImportPackage{
				Pkg:     pkg,
				PkgName: name,
			}
		}
	}
	return packages
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		var packages map[string]*ImportPackage
		switch n := n.(type) {
		case *ast.GenDecl:
			packages = ParseImport(n.Specs)
		case *ast.FuncDecl:
			args := ParseFunc(n)
			for _, rule := range config.Rules {
				failedRuleArgs := rule.Args.Match(args)
				if len(failedRuleArgs) != 0 {
					pass.Reportf(n.Pos(), failedRuleArgs.ErrorMsg(n.Name.Name))
				}
			}
		}
		for name, pkg := range packages {
			fmt.Printf("%s: %#v\n", name, pkg)
		}
	})

	return nil, nil
}
