package cache

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Item struct {
	sync.RWMutex

	key  interface{}
	data interface{}
	// 缓存时间
	duration time.Duration
	// 缓存创建时间
	createTime time.Time
	// 缓存最后使用时间
	accessTime time.Time
	// 缓存访问次数
	count int64
	// 删除时的回调函数
	deleteCallback func(key interface{})
}

func NewItem(key interface{}, duration time.Duration, data interface{}) *Item {
	t := time.Now()
	return &Item{
		key:            key,
		duration:       duration,
		createTime:     t,
		accessTime:     t,
		count:          0,
		deleteCallback: nil,
		data:           data,
	}
}

// 缓存保活
func (item *Item) KeepAlive() {
	item.Lock()
	defer item.Unlock()

	item.accessTime = time.Now()
	item.count++
}

func (item *Item) Duration() time.Duration {
	return item.duration
}

func (item *Item) AccessTime() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.accessTime
}

func (item *Item) CreateTime() time.Time {
	return item.createTime
}

func (item *Item) Count() int64 {
	item.RLock()
	defer item.RUnlock()
	return item.count
}

func (item *Item) Key() interface{} {
	return item.key
}

func (item *Item) Data() interface{} {
	return item.data
}

func (item *Item) SetDeleteCallback(f func(interface{})) {
	item.Lock()
	defer item.Unlock()
	item.deleteCallback = f
}

// 管理缓存
type Table struct {
	sync.RWMutex
	items map[interface{}]*Item
	// 清除缓存的触发器
	cleanUpTimer *time.Timer
	// 缓存的保存时间
	cleanUpInterval time.Duration
	// 获取缓存数据
	loadData func(key interface{}, args ...interface{}) *Item

	// 添加缓存后的回调函数
	addedCallback func(item *Item)
	// 删除缓存后的回调函数
	deleteCallback func(item *Item)
}

func (table *Table) Count() int {
	table.RLock()
	defer table.RUnlock()
	return len(table.items)
}

func (table *Table) Foreach(trans func(key interface{}, item *Item)) {
	table.RLock()
	defer table.RUnlock()

	for k, v := range table.items {
		trans(k, v)
	}
}

func (table *Table) SetDataLoader(f func(key interface{}, args ...interface{}) *Item) {
	table.Lock()
	defer table.Unlock()
	table.loadData = f
}

func (table *Table) SetAddCallback(f func(item *Item)) {
	table.Lock()
	defer table.Unlock()
	table.addedCallback = f
}

func (table *Table) SetDeleteCallback(f func(item *Item)) {
	table.Lock()
	defer table.Unlock()
	table.deleteCallback = f
}

func (table *Table) RunWithRecovery(f func()) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("occur error %v \r\n", e)
		}
	}()

	f()
}

// 定时清除缓存中的过期项
func (table *Table) checkExpire() {
	table.Lock()
	if table.cleanUpTimer != nil {
		// 清除过期缓存时，先暂停计时器
		table.cleanUpTimer.Stop()
	}

	if table.cleanUpInterval > 0 {
		table.log("Expiration check triggered after %v for table", table.cleanUpInterval)
	} else {
		table.log("Expiration check installed for table")
	}

	items := table.items
	table.Unlock()

	now := time.Now()
	// smallestDuration：所有缓存中最近将要过期的一个时间，将这个时间作为下一次触发清除动作的时间
	smallestDuration := 0 * time.Second
	for key, item := range items {
		item.RLock()
		duration := item.duration
		accessTime := item.accessTime
		item.RUnlock()

		// 过期时间设置为0，表示该缓存永不过期
		if duration == 0 {
			continue
		}

		// 缓存在过期时间内都没有被使用过，则删除该缓存
		if now.Sub(accessTime) >= duration {
			_, e := table.Delete(key)
			if e != nil {
				table.log("occur error while deleting %v", e.Error())
			}
		} else {
			// 找到下一个缓存要过期的时间
			if smallestDuration == 0 || duration-now.Sub(accessTime) < smallestDuration {
				smallestDuration = duration - now.Sub(accessTime)
			}
		}
	}

	// 设置下一个清除缓存的定时器
	table.Lock()
	table.cleanUpInterval = smallestDuration
	if smallestDuration > 0 {
		table.cleanUpTimer = time.AfterFunc(smallestDuration, func() {
			go table.RunWithRecovery(table.checkExpire)
		})
	}
	table.Unlock()
}

