package minusindex

func f1(a, b, c string, n int) {
}

func f2(a, b, c string) { // want "no int type arg at index -1 found for func f2"
}

func f3(a string) { // want "no int type arg at index -1 found for func f3"
}
