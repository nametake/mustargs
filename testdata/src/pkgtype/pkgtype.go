package pkgtype

import (
	"context"
	"io"
)

// TODO: support package name
func f1(ctx context.Context, w io.Writer) { // want "func f1 not found arg Reader"
}

func f2(ctx context.Context, r io.Reader) {
}
