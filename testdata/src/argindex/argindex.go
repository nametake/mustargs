package argindex

func f1(num int) { // want "func f1 not found arg string"
}

func f2(num1, num2 int) { // want "func f2 not found arg string"
}

func f3(str string) { // want "func f3 not found arg string"
}

func f4(str1, str2 string) {
}

func f5(num int, str string) {
}

func f6(num int, str string, b bool) {
}