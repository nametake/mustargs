package unsupporttype

func f1(m map[interface{}]interface{}, f func() error, c chan int, s struct{}, i interface{}, e ...int) { // want "no string type arg found for func f1"
}
