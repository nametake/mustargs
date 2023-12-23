package multitype

func f1(a string, b int) {
}

func f2(a int, b string) { // want "no string type arg found for func f2" "no int type arg found for func f2"
}
