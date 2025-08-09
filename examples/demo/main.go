package main

import (
	"fmt"
	"time"

	"gobatis"
	"gobatis/examples"
)

func main() {
	fmt.Println("🚀 gobatis 框架功能演示")
	fmt.Println("================================")

	// 演示基本用法
	demonstrateBasicUsage()

	// 演示配置用法
	demonstrateConfigurationUsage()

	// 演示参数绑定
	demonstrateParameterBinding()

	// 演示事务用法
	demonstrateTransactionUsage()

	// 演示 Mapper 代理
	demonstrateMapperProxy()

	// 演示错误处理
	demonstrateErrorHandling()

	// 演示数据类型
	demonstrateDataTypes()

	// 演示插件系统
	runPluginDemonstrations()

	fmt.Println("\n================================")
	fmt.Println("🎉 所有演示完成！")
}

// demonstrateBasicUsage 演示基本用法
func demonstrateBasicUsage() {
	fmt.Println("\n=== 基本用法演示 ===")

	// 使用 Mock 数据进行演示
	mockMapper := examples.NewMockUserMapper()
	userService := examples.NewUserService(mockMapper)

	// 1. 创建用户
	fmt.Println("\n1. 创建用户:")
	user, err := userService.CreateUser("demo_user", "demo@example.com")
	if err != nil {
		fmt.Printf("❌ 创建用户失败: %v\n", err)
	} else {
		fmt.Printf("✅ 创建用户成功: ID=%d, Username=%s, Email=%s\n",
			user.ID, user.Username, user.Email)
	}

	// 2. 查询用户
	fmt.Println("\n2. 查询用户:")
	foundUser, err := userService.GetUser(user.ID)
	if err != nil {
		fmt.Printf("❌ 查询用户失败: %v\n", err)
	} else if foundUser != nil {
		fmt.Printf("✅ 查询用户成功: %+v\n", foundUser)
	} else {
		fmt.Println("⚠️  用户不存在")
	}

	// 3. 更新用户邮箱
	fmt.Println("\n3. 更新用户邮箱:")
	err = userService.UpdateUserEmail(user.ID, "updated@example.com")
	if err != nil {
		fmt.Printf("❌ 更新用户失败: %v\n", err)
	} else {
		fmt.Println("✅ 更新用户成功")
	}

	// 4. 搜索用户
	fmt.Println("\n4. 搜索用户:")
	users, err := userService.SearchUsers("demo")
	if err != nil {
		fmt.Printf("❌ 搜索用户失败: %v\n", err)
	} else {
		fmt.Printf("✅ 搜索到 %d 个用户\n", len(users))
	}

	// 5. 获取用户总数
	fmt.Println("\n5. 获取用户总数:")
	count, err := userService.GetUserCount()
	if err != nil {
		fmt.Printf("❌ 获取用户总数失败: %v\n", err)
	} else {
		fmt.Printf("✅ 用户总数: %d\n", count)
	}

	// 6. 删除用户
	fmt.Println("\n6. 删除用户:")
	err = userService.DeleteUser(user.ID)
	if err != nil {
		fmt.Printf("❌ 删除用户失败: %v\n", err)
	} else {
		fmt.Println("✅ 删除用户成功")
	}
}

// demonstrateConfigurationUsage 演示配置用法
func demonstrateConfigurationUsage() {
	fmt.Println("\n=== 配置系统演示 ===")

	// 创建配置
	config := gobatis.NewConfiguration()
	fmt.Println("✅ 创建配置成功")

	// 配置数据源（演示用，实际可能会失败）
	err := config.SetDataSource(
		"mysql",
		"root:password@tcp(localhost:3306)/gobatis_demo?charset=utf8mb4&parseTime=True&loc=Local",
	)
	if err != nil {
		fmt.Printf("⚠️  数据源配置失败 (这是正常的): %v\n", err)
	} else {
		fmt.Println("✅ 数据源配置成功")
	}

	// 添加 Mapper XML
	err = config.AddMapperXML("examples/user_mapper.xml")
	if err != nil {
		fmt.Printf("⚠️  Mapper XML 配置失败: %v\n", err)
	} else {
		fmt.Println("✅ Mapper XML 配置成功")
	}

	// 创建会话工厂
	factory := gobatis.NewSqlSessionFactory(config)
	fmt.Println("✅ 会话工厂创建成功")

	// 打开会话
	session := factory.OpenSession()
	defer session.Close()
	fmt.Println("✅ 会话打开成功")
}

