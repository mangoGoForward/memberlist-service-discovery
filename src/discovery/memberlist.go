package discovery

import (
	"log"
	"sync"

	"github.com/hashicorp/memberlist"
)

type MemberlistService struct {
	list    *memberlist.Memberlist
	members map[string]string // 存储节点信息
	mu      sync.RWMutex
}

// NewMemberlistService 初始化 Memberlist
func NewMemberlistService(config *memberlist.Config) (*MemberlistService, error) {
	service := &MemberlistService{
		members: make(map[string]string),
	}
	config.Events = service // 设置事件回调
	list, err := memberlist.Create(config)
	if err != nil {
		return nil, err
	}
	service.list = list
	return service, nil
}

// Join 加入现有的集群
func (m *MemberlistService) Join(existing []string) error {
	_, err := m.list.Join(existing)
	return err
}

// Members 获取当前集群中的所有成员
func (m *MemberlistService) Members() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	addresses := make(map[string]string, len(m.members))
	for nodeName, addr := range m.members {
		addresses[nodeName] = addr
		//addresses = append(addresses, addr)
	}
	return addresses
}

// Shutdown 关闭 Memberlist
func (m *MemberlistService) Shutdown() error {
	return m.list.Shutdown()
}

// NotifyJoin 处理节点加入事件
func (m *MemberlistService) NotifyJoin(node *memberlist.Node) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.members[node.Name] = node.Address()
	log.Printf("服务注册: %s (%s)", node.Name, node.Address())
}

// NotifyLeave 处理节点离开事件
func (m *MemberlistService) NotifyLeave(node *memberlist.Node) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.members, node.Name)
	log.Printf("服务离开: %s (%s)", node.Name, node.Address())
}

// NotifyUpdate 处理节点更新事件
func (m *MemberlistService) NotifyUpdate(node *memberlist.Node) {
	log.Printf("服务更新更新: %s (%s)", node.Name, node.Address())
}

// DefaultConfig 返回默认的 Memberlist 配置
func DefaultConfig(nodeName string) *memberlist.Config {
	config := memberlist.DefaultLANConfig()
	config.Name = nodeName
	config.BindPort = 6789
	config.Logger = log.Default()
	return config
}
