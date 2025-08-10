package plugins

import (
	"reflect"
	"sort"
	"sync"
)

// PluginManager 插件管理器 - 线程安全版本
type PluginManager struct {
	plugins []Plugin
	mutex   sync.RWMutex
}

// NewPluginManager 创建插件管理器
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make([]Plugin, 0),
	}
}

// AddPlugin 添加插件 - 线程安全
func (pm *PluginManager) AddPlugin(plugin Plugin) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.plugins = append(pm.plugins, plugin)
	// 按优先级排序
	sort.Slice(pm.plugins, func(i, j int) bool {
		return pm.plugins[i].GetOrder() < pm.plugins[j].GetOrder()
	})
}

// RemovePlugin 移除插件 - 线程安全
func (pm *PluginManager) RemovePlugin(pluginType reflect.Type) bool {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	for i, plugin := range pm.plugins {
		if reflect.TypeOf(plugin) == pluginType {
			pm.plugins = append(pm.plugins[:i], pm.plugins[i+1:]...)
			return true
		}
	}
	return false
}

// GetPlugins 获取所有插件 - 返回副本确保线程安全
func (pm *PluginManager) GetPlugins() []Plugin {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	// 返回插件副本，避免并发修改
	pluginsCopy := make([]Plugin, len(pm.plugins))
	copy(pluginsCopy, pm.plugins)
	return pluginsCopy
}

// HasPlugins 检查是否有插件
func (pm *PluginManager) HasPlugins() bool {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return len(pm.plugins) > 0
}

// GetPluginCount 获取插件数量
func (pm *PluginManager) GetPluginCount() int {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return len(pm.plugins)
}

// Size 获取插件数量 - 向后兼容
func (pm *PluginManager) Size() int {
	return pm.GetPluginCount()
}

// InterceptorChain 拦截器链 - 已弃用，使用 PluginChain 替代
type InterceptorChain struct {
	plugins []Plugin
	index   int
	target  interface{}
}

// CreateInterceptorChain 创建拦截器链 - 已弃用，保留向后兼容
func (pm *PluginManager) CreateInterceptorChain(target interface{}) *InterceptorChain {
	plugins := pm.GetPlugins()
	return &InterceptorChain{
		plugins: plugins,
		index:   0,
		target:  target,
	}
}

// InterceptMethod 拦截方法调用 - 使用新的线程安全 PluginChain
func (pm *PluginManager) InterceptMethod(target interface{}, method reflect.Method, args []interface{}, statementId string, proceed func() (interface{}, error)) (interface{}, error) {
	plugins := pm.GetPlugins()
	if len(plugins) == 0 {
		return proceed()
	}

	// 创建新的线程安全插件链
	chain := NewPluginChain(plugins, target)

	// 创建调用信息
	invocation := &Invocation{
		Target:      target,
		Method:      method,
		Args:        args,
		StatementId: statementId,
		Properties:  make(map[string]interface{}),
		Proceed:     proceed,
		Context:     NewInvocationContext(),
	}

	// 执行插件链
	return chain.Proceed(invocation)
}

// PluginRegistry 插件注册表 - 管理多个插件管理器
type PluginRegistry struct {
	managers map[string]*PluginManager
	mutex    sync.RWMutex
}

// NewPluginRegistry 创建插件注册表
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		managers: make(map[string]*PluginManager),
	}
}

// RegisterManager 注册插件管理器
func (pr *PluginRegistry) RegisterManager(name string, manager *PluginManager) {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()
	pr.managers[name] = manager
}

// GetManager 获取插件管理器
func (pr *PluginRegistry) GetManager(name string) (*PluginManager, bool) {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()
	manager, exists := pr.managers[name]
	return manager, exists
}

// RemoveManager 移除插件管理器
func (pr *PluginRegistry) RemoveManager(name string) bool {
	pr.mutex.Lock()
	defer pr.mutex.Unlock()
	if _, exists := pr.managers[name]; exists {
		delete(pr.managers, name)
		return true
	}
	return false
}

// GetAllManagers 获取所有插件管理器
func (pr *PluginRegistry) GetAllManagers() map[string]*PluginManager {
	pr.mutex.RLock()
	defer pr.mutex.RUnlock()

	// 返回副本
	result := make(map[string]*PluginManager)
	for name, manager := range pr.managers {
		result[name] = manager
	}
	return result
}

// 全局插件注册表
var GlobalPluginRegistry = NewPluginRegistry()

// PluginBuilder 插件构建器 - 用于方便地构建插件管理器
type PluginBuilder struct {
	manager *PluginManager
}

// NewPluginBuilder 创建插件构建器
func NewPluginBuilder() *PluginBuilder {
	return &PluginBuilder{
		manager: NewPluginManager(),
	}
}

// WithPagination 添加分页插件
func (pb *PluginBuilder) WithPagination() *PluginBuilder {
	plugin := NewPaginationPlugin()
	pb.manager.AddPlugin(plugin)
	return pb
}

// WithCustomPlugin 添加自定义插件
func (pb *PluginBuilder) WithCustomPlugin(plugin Plugin) *PluginBuilder {
	pb.manager.AddPlugin(plugin)
	return pb
}

// Build 构建插件管理器
func (pb *PluginBuilder) Build() *PluginManager {
	return pb.manager
}

// BuildAndRegister 构建并注册插件管理器
func (pb *PluginBuilder) BuildAndRegister(name string) *PluginManager {
	manager := pb.Build()
	GlobalPluginRegistry.RegisterManager(name, manager)
	return manager
}
