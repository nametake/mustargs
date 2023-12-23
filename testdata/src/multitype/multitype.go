package multitype

func f1(a string, b int) {
}

func f2(a int, b string) { // want "func f2 not found arg int" "func f2 not found arg string"
}
