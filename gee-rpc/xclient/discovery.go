package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect     SelectMode = iota // 随机策略
	RoundRobinSelect                   // 轮询策略
)

// 服务发现功能的接口
type Discovery interface {
	Refresh() error                      // 从注册中心更新服务列表
	Update(servers []string) error       // 手动更新服务列表
	Get(mode SelectMode) (string, error) // 根据负载均衡策略，选择一个服务实例
	GetAll() ([]string, error)           // 返回所有的服务实例
}

// 存储服务数据的结构，该结构存放于客户端内部，客户端自己维护
type MultiServerDiscovery struct {
	r       *rand.Rand    // 生成一个随机数
	mu      *sync.RWMutex // 读写锁，保证并发写入正确
	servers []string      // 服务器列表
	index   int           // 记录轮询算法已经轮询到的位置，可以避免每次都从0开始，初始化时随机指定一个值
}

func (d *MultiServerDiscovery) Refresh() error {
	return nil
}

func (d *MultiServerDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	return nil
}

func (d *MultiServerDiscovery) Get(mode SelectMode) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	n := len(d.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available servers")
	}
	switch mode {
	case RandomSelect: // 如果是随机模式的话，就随机返回一个
		return d.servers[d.r.Intn(n)], nil
	case RoundRobinSelect:
		// 此处取模是为了防止数组越界，index在上一次计算后，值为5，但是servers列表更新过了一次，少了两个server，长度为3
		// 不取模运算的话，就会出问题
		s := d.servers[d.index%n]
		d.index = (d.index + 1) % n // 重新安全计算index的值
		return s, nil
	default:
		return "", errors.New("rpc discovery: not supported select mode")
	}
}

func (d *MultiServerDiscovery) GetAll() ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}

func NewMultiServerDiscovery(servers []string) *MultiServerDiscovery {
	d := &MultiServerDiscovery{
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
		servers: servers,
		mu:      &sync.RWMutex{},
	}
	d.index = d.r.Intn(math.MaxInt32 - 1)
	return d
}
