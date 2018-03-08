package domain

type Pod struct {
	Name  string
	Kind  string
}

type Pods *[]Pod

func NewPod() *Pod {
	return &Pod{}
}
