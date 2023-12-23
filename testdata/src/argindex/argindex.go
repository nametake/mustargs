package argindex

func f1(num int) { // want "no string type arg found for func f1"
}

func f2(num1, num2 int) { // want "no string type arg found for func f2"
}

func f3(str string) { // want "no string type arg found for func f3"
}

func f4(str1, str2 string) {
}

func f5(num int, str string) {
}

func f6(num int, str string, b bool) {
}
