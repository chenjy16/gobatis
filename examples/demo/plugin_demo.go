package main

import (
	"fmt"
	"reflect"
	"time"

	"gobatis"
	"gobatis/plugins"
)

// demonstratePluginManager 演示插件管理器
func demonstratePluginManager() {
	fmt.Println("\n=== 插件管理器演示 ===")

	// 创建插件管理器
	manager := plugins.NewPluginManager()
	fmt.Println("✅ 创建插件管理器")

	// 添加分页插件
	paginationPlugin := plugins.NewPaginationPlugin()
	manager.AddPlugin(paginationPlugin)
	fmt.Println("✅ 添加分页插件")

	// 显示插件信息
	fmt.Printf("插件数量: %d\n", manager.GetPluginCount())
	fmt.Printf("是否有插件: %t\n", manager.HasPlugins())

	// 获取所有插件
	allPlugins := manager.GetPlugins()
	for i, plugin := range allPlugins {
		fmt.Printf("插件 %d: 优先级 %d, 类型: %T\n", i+1, plugin.GetOrder(), plugin)
	}

	// 使用插件构建器
	fmt.Println("\n--- 插件构建器演示 ---")
	builder := plugins.NewPluginBuilder()
	builtManager := builder.
		WithPagination().
		Build()

	fmt.Printf("✅ 通过构建器创建的插件管理器，插件数量: %d\n", builtManager.GetPluginCount())
}

// demonstratePaginationPlugin 演示分页插件
func demonstratePaginationPlugin() {
	fmt.Println("\n=== 分页插件演示 ===")

	// 创建分页请求
	pageRequest := &plugins.PageRequest{
		Page: 1,
		Size: 10,
	}

	fmt.Printf("✅ 分页请求: 第 %d 页，每页 %d 条\n", pageRequest.Page, pageRequest.Size)

	// 模拟分页结果
	mockData := []interface{}{
		"用户1", "用户2", "用户3", "用户4", "用户5",
		"用户6", "用户7", "用户8", "用户9", "用户10",
	}

	pageResult := &plugins.PageResult{
		Data:       mockData,
		Page:       pageRequest.Page,
		Size:       pageRequest.Size,
		Total:      100, // 假设总共100条记录
		TotalPages: 10,  // 总共10页
	}

	fmt.Printf("✅ 分页结果:\n")
	fmt.Printf("  当前页: %d\n", pageResult.Page)
	fmt.Printf("  每页大小: %d\n", pageResult.Size)
	fmt.Printf("  总记录数: %d\n", pageResult.Total)
	fmt.Printf("  总页数: %d\n", pageResult.TotalPages)
	fmt.Printf("  当前页数据: %v\n", pageResult.Data)

	// 演示分页计算
	fmt.Println("\n--- 分页计算演示 ---")
	for page := 1; page <= 3; page++ {
		offset := (page - 1) * pageRequest.Size
		fmt.Printf("第 %d 页: OFFSET %d, LIMIT %d\n", page, offset, pageRequest.Size)
	}
}

// demonstratePluginConfiguration 演示插件配置
func demonstratePluginConfiguration() {
	fmt.Println("\n=== 插件配置演示 ===")

	// 创建配置
	config := gobatis.NewConfiguration()
	fmt.Println("✅ 创建配置")

	// 创建分页插件并设置属性
	paginationPlugin := plugins.NewPaginationPlugin()
	
	// 设置插件属性
	properties := map[string]string{
		"defaultPageSize": "20",
		"maxPageSize":     "100",
	}
	paginationPlugin.SetProperties(properties)
	
	fmt.Printf("✅ 设置插件属性: %+v\n", properties)
	fmt.Printf("插件优先级: %d\n", paginationPlugin.GetOrder())

	// 创建会话工厂
	factory := gobatis.NewSqlSessionFactory(config)
	session := factory.OpenSession()
	defer session.Close()

	fmt.Println("✅ 会话创建成功，插件已集成")
}

// demonstratePluginChain 演示插件链
func demonstratePluginChain() {
	fmt.Println("\n=== 插件链演示 ===")

	// 创建多个插件
	pluginList := []plugins.Plugin{
		plugins.NewPaginationPlugin(),
	}

	// 创建插件链
	target := "目标对象"
	_ = plugins.NewPluginChain(pluginList, target)

	fmt.Printf("✅ 创建插件链，包含 %d 个插件\n", len(pluginList))

	// 模拟调用上下文
	context := plugins.NewInvocationContext()
	
	// 设置插件状态
	context.SetPluginState("pagination", map[string]interface{}{
		"page": 1,
		"size": 10,
	})

	fmt.Println("✅ 设置插件状态")

	// 获取插件状态
	if state, exists := context.GetPluginState("pagination"); exists {
		fmt.Printf("分页插件状态: %+v\n", state)
	}

	// 添加回滚函数
	context.AddRollbackFunc(func() error {
		fmt.Println("🔄 执行回滚操作: 清理分页缓存")
		return nil
	})

	fmt.Println("✅ 添加回滚函数")

	// 模拟执行时间
	fmt.Printf("执行开始时间: %v\n", context.StartTime.Format("15:04:05.000"))
	time.Sleep(1 * time.Millisecond)
	fmt.Printf("执行耗时: %v\n", time.Since(context.StartTime))

	// 执行回滚演示
	errors := context.ExecuteRollback()
	if len(errors) == 0 {
		fmt.Println("✅ 回滚操作执行成功")
	} else {
		fmt.Printf("❌ 回滚操作有错误: %v\n", errors)
	}
}

