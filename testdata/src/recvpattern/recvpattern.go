package recvpattern

type TargetStruct struct{}

func (s *Struct) TargetMethod1(num int) { // want "no string type arg at index 1 found for func TargetMethod1"
}

type Struct struct{}

func (s *Struct) Method1(num int) {
}

type NoTargetStruct struct{}

func (s *NoTargetStruct) NoTargetMethod1(num int) {
}
