package pattern

func f1(num int) { // want "no string type arg found for func f1"
}

func f2(num1, num2 int) { // want "no string type arg found for func f2"
}

func f3(str string) { // want "no int type arg found for func f3"
}

func f4(num int, str string) {
}

func errorFn(num int, str string, err error) {
}
