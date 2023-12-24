package recvpattern

type TargetStruct struct{}

func (s *TargetStruct) TargetMethod1(num int) { // want "no string type arg at index 1 found for func TargetMethod1"
}

func (TargetStruct) TargetMethod2(num int) { // want "no string type arg at index 1 found for func TargetMethod2"
}

type Struct struct{}

func (s *Struct) Method1(num int) {
}

type NoTargetStruct struct{}

func (s *NoTargetStruct) NoTargetMethod1(num int) {
}
