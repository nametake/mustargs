package mustargs

import (
	"go/ast"
	"go/importer"
	"go/types"

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

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	nodeFilter := []ast.Node{
		// (*ast.GenDecl)(nil),
		(*ast.FuncDecl)(nil),
		// (*ast.Ident)(nil),
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}

	if _, err := conf.Check(pass.Pkg.Path(), pass.Fset, pass.Files, info); err != nil {
		return nil, err
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			def, ok := info.Defs[n.Name]
			if !ok {
				return
			}

			signature, ok := def.Type().(*types.Signature)
			if !ok {
				return
			}
			fileName := pass.Fset.File(n.Pos()).Name()

			funcName := def.Name()
			if funcName == "init" || funcName == "main" {
				return
			}

			recvName := recvName(signature)

			args := newAstArgsBySignature(signature)

			for _, rule := range config.Rules {
				isTargetFile, err := rule.IsTargetFile(fileName)
				if err != nil {
					pass.Reportf(n.Pos(), err.Error())
					return
				}
				if !isTargetFile {
					continue
				}

				isIgnoreFile, err := rule.IsIgnoreFile(fileName)
				if err != nil {
					pass.Reportf(n.Pos(), err.Error())
					return
				}
				if isIgnoreFile {
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

				isIgnoreFunc, err := rule.IsIgnoreFunc(funcName)
				if err != nil {
					pass.Reportf(n.Pos(), err.Error())
					return
				}
				if isIgnoreFunc {
					continue
				}

				isTargetRecv, err := rule.IsTargetRecv(recvName)
				if err != nil {
					pass.Reportf(n.Pos(), err.Error())
				}
				if !isTargetRecv {
					continue
				}

				isIgnoreRecv, err := rule.IsIgnoreRecv(recvName)
				if err != nil {
					pass.Reportf(n.Pos(), err.Error())
					return
				}
				if isIgnoreRecv {
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