func (table *Table) Add(key interface{}, duration time.Duration, data interface{}) *Item {
	item := NewItem(key, duration, data)

	table.Lock()
	table.log("Adding item with key %v and lifespan of %d to table", key, duration)
	table.items[key] = item

	expDuration := table.cleanUpInterval
	addedItem := table.addedCallback
	table.Unlock()

	if addedItem != nil {
		// 添加缓存后的回调函数
		addedItem(item)
	}

	// 更新缓存容器的刷新时间
	if duration > 0 && (expDuration == 0 || duration < expDuration) {
		table.checkExpire()
	}

	return item
}

func (table *Table) Delete(key interface{}) (*Item, error) {
	table.RLock()
	r, ok := table.items[key]
	if !ok {
		table.RUnlock()
		return nil, errors.New(fmt.Sprintf("no item with key %s", key))
	}

	deleteCallback := table.deleteCallback
	table.RUnlock()

	if deleteCallback != nil {
		deleteCallback(r)
	}

	r.RLock()
	defer r.RUnlock()
	if r.deleteCallback != nil {
		r.deleteCallback(key)
	}

	table.Lock()
	defer table.Unlock()
	table.log("Deleting item with key %v created on %s and hit %d times from table", key, r.createTime, r.count)
	delete(table.items, key)
	return r, nil
}

// 检查缓存是否存在
func (table *Table) Exists(key interface{}) bool {
	table.RLock()
	defer table.RUnlock()
	_, ok := table.items[key]
	return ok
}

// 如果key存在，返回false，如果key不存在，添加该缓存，并返回true
func (table *Table) NotFoundAdd(key interface{}, lifeSpan time.Duration, data interface{}) bool {
	table.Lock()

	if _, ok := table.items[key]; ok {
		return false
	}

	item := NewItem(key, lifeSpan, data)
	table.log("Adding item key %v and lifespan of %d to table", key, lifeSpan)
	table.items[key] = item

	expDur := table.cleanUpInterval
	addedItem := table.addedCallback
	table.Unlock()

	if addedItem != nil {
		addedItem(item)
	}

	if lifeSpan > 0 && (expDur == 0 || lifeSpan < expDur) {
		table.checkExpire()
	}
	return true
}

func (table *Table) Value(key interface{}, args ...interface{}) (*Item, error) {
	table.RLock()
	r, ok := table.items[key]
	loadData := table.loadData
	table.RUnlock()
	if ok {
		r.KeepAlive()
		return r, nil
	}

	if loadData != nil {
		item := loadData(key, args...)
		if item != nil {
			table.Add(key, item.duration, item.data)
			return item, nil
		}

		return nil, errors.New("cannot load item")
	}

	return nil, nil
}

func (table *Table) Truncate() {
	table.Lock()
	defer table.Unlock()

	table.items = make(map[interface{}]*Item)
	table.cleanUpInterval = 0
	if table.cleanUpTimer != nil {
		table.cleanUpTimer.Stop()
	}
}

type ItemPair struct {
	Key         interface{}
	AccessCount int64
}

type ItemPairList []ItemPair

func (p ItemPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ItemPairList) Len() int           { return len(p) }
func (p ItemPairList) Less(i, j int) bool { return p[i].AccessCount > p[j].AccessCount }

// 找到前count位最多访问的缓存
func (table *Table) MostAccessed(count int64) []*Item {
	table.RLock()
	defer table.RUnlock()

	p := make(ItemPairList, len(table.items))
	i := 0
	for k, v := range table.items {
		p[i] = ItemPair{k, v.count}
		i++
	}

	sort.Sort(p)

	var r []*Item
	c := int64(0)
	for _, v := range p {
		if c > count {
			break
		}
		item, ok := table.items[v.Key]
		if ok {
			r = append(r, item)
		}
		c++
	}
	return r
}

func (table *Table) log(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v)
}
