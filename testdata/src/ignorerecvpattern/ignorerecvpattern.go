package ignorerecvpattern

type Struct struct{}

func (s *Struct) Method(num int) { // want "no string type arg at index 1 found for func Method"
}

type IgnoreStruct struct{}

func (s *IgnoreStruct) Method(num int) {
}