// demonstratePluginRegistry 演示插件注册表
func demonstratePluginRegistry() {
	fmt.Println("\n=== 插件注册表演示 ===")

	// 创建插件注册表
	registry := plugins.NewPluginRegistry()
	fmt.Println("✅ 创建插件注册表")

	// 为不同的 Mapper 创建不同的插件管理器
	userManager := plugins.NewPluginManager()
	userManager.AddPlugin(plugins.NewPaginationPlugin())
	registry.RegisterManager("UserMapper", userManager)

	orderManager := plugins.NewPluginManager()
	orderManager.AddPlugin(plugins.NewPaginationPlugin())
	registry.RegisterManager("OrderMapper", orderManager)

	fmt.Println("✅ 注册插件管理器")

	// 获取特定 Mapper 的插件管理器
	if manager, exists := registry.GetManager("UserMapper"); exists {
		fmt.Printf("UserMapper 插件数量: %d\n", manager.GetPluginCount())
	}

	if manager, exists := registry.GetManager("OrderMapper"); exists {
		fmt.Printf("OrderMapper 插件数量: %d\n", manager.GetPluginCount())
	}

	// 使用全局插件注册表
	fmt.Println("\n--- 全局插件注册表演示 ---")
	globalManager := plugins.NewPluginBuilder().
		WithPagination().
		BuildAndRegister("GlobalMapper")

	fmt.Printf("✅ 全局注册插件管理器，插件数量: %d\n", globalManager.GetPluginCount())
}

// demonstrateAdvancedFeatures 演示高级功能
func demonstrateAdvancedFeatures() {
	fmt.Println("\n=== 高级功能演示 ===")

	// 演示插件移除
	fmt.Println("\n--- 插件移除演示 ---")
	manager := plugins.NewPluginManager()
	
	// 添加插件
	paginationPlugin := plugins.NewPaginationPlugin()
	manager.AddPlugin(paginationPlugin)
	fmt.Printf("添加插件后数量: %d\n", manager.GetPluginCount())

	// 移除插件
	pluginType := reflect.TypeOf(paginationPlugin)
	removed := manager.RemovePlugin(pluginType)
	fmt.Printf("移除插件成功: %t, 剩余插件数量: %d\n", removed, manager.GetPluginCount())

	// 演示线程安全性
	fmt.Println("\n--- 线程安全演示 ---")
	safeManager := plugins.NewPluginManager()
	
	// 模拟并发添加插件
	done := make(chan bool, 2)
	
	go func() {
		for i := 0; i < 3; i++ {
			plugin := plugins.NewPaginationPlugin()
			safeManager.AddPlugin(plugin)
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()
	
	go func() {
		for i := 0; i < 3; i++ {
			count := safeManager.GetPluginCount()
			fmt.Printf("当前插件数量: %d\n", count)
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	// 等待完成
	<-done
	<-done
	
	fmt.Printf("✅ 并发操作完成，最终插件数量: %d\n", safeManager.GetPluginCount())

	// 演示插件排序
	fmt.Println("\n--- 插件排序演示 ---")
	sortManager := plugins.NewPluginManager()
	
	// 添加多个相同类型的插件（实际使用中不常见，仅用于演示）
	for i := 0; i < 3; i++ {
		plugin := plugins.NewPaginationPlugin()
		sortManager.AddPlugin(plugin)
	}
	
	plugins := sortManager.GetPlugins()
	fmt.Printf("插件按优先级排序，数量: %d\n", len(plugins))
	for i, plugin := range plugins {
		fmt.Printf("  插件 %d: 优先级 %d\n", i+1, plugin.GetOrder())
	}
}

// runPluginDemonstrations 运行所有插件演示
func runPluginDemonstrations() {
	fmt.Println("\n🔌 gobatis 插件系统演示")
	fmt.Println("================================")

	// 演示插件管理器
	demonstratePluginManager()

	// 演示分页插件
	demonstratePaginationPlugin()

	// 演示插件配置
	demonstratePluginConfiguration()

	// 演示插件链
	demonstratePluginChain()

	// 演示插件注册表
	demonstratePluginRegistry()

	// 演示高级功能
	demonstrateAdvancedFeatures()

	fmt.Println("\n================================")
	fmt.Println("🎉 插件演示完成！")
}