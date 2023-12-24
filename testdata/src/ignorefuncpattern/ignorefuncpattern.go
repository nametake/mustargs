package ignorefuncpattern

func IgnoreFunc1(num int) {
}

func Func1(num int) { // want "no string type arg at index 1 found for func Func1"
}
