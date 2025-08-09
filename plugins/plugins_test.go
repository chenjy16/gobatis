package plugins

import (
	"reflect"
	"testing"
)

// TestPaginationPlugin 测试分页插件
func TestPaginationPlugin(t *testing.T) {
	plugin := NewPaginationPlugin()

	// 创建分页请求
	pageRequest := &PageRequest{
		Page: 1,
		Size: 10,
	}

	method := reflect.Method{
		Name: "SelectUsers",
		Func: reflect.ValueOf(func() []string { return []string{"user1", "user2"} }),
	}

	invocation := &Invocation{
		Target:      &struct{}{},
		Method:      method,
		Args:        []interface{}{pageRequest},
		StatementId: "selectUsers",
		Properties:  map[string]interface{}{"sql": "SELECT * FROM users"},
	}

	// 模拟 Proceed 方法
	invocation.Proceed = func() (interface{}, error) {
		return []string{"user1", "user2"}, nil
	}

	result, err := plugin.Intercept(invocation)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	pageResult, ok := result.(*PageResult)
	if !ok {
		t.Errorf("Expected PageResult, got %T", result)
	}

	if pageResult.Page != 1 || pageResult.Size != 10 {
		t.Errorf("Expected page=1, size=10, got page=%d, size=%d", pageResult.Page, pageResult.Size)
	}
}



// TestPluginManager 测试插件管理器
func TestPluginManager(t *testing.T) {
	manager := NewPluginManager()

	// 创建测试插件
	plugin1 := &TestPlugin{order: 10}
	plugin2 := &TestPlugin{order: 5}
	plugin3 := &TestPlugin{order: 15}

	// 添加插件
	manager.AddPlugin(plugin1)
	manager.AddPlugin(plugin2)
	manager.AddPlugin(plugin3)

	plugins := manager.GetPlugins()
	if len(plugins) != 3 {
		t.Errorf("Expected 3 plugins, got %d", len(plugins))
	}

	// 验证插件按优先级排序
	if plugins[0].GetOrder() != 5 || plugins[1].GetOrder() != 10 || plugins[2].GetOrder() != 15 {
		t.Error("Plugins not sorted by order")
	}

	// 测试移除插件
	manager.RemovePlugin(reflect.TypeOf(plugin2))
	plugins = manager.GetPlugins()
	if len(plugins) != 2 {
		t.Errorf("Expected 2 plugins after removal, got %d", len(plugins))
	}
}

// TestPlugin 测试插件实现
type TestPlugin struct {
	order int
}

func (p *TestPlugin) Intercept(invocation *Invocation) (interface{}, error) {
	return invocation.Proceed()
}

func (p *TestPlugin) SetProperties(properties map[string]string) {
	// 空实现
}

func (p *TestPlugin) GetOrder() int {
	return p.order
}

// TestPluginBuilder 测试插件构建器
func TestPluginBuilder(t *testing.T) {
	builder := NewPluginBuilder()
	manager := builder.
		WithPagination().
		Build()

	if manager.Size() != 1 {
		t.Errorf("Expected 1 plugin, got %d", manager.Size())
	}

	plugins := manager.GetPlugins()
	
	// 验证插件类型
	if _, ok := plugins[0].(*PaginationPlugin); !ok {
		t.Error("Expected plugin to be PaginationPlugin")
	}
}

// TestPluginRegistry 测试插件注册表
func TestPluginRegistry(t *testing.T) {
	registry := NewPluginRegistry()

	// 测试获取不存在的管理器
	manager, exists := registry.GetManager("nonexistent")
	if exists || manager != nil {
		t.Error("Expected nil for nonexistent manager")
	}

	// 测试注册管理器
	newManager := NewPluginManager()
	registry.RegisterManager("custom", newManager)
	
	retrievedManager, exists := registry.GetManager("custom")
	if !exists || retrievedManager != newManager {
		t.Error("Expected registered manager to be retrieved")
	}

	// 测试移除管理器
	removed := registry.RemoveManager("custom")
	if !removed {
		t.Error("Expected manager to be removed")
	}

	_, exists = registry.GetManager("custom")
	if exists {
		t.Error("Expected manager to be removed")
	}
}

// TestGlobalRegistry 测试全局注册表
func TestGlobalRegistry(t *testing.T) {
	// 测试全局插件注册表
	manager := NewPluginManager()
	plugin := &TestPlugin{order: 1}
	manager.AddPlugin(plugin)

	GlobalPluginRegistry.RegisterManager("global", manager)

	retrievedManager, exists := GlobalPluginRegistry.GetManager("global")
	if !exists || retrievedManager.Size() != 1 {
		t.Errorf("Expected 1 plugin in global manager, got %d", retrievedManager.Size())
	}

	plugins := retrievedManager.GetPlugins()
	if plugins[0] != plugin {
		t.Error("Expected registered plugin to be in global manager")
	}
}

// TestPageRequest 测试分页请求
func TestPageRequest(t *testing.T) {
	pageReq := &PageRequest{
		Page: 2,
		Size: 20,
	}

	// 测试偏移量计算
	plugin := NewPaginationPlugin()
	extractedReq := plugin.extractPageRequest([]interface{}{pageReq})
	
	if extractedReq == nil {
		t.Error("Expected to extract page request")
	}
	
	expectedOffset := (2 - 1) * 20
	if extractedReq.Offset != expectedOffset {
		t.Errorf("Expected offset %d, got %d", expectedOffset, extractedReq.Offset)
	}
}