// demonstrateParameterBinding 演示参数绑定
func demonstrateParameterBinding() {
	fmt.Println("\n=== 参数绑定演示 ===")

	// 演示结构体参数绑定
	type UserQuery struct {
		Username string `db:"username"`
		Email    string `db:"email"`
		MinAge   int    `db:"min_age"`
		MaxAge   int    `db:"max_age"`
	}

	query := &UserQuery{
		Username: "john%",
		Email:    "john@example.com",
		MinAge:   18,
		MaxAge:   65,
	}

	fmt.Printf("✅ 查询条件结构体: %+v\n", query)

	// 演示 Map 参数绑定
	params := map[string]interface{}{
		"username": "john_doe",
		"email":    "john@example.com",
		"page":     1,
		"size":     10,
	}

	fmt.Printf("✅ Map 参数: %+v\n", params)

	// 演示单个参数绑定
	userId := int64(123)
	fmt.Printf("✅ 单个参数: %d\n", userId)
}

// demonstrateTransactionUsage 演示事务用法
func demonstrateTransactionUsage() {
	fmt.Println("\n=== 事务管理演示 ===")

	config := gobatis.NewConfiguration()
	factory := gobatis.NewSqlSessionFactory(config)

	// 手动事务管理
	session := factory.OpenSessionWithAutoCommit(false) // 关闭自动提交
	defer session.Close()

	fmt.Println("✅ 开启手动事务模式")

	// 模拟事务操作
	fmt.Println("📝 执行事务操作...")
	fmt.Println("  - 插入用户1")
	fmt.Println("  - 插入用户2")
	fmt.Println("  - 更新用户1")

	// 模拟提交事务
	err := session.Commit()
	if err != nil {
		fmt.Printf("❌ 提交事务失败: %v\n", err)
		session.Rollback()
		fmt.Println("🔄 已回滚事务")
	} else {
		fmt.Println("✅ 事务提交成功")
	}
}

// demonstrateMapperProxy 演示 Mapper 代理
func demonstrateMapperProxy() {
	fmt.Println("\n=== Mapper 代理演示 ===")

	config := gobatis.NewConfiguration()
	factory := gobatis.NewSqlSessionFactory(config)
	session := factory.OpenSession()
	defer session.Close()

	// 获取 Mapper 代理
	userMapper := session.GetMapper((*examples.UserMapper)(nil))
	if userMapper != nil {
		fmt.Println("✅ Mapper 代理创建成功")
		fmt.Printf("代理类型: %T\n", userMapper)
	} else {
		fmt.Println("❌ Mapper 代理创建失败")
	}
}

// demonstrateErrorHandling 演示错误处理
func demonstrateErrorHandling() {
	fmt.Println("\n=== 错误处理演示 ===")

	mockMapper := examples.NewMockUserMapper()
	userService := examples.NewUserService(mockMapper)

	// 模拟数据库错误
	mockMapper.SetError(true)

	// 测试各种错误情况
	fmt.Println("\n1. 测试创建用户时的错误:")
	_, err := userService.CreateUser("error_user", "error@example.com")
	if err != nil {
		fmt.Printf("✅ 正确捕获错误: %v\n", err)
	}

	fmt.Println("\n2. 测试查询用户时的错误:")
	_, err = userService.GetUser(1)
	if err != nil {
		fmt.Printf("✅ 正确捕获错误: %v\n", err)
	}

	fmt.Println("\n3. 测试更新用户时的错误:")
	err = userService.UpdateUserEmail(1, "new@example.com")
	if err != nil {
		fmt.Printf("✅ 正确捕获错误: %v\n", err)
	}

	// 恢复正常状态
	mockMapper.SetError(false)
	fmt.Println("\n✅ 错误处理演示完成")
}

// demonstrateDataTypes 演示数据类型处理
func demonstrateDataTypes() {
	fmt.Println("\n=== 数据类型演示 ===")

	// 创建用户实例
	user := &examples.User{
		ID:       1,
		Username: "test_user",
		Email:    "test@example.com",
		CreateAt: time.Now(),
	}

	fmt.Printf("✅ 用户实体: %+v\n", user)
	fmt.Printf("ID 类型: %T\n", user.ID)
	fmt.Printf("Username 类型: %T\n", user.Username)
	fmt.Printf("Email 类型: %T\n", user.Email)
	fmt.Printf("CreateAt 类型: %T\n", user.CreateAt)

	// 演示时间格式化
	fmt.Printf("创建时间格式化: %s\n", user.CreateAt.Format("2006-01-02 15:04:05"))
}
