package core

type Bean interface {
	// 构造时创建bean
	Init()

	Cleanup()

	Bootstrap()

	PanicError(err error)
}
