package stub

import "Project/LearnGo/Learn/webcrawler/module"

// 组件内部基础接口类型
type ModuleInternal interface {
	module.Module
	// 将组件调用次数加1
	IncrCalledCount()
	// 将组件接口调用次数加1
	IncrAcceptedCount()
	// 将组件完成调用次数加1
	IncrCompletedCount()
	// 将组件实时处理数加1
	IncrHandlingNumber()
	// 将组件实时处理数减1
	DecrHandlingNumber()
	// 清空所有计数
	Clear()
}
