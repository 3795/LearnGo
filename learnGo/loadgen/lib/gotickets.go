package lib

import (
	"errors"
	"fmt"
)

// 表示Goroutine票池的接口
type GoTickets interface {
	// 取走一个资源
	Take()
	// 返还一个资源
	Return()
	// 资源池是否已被激活
	Active() bool
	// 资源总数
	Total() uint32
	// 剩余的资源数
	Remainder() uint32
}

type myGoTickets struct {
	total    uint32        // 资源总数
	ticketCh chan struct{} // 资源容器
	active   bool          // 资源池是否被激活
}

func NewGoTickets(total uint32) (GoTickets, error) {
	gt := myGoTickets{}
	if !gt.init(total) {
		errMsg :=
			fmt.Sprintf("The goroutine ticket pool can NOT be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return &gt, nil
}

func (gt *myGoTickets) init(total uint32) bool {
	if gt.active {
		return false
	}
	if total == 0 {
		return false
	}

	ch := make(chan struct{}, total)
	n := int(total)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
	}
	gt.ticketCh = ch
	gt.total = total
	gt.active = true
	return true
}

func (gt *myGoTickets) Take() {
	<-gt.ticketCh
}

func (gt *myGoTickets) Return() {
	gt.ticketCh <- struct{}{}
}

func (gt *myGoTickets) Active() bool {
	return gt.active
}

func (gt *myGoTickets) Total() uint32 {
	return gt.total
}

func (gt *myGoTickets) Remainder() uint32 {
	return uint32(len(gt.ticketCh))
}
