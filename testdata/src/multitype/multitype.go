package multitype

func f1(a string, b int) {
}

func f2(a int, b string) { // want "no string type arg at index 0, no int type arg at index 1 found for func f2"
}
