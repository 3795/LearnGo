package common

// 存储类型（文件存储到何处）
type StoreType int

const (
	_ StoreType = iota
	// 本地节点
	StoreLocal
	// Ceph集群
	StoreCeph
	//阿里OSS
	StoreOSS
	// 混合存储（Ceph和OSS）
	StoreMix
	// 所有的类型都存储一份儿
	StoreAll
)
