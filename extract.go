package mustargs

import (
	"go/ast"
	"strings"
)

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

func extractPkgName(importPath string) string {
	parts := strings.Split(importPath, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
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
