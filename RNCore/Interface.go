package RNCore

type IName interface {
	Name() string
	Type_Name() string
}

type IRun interface {
	Run()
}
