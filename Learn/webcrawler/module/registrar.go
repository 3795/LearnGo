package module

import (
	"Project/LearnGo/Learn/webcrawler/errors"
	er "errors"
	"fmt"
	"sync"
)

var ErrNotFoundModuleInstance = er.New("not found module instance")

// 组件注册的接口
type Registrar interface {
	// 组件实例注册
	Register(module Module) (bool, error)
	// 注销组件实例
	UnRegister(mid MID) (bool, error)
	// 获取一个指定类型的组件的实例，应基于负载均衡策略返回实例
	Get(moduleType Type) (Module, error)
	// 返回指定类型的所有组件实例
	GetAllByType(moduleType Type) (map[MID]Module, error)
	// 获取所有组件实例
	GetAll() map[MID]Module
	// 清楚所有的组件注册记录
	Clear()
}

type myRegistrar struct {
	// 组件类型与对应组件实例之间的映射
	moduleTypeMap map[Type]map[MID]Module
	// 读写锁
	rwLock sync.RWMutex
}

func NewRegistrar() Registrar {
	return &myRegistrar{
		moduleTypeMap: map[Type]map[MID]Module{},
	}
}

func (registrar *myRegistrar) Register(module Module) (bool, error) {
	if module == nil {
		return false, errors.NewIllegalParameterError("nil module instance")
	}
	mid := module.ID()
	parts, err := SplitMID(mid)
	if err != nil {
		return false, err
	}
	moduleType := legalLetterTypeMap[parts[0]]
	if !CheckType(moduleType, module) {
		errMsg := fmt.Sprintf("incorrect module type: %s", moduleType)
		return false, errors.NewIllegalParameterError(errMsg)
	}

	registrar.rwLock.Lock()
	defer registrar.rwLock.Unlock()

	modules := registrar.moduleTypeMap[moduleType]
	// 如果还没有该类型的组件注册，则初始化容器
	if modules == nil {
		modules = map[MID]Module{}
	}
	// 如果该组件已经注册过，则返回false
	if _, ok := modules[mid]; ok {
		return false, nil
	}
	// 将组件注册进容器，并更新该容器
	modules[mid] = module
	registrar.moduleTypeMap[moduleType] = modules
	return true, nil
}

func (registrar *myRegistrar) UnRegister(mid MID) (bool, error) {
	parts, err := SplitMID(mid)
	if err != nil {
		return false, err
	}
	moduleType := legalLetterTypeMap[parts[0]]
	var deleted bool

	registrar.rwLock.Lock()
	defer registrar.rwLock.Unlock()

	if modules, ok := registrar.moduleTypeMap[moduleType]; ok {
		if _, ok := modules[mid]; ok {
			delete(modules, mid)
			deleted = true
		}
	}
	return deleted, nil
}

// 获取一个指定类型的组件的实例
// 基于负载均衡策略返回实例，选择得分最小的，即负载最小的组件
func (registrar *myRegistrar) Get(moduleType Type) (Module, error) {
	modules, err := registrar.GetAllByType(moduleType)
	if err != nil {
		return nil, err
	}
	minScore := uint64(0)
	var selectedModule Module
	for _, module := range modules {
		SetScore(module)
		score := module.Score()
		if minScore == 0 || score < minScore {
			selectedModule = module
			minScore = score
		}
	}
	return selectedModule, nil
}

func (registrar *myRegistrar) GetAllByType(moduleType Type) (map[MID]Module, error) {
	if !LegalType(moduleType) {
		errMsg := fmt.Sprintf("illagal module type: %s", moduleType)
		return nil, errors.NewIllegalParameterError(errMsg)
	}

	registrar.rwLock.RLock()
	defer registrar.rwLock.RUnlock()

	modules := registrar.moduleTypeMap[moduleType]
	if len(modules) == 0 {
		return nil, ErrNotFoundModuleInstance
	}
	result := map[MID]Module{}
	for mid, module := range modules {
		result[mid] = module
	}
	return result, nil
}

func (registrar *myRegistrar) GetAll() map[MID]Module {
	result := map[MID]Module{}

	registrar.rwLock.RLock()
	defer registrar.rwLock.RUnlock()

	for _, modules := range registrar.moduleTypeMap {
		for mid, module := range modules {
			result[mid] = module
		}
	}
	return result
}

// 清除所有的组件注册记录
func (registrar *myRegistrar) Clear() {
	registrar.rwLock.Lock()
	defer registrar.rwLock.Unlock()
	registrar.moduleTypeMap = map[Type]map[MID]Module{}
}
