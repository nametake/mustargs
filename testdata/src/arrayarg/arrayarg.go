package arrayarg

import (
	"os"
)

func f1(a []string, b []*os.File) {
}

func f2(a string, b *os.File) { // want "no \\[\\]string type arg at index 0, no \\[\\]\\*os.File type arg at index 1 found for func f2"
}
