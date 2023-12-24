package pointerarg

func f1(a *string) {
}

func f2(a string) { // want "no \\*string type arg found for func f1"
}
