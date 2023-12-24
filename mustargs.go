package mustargs

import (
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

func ExtractAstArg(expr ast.Expr, index int, ptr bool, packages map[string]string) *AstArg {
	switch typ := expr.(type) {
	case *ast.Ident:
		return &AstArg{
			Index:   index,
			Type:    typ.Name,
			Ptr:     ptr,
			Pkg:     "",
			PkgName: "",
		}
	case *ast.SelectorExpr:
		name := typ.X.(*ast.Ident).Name
		return &AstArg{
			Index:   index,
			Type:    typ.Sel.Name,
			Ptr:     ptr,
			Pkg:     packages[name],
			PkgName: name,
		}
	}
	return nil
}

func ExtractAstArgs(funcDecl *ast.FuncDecl, packages map[string]string) []*AstArg {
	var args []*AstArg
	for i, list := range funcDecl.Type.Params.List {
		for j := range list.Names {
			switch typ := list.Type.(type) {
			case *ast.StarExpr:
				args = append(args, ExtractAstArg(typ.X, i+j, true, packages))
			default:
				args = append(args, ExtractAstArg(typ, i+j, false, packages))
			}
		}
	}
	return args
}

func trimQuotes(str string) string {
	replacer := strings.NewReplacer("\"", "", "'", "")
	return replacer.Replace(str)
}

func ExtractImportPackages(specs []ast.Spec) map[string]string {
	packages := make(map[string]string)
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
			packages[name] = pkg
		}
	}
	return packages
}

func ExtractRecvName(recv *ast.FieldList) string {
	if recv == nil {
		return ""
	}
	switch typ := recv.List[0].Type.(type) {
	case *ast.Ident:
		return typ.Name
	case *ast.StarExpr:
		return typ.X.(*ast.Ident).Name
	}
	return ""
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

	var packages map[string]string
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fileName := pass.Fset.File(n.Pos()).Name()

		switch n := n.(type) {
		case *ast.GenDecl:
			packages = ExtractImportPackages(n.Specs)
		case *ast.FuncDecl:
			funcName := n.Name.Name
			if funcName == "init" || funcName == "main" {
				return
			}
			recvName := ExtractRecvName(n.Recv)
			args := ExtractAstArgs(n, packages)
			for _, rule := range config.Rules {
				isTargetFile, err := rule.IsTargetFile(fileName)
				if err != nil {
					pass.Reportf(n.Pos(), err.Error())
					return
				}
				if !isTargetFile {
					continue
				}

				isTargetFunc, err := rule.IsTargetFunc(funcName)
				if err != nil {
					pass.Reportf(n.Pos(), err.Error())
					return
				}
				if !isTargetFunc {
					continue
				}

				isTargetRecv, err := rule.IsTargetRecv(recvName)
				if err != nil {
					pass.Reportf(n.Pos(), err.Error())
				}
				if !isTargetRecv {
					continue
				}

				unmatchedRules := rule.Args.Match(args)
				if len(unmatchedRules) != 0 {
					pass.Reportf(n.Pos(), unmatchedRules.ErrorMsg(n.Name.Name))
				}
			}
		}
	})

	return nil, nil
}
