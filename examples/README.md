# gobatis 示例代码

本目录包含了 gobatis ORM 框架的完整示例代码，演示了框架的各种功能和用法。gobatis 是一个类似 MyBatis 的 Go 语言 ORM 框架，提供 SQL 与业务逻辑解耦、自动参数绑定、结构体结果映射、插件扩展等功能。



## 🚀 快速开始

### 运行完整演示

```bash
cd /Users/chenjianyu/GolandProjects/gobatis/examples/demo
go run .
```

### 运行单元测试

```bash
cd /Users/chenjianyu/GolandProjects/gobatis/examples
go test -v
```

## 📋 功能演示

### 1. 基础功能演示 (`demo/main.go`)

#### 1.1 基本配置和会话管理

```go
// 创建配置
config := gobatis.NewConfiguration()

// 配置数据源
err := config.SetDataSource(
    "mysql",
    "root:password@tcp(localhost:3306)/gobatis_demo?charset=utf8mb4&parseTime=True&loc=Local",
)

// 注册 Mapper
config.RegisterMapper("examples.UserMapper", &examples.UserMapper{})

// 创建会话工厂
factory := gobatis.NewSqlSessionFactory(config)
session := factory.OpenSession()
defer session.Close()
```

#### 1.2 基本 CRUD 操作

```go
// 获取 Mapper 代理
userMapper := session.GetMapper((*examples.UserMapper)(nil))

// 查询单个用户
user, err := userMapper.FindByID(1)

// 查询所有用户
users, err := userMapper.FindAll()

// 插入用户
newUser := &examples.User{
    Username: "john_doe",
    Email:    "john@example.com",
    CreateAt: time.Now(),
}
err = userMapper.Insert(newUser)

// 更新用户
user.Email = "newemail@example.com"
err = userMapper.Update(user)

// 删除用户
err = userMapper.Delete(1)
```

#### 1.3 事务管理

```go
// 开启事务
session := factory.OpenSessionWithAutoCommit(false)
defer session.Close()

// 执行多个操作
err1 := userMapper.Insert(user1)
err2 := userMapper.Insert(user2)

if err1 != nil || err2 != nil {
    // 回滚事务
    session.Rollback()
} else {
    // 提交事务
    session.Commit()
}
```

### 2. 插件系统演示 (`demo/plugin_demo.go`)

#### 2.1 分页插件配置和使用

##### 基本配置

```go
import "gobatis/plugins"

// 方式1: 直接创建插件
paginationPlugin := plugins.NewPaginationPlugin()

// 设置插件属性
properties := map[string]string{
    "defaultPageSize": "20",
    "maxPageSize":     "100",
}
paginationPlugin.SetProperties(properties)

// 添加到配置
config := gobatis.NewConfiguration()
config.AddPlugin(paginationPlugin)
```

```go
// 方式2: 使用插件构建器
manager := plugins.NewPluginBuilder().
    WithPagination().
    Build()
```

##### 分页请求结构

```go
// PageRequest 分页请求参数
type PageRequest struct {
    Page     int    `json:"page"`     // 页码（从1开始）
    Size     int    `json:"size"`     // 每页大小
    Offset   int    `json:"offset"`   // 偏移量（自动计算）
    SortBy   string `json:"sortBy"`   // 排序字段
    SortDir  string `json:"sortDir"`  // 排序方向（ASC/DESC）
}

// PageResult 分页结果
type PageResult struct {
    Data       interface{} `json:"data"`       // 数据列表
    Total      int64       `json:"total"`      // 总记录数
    Page       int         `json:"page"`       // 当前页码
    Size       int         `json:"size"`       // 每页大小
    TotalPages int         `json:"totalPages"` // 总页数
    HasNext    bool        `json:"hasNext"`    // 是否有下一页
    HasPrev    bool        `json:"hasPrev"`    // 是否有上一页
}
```

##### 在 Mapper 中使用分页

```go
// 1. 在 Mapper 接口中定义分页方法
type UserMapper interface {
    // 普通查询
    FindAll() ([]*User, error)
    
    // 分页查询 - 方式1：直接传入 PageRequest
    FindAllWithPage(pageReq *plugins.PageRequest) (*plugins.PageResult, error)
    
    // 分页查询 - 方式2：传入包含分页信息的结构体
    FindByCondition(condition *UserSearchCondition) (*plugins.PageResult, error)
}

// 用户搜索条件（包含分页信息）
type UserSearchCondition struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Page     int    `json:"page"`     // 分页插件会自动识别
    Size     int    `json:"size"`     // 分页插件会自动识别
}
```

##### 分页查询示例

