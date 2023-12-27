package mustargs

import (
	"fmt"
	"go/ast"
	"strings"
)

func NewAstArgs(funcDecl *ast.FuncDecl, packages map[string]string) []*AstArg {
	var args []*AstArg
	for i, list := range funcDecl.Type.Params.List {
		for j := range list.Names {
			args = append(args, checkLiteralExpr(list.Type, WithIndex(i+j), WithPkg(packages)))
		}
	}
	return args
}

func checkLiteralExpr(expr ast.Expr, options ...Option) *AstArg {
	switch typ := expr.(type) {
	case *ast.MapType, *ast.Ellipsis, *ast.InterfaceType, *ast.ChanType, *ast.FuncType, *ast.StructType:
		// TODO support
	case *ast.ArrayType:
		return checkStarExpr(typ.Elt, append(options, WithIsArray())...)
	default:
		return checkStarExpr(typ, options...)
	}
	panic(fmt.Sprintf("unsupported arg type: ast.Expr type = %T", expr))
}

func checkStarExpr(expr ast.Expr, options ...Option) *AstArg {
	switch typ := expr.(type) {
	case *ast.StarExpr:
		return checkSelectorExpr(typ.X, append(options, WithPtr())...)
	}
	return checkSelectorExpr(expr, options...)
}

func checkSelectorExpr(expr ast.Expr, options ...Option) *AstArg {
	switch typ := expr.(type) {
	case *ast.Ident:
		return NewAstArg(typ.Name, "", options...)
	case *ast.SelectorExpr:
		name := typ.X.(*ast.Ident).Name
		return NewAstArg(typ.Sel.Name, name, options...)
	}
	panic(fmt.Sprintf("unsupported arg type: ast.Expr type = %T", expr))
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
