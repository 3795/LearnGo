package core

// 恢复系统中的Panic
func RunWithRecovery(f func()) {
	defer func() {
		if err := recover(); err != nil {
			LOGGER.Error("error in async method: %v", err)
		}
	}()

	f()
}

// 快速失败
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