```go
// 示例1: 基本分页查询
pageRequest := &plugins.PageRequest{
    Page:    1,      // 第1页
    Size:    10,     // 每页10条
    SortBy:  "id",   // 按ID排序
    SortDir: "DESC", // 降序
}

result, err := userMapper.FindAllWithPage(pageRequest)
if err != nil {
    log.Fatal(err)
}

// 处理分页结果
fmt.Printf("总记录数: %d\n", result.Total)
fmt.Printf("当前页: %d/%d\n", result.Page, result.TotalPages)
fmt.Printf("是否有下一页: %t\n", result.HasNext)

// 获取数据
users := result.Data.([]*User)
for _, user := range users {
    fmt.Printf("用户: %s\n", user.Username)
}
```

```go
// 示例2: 条件查询 + 分页
condition := &UserSearchCondition{
    Username: "john",
    Page:     2,
    Size:     5,
}

result, err := userMapper.FindByCondition(condition)
// 处理结果...
```

##### 分页 SQL 自动转换

分页插件会自动将原始 SQL 转换为分页 SQL：

```sql
-- 原始 SQL
SELECT * FROM users WHERE username LIKE ?

-- 自动转换为计数 SQL
SELECT COUNT(*) FROM users WHERE username LIKE ?

-- 自动转换为分页 SQL
SELECT * FROM users WHERE username LIKE ? ORDER BY id DESC LIMIT 10 OFFSET 0
```

#### 2.2 插件管理器使用

```go
// 创建插件管理器
manager := plugins.NewPluginManager()

// 添加插件
paginationPlugin := plugins.NewPaginationPlugin()
manager.AddPlugin(paginationPlugin)

// 查询插件信息
fmt.Printf("插件数量: %d\n", manager.GetPluginCount())
fmt.Printf("是否有插件: %t\n", manager.HasPlugins())

// 获取所有插件
allPlugins := manager.GetPlugins()
for i, plugin := range allPlugins {
    fmt.Printf("插件 %d: 优先级 %d, 类型: %T\n", 
        i+1, plugin.GetOrder(), plugin)
}

// 移除插件
pluginType := reflect.TypeOf(paginationPlugin)
removed := manager.RemovePlugin(pluginType)
```

#### 2.3 插件注册表

```go
// 创建插件注册表
registry := plugins.NewPluginRegistry()

// 为不同的 Mapper 注册不同的插件管理器
userManager := plugins.NewPluginManager()
userManager.AddPlugin(plugins.NewPaginationPlugin())
registry.RegisterManager("UserMapper", userManager)

orderManager := plugins.NewPluginManager()
orderManager.AddPlugin(plugins.NewPaginationPlugin())
registry.RegisterManager("OrderMapper", orderManager)

// 获取特定 Mapper 的插件管理器
if manager, exists := registry.GetManager("UserMapper"); exists {
    fmt.Printf("UserMapper 插件数量: %d\n", manager.GetPluginCount())
}
```

#### 2.4 自定义插件开发

```go
// 自定义日志插件
type LoggingPlugin struct {
    properties map[string]string
    order      int
}

func NewLoggingPlugin() *LoggingPlugin {
    return &LoggingPlugin{
        properties: make(map[string]string),
        order:      50, // 中等优先级
    }
}

func (p *LoggingPlugin) Intercept(invocation *plugins.Invocation) (interface{}, error) {
    startTime := time.Now()
    
    // 记录方法调用开始
    fmt.Printf("🔍 开始执行方法: %s\n", invocation.Method.Name)
    
    // 执行原方法
    result, err := invocation.Proceed()
    
    // 记录执行时间
    duration := time.Since(startTime)
    if err != nil {
        fmt.Printf("❌ 方法执行失败: %s, 耗时: %v, 错误: %v\n",
            invocation.Method.Name, duration, err)
    } else {
        fmt.Printf("✅ 方法执行成功: %s, 耗时: %v\n",
            invocation.Method.Name, duration)
    }
    
    return result, err
}

func (p *LoggingPlugin) SetProperties(properties map[string]string) {
    p.properties = properties
}

func (p *LoggingPlugin) GetOrder() int {
    return p.order
}

// 使用自定义插件
loggingPlugin := NewLoggingPlugin()
manager.AddPlugin(loggingPlugin)
```

## 🏗️ 核心组件

### User 实体定义

```go
type User struct {
    ID       int64     `json:"id" db:"id"`
    Username string    `json:"username" db:"username"`
    Email    string    `json:"email" db:"email"`
    CreateAt time.Time `json:"create_at" db:"create_at"`
}
```

### UserMapper 接口

