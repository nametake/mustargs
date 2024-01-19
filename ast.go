package mustargs

import (
	"go/types"
)

type option func(*astArg)

type astArg struct {
	Index   int
	Type    string
	Pkg     string
	IsPtr   bool
	IsArray bool
}

func newAstArgs(signature *types.Signature) []*astArg {
	var args []*astArg
	for i := 0; i < signature.Params().Len(); i++ {
		opts := []option{withIndex(i)}
		argType := signature.Params().At(i).Type()

		if typ, ok := arrayType(argType); ok {
			opts = append(opts, withIsArray())
			argType = typ
		}

		if typ, ok := pointerType(argType); ok {
			opts = append(opts, withIsPtr())
			argType = typ
		}

		switch u := argType.(type) {
		case *types.Named:
			opts = append(opts, withType(u.Obj().Name()))
			// for primitive interface (e.g. error)
			if u.Obj().Pkg() != nil {
				opts = append(opts, withPkg(u.Obj().Pkg().Path()))
			}
		case *types.Basic:
			opts = append(opts, withType(u.Name()))
		default:
			continue
		}
		args = append(args, newAstArg(opts...))
	}
	return args
}

func newAstArg(options ...option) *astArg {
	astArg := &astArg{}
	for _, option := range options {
		option(astArg)
	}
	return astArg
}

func withType(typ string) option {
	return func(arg *astArg) {
		arg.Type = typ
	}
}

func withIndex(index int) option {
	return func(arg *astArg) {
		arg.Index = index
	}
}

func withIsPtr() option {
	return func(arg *astArg) {
		arg.IsPtr = true
	}
}

func withPkg(pkg string) option {
	return func(arg *astArg) {
		arg.Pkg = pkg
	}
}

func withIsArray() option {
	return func(arg *astArg) {
		arg.IsArray = true
	}
}

func pointerType(typ types.Type) (types.Type, bool) {
	if typ, ok := typ.Underlying().(*types.Pointer); ok {
		return typ.Elem(), true
	}
	return typ, false
}

func arrayType(typ types.Type) (types.Type, bool) {
	switch t := typ.(type) {
	case *types.Array:
		return t.Elem().Underlying(), true
	case *types.Slice:
		return t.Elem().Underlying(), true
	}
	return typ, false
}

func recvName(sig *types.Signature) string {
	if sig == nil {
		return ""
	}
	recv := sig.Recv()
	if recv == nil {
		return ""
	}
	recvType := recv.Type()

	switch typ := recvType.(type) {
	case *types.Pointer:
		if named, ok := typ.Elem().(*types.Named); ok {
			return named.Obj().Name()
		}
	case *types.Named:
		return typ.Obj().Name()
	}
	return ""
}
