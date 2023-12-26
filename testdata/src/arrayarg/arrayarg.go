package arrayarg

import (
	"os"
)

func f1(a []string, b []*os.File) {
}

func f2(a string, b *os.File) { // want "no \\[\\]string type arg, no \\[\\]\\*os.File type arg found for func f2"
}