```go
type UserMapper interface {
    FindByID(id int64) (*User, error)
    FindByUsername(username string) (*User, error)
    FindAll() ([]*User, error)
    Insert(user *User) error
    Update(user *User) error
    Delete(id int64) error
    Count() (int64, error)
    
    // 分页查询方法
    FindAllWithPage(pageReq *plugins.PageRequest) (*plugins.PageResult, error)
    FindByCondition(condition *UserSearchCondition) (*plugins.PageResult, error)
}
```

### Mapper XML 配置

```xml
<?xml version="1.0" encoding="UTF-8"?>
<mapper namespace="examples.UserMapper">
    <select id="FindByID" resultType="examples.User">
        SELECT id, username, email, create_at FROM users WHERE id = ?
    </select>
    
    <select id="FindAll" resultType="examples.User">
        SELECT id, username, email, create_at FROM users ORDER BY id
    </select>
    
    <select id="FindAllWithPage" resultType="examples.User">
        SELECT id, username, email, create_at FROM users ORDER BY id
    </select>
    
    <insert id="Insert">
        INSERT INTO users (username, email, create_at) VALUES (?, ?, ?)
    </insert>
    
    <update id="Update">
        UPDATE users SET username = ?, email = ? WHERE id = ?
    </update>
    
    <delete id="Delete">
        DELETE FROM users WHERE id = ?
    </delete>
</mapper>
```

## 🔧 高级功能

### 参数绑定

```go
// 支持多种参数类型
// 1. 基本类型
userMapper.FindByID(123)

// 2. 结构体
condition := &UserSearchCondition{
    Username: "john",
    Email:    "john@example.com",
}
userMapper.FindByCondition(condition)

// 3. Map
params := map[string]interface{}{
    "username": "john",
    "email":    "john@example.com",
}

// 4. 切片
ids := []int64{1, 2, 3, 4, 5}
```

### 结果映射

```go
// 自动映射到结构体
user, err := userMapper.FindByID(1)

// 映射到切片
users, err := userMapper.FindAll()

// 映射到分页结果
pageResult, err := userMapper.FindAllWithPage(pageRequest)
```

### 错误处理

```go
user, err := userMapper.FindByID(1)
if err != nil {
    switch {
    case errors.Is(err, sql.ErrNoRows):
        fmt.Println("用户不存在")
    case strings.Contains(err.Error(), "connection"):
        fmt.Println("数据库连接错误")
    default:
        fmt.Printf("其他错误: %v\n", err)
    }
}
```

## 📊 性能优化

### 连接池配置

```go
config := gobatis.NewConfiguration()
config.SetDataSource(
    "mysql",
    "root:password@tcp(localhost:3306)/gobatis_demo?charset=utf8mb4&parseTime=True&loc=Local",
)

// 配置连接池（如果支持）
config.SetMaxOpenConns(100)
config.SetMaxIdleConns(10)
config.SetConnMaxLifetime(time.Hour)
```

### 插件优化

```go
// 设置插件执行顺序（数字越小优先级越高）
plugin1.order = 10  // 高优先级
plugin2.order = 50  // 中等优先级
plugin3.order = 100 // 低优先级
```

## 🧪 测试

### Mock 实现

项目提供了完整的 Mock 实现用于测试：

```go
// 创建 Mock Mapper
mockMapper := examples.NewMockUserMapper()

// 模拟数据
user := &examples.User{
    ID:       1,
    Username: "test_user",
    Email:    "test@example.com",
    CreateAt: time.Now(),
}

// 执行操作
err := mockMapper.Insert(user)
foundUser, err := mockMapper.FindByID(1)
```

### 单元测试

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestUserMapper

# 运行基准测试
go test -bench=.
```

## 📝 最佳实践

1. **配置管理**: 使用配置文件管理数据源和 Mapper
2. **事务处理**: 合理使用事务确保数据一致性
3. **错误处理**: 完善的错误处理和日志记录
4. **插件使用**: 根据需要选择和配置插件
5. **性能优化**: 使用连接池和缓存机制
6. **分页查询**: 大数据量查询时使用分页插件
7. **SQL 优化**: 合理设计 SQL 语句和索引

## ⚠️ 注意事项

- 确保数据库连接配置正确
- 注意事务的正确使用
- 插件的执行顺序很重要
- Mock 实现仅用于测试和演示
- 分页插件会自动修改 SQL，注意 SQL 兼容性
- 大数据量分页时注意性能影响

## 🔮 扩展示例

基于这些示例可以进一步扩展：

1. **多表关联查询**
2. **复杂条件查询**
3. **批量操作**
4. **缓存集成**
5. **性能监控插件**
6. **数据库迁移工具**
7. **读写分离**

---

这些示例展示了 gobatis 框架的强大功能和灵活性，特别是分页插件的便捷使用，可以作为学习和开发的完整参考。