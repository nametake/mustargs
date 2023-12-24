package pkgtypenopkg

import (
	"context"
	"io"
	stdos "os"
)

func f1(ctx context.Context, r io.Reader, f *stdos.File) {
}

func f2(ctx context.Context) { // want "no Reader type arg at index 1, no \\*File type arg at index 2 found for func f2"
}

func f3(ctx context.Context, w io.Writer, f *stdos.File) { // want "no Reader type arg at index 1 found for func f3"
}

func f4(ctx context.Context, r io.Reader, f stdos.File) { // want "no \\*File type arg at index 2 found for func f4"
}
