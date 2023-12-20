package pattern

func f1(num int) { // want "func f1 not found arg string"
}

func f2(num1, num2 int) { // want "func f2 not found arg string"
}

func f3(str string) {
}

func f4(num int, str string) {
}
