package module

import (
	"math"
	"sync"
)

type SNGenerator interface {
	// 获取预设的最小序列号
	Start() uint64
	// 获取预设的最大序列号
	Max() uint64
	// 获取下一个序列号
	Next() uint64
	// 获取循环计数
	CycleCount() uint64
	// 获取一个序列号并准备下一个序列号
	Get() uint64
}

type mySNGenerator struct {
	// 序列号的最小值
	start uint64
	// 序列号的最大值
	max uint64
	// 下一个序列号
	next uint64
	// 循环的计数
	cycleCount uint64
	// 读写锁
	lock sync.RWMutex
}

func NewSNGenerator(start, max uint64) SNGenerator {
	if max == 0 {
		max = math.MaxUint64
	}
	return &mySNGenerator{start: start, max: max}
}

func (gen *mySNGenerator) Start() uint64 {
	return gen.start
}

func (gen *mySNGenerator) Max() uint64 {
	return gen.max
}

func (gen *mySNGenerator) Next() uint64 {
	gen.lock.RLock()
	defer gen.lock.RUnlock()
	return gen.next
}

func (gen *mySNGenerator) CycleCount() uint64 {
	gen.lock.RLock()
	defer gen.lock.RUnlock()
	return gen.cycleCount
}

func (gen *mySNGenerator) Get() uint64 {
	gen.lock.Lock()
	defer gen.lock.Unlock()
	id := gen.next
	if id == gen.max {
		gen.next = gen.start
		gen.cycleCount++
	} else {
		gen.next++
	}
	return id
}
