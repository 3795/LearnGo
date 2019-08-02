package toolkit

// 数据缓冲池接口
type pool interface {
	// 获取池中缓冲器的统一容量
	BufferCap() uint32
	// 获取池中缓冲器的最大数量
	MaxBufferNumber() uint32
	// 获取池中缓冲器的数量
	BufferNumber() uint32
	// 获取缓冲池中数据的总数
	Total() uint64
	// 向缓冲池中放入数据
	Put(datum interface{}) error
	// 从缓冲池中获取数据
	Get(datum interface{}, err error)
	// 关闭缓冲池
	Close() bool
	// 判断缓冲池是否已关闭
	Closed() bool
}

type myPool struct {
	// 缓冲器的容量
	bufferCap uint32
	// 缓冲器的最大数量
	maxBufferNumber uint32
	// 缓冲器的实际数量
	bufferNumber uint32
}
