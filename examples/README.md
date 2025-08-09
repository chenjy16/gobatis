# GoBatis Example 完整示例

这个目录包含了 GoBatis Example 功能的完整使用示例，从数据库表结构到 Go 代码实现的完整演示。

## 文件结构

```
examples/
├── README.md                    # 本说明文件
├── schema.sql                   # 数据库表结构和测试数据
├── models.go                    # Go 结构体定义
├── dao.go                       # 数据访问层接口和实现
├── complete_example.go          # 完整的使用示例
├── complete_example_test.go     # 测试文件
├── example_demo.go              # 基础演示代码
├── example_demo_test.go         # 基础演示测试
└── main.go                      # 主程序入口
```

## 功能特性

### 1. 数据库表结构 (schema.sql)
- **users**: 用户表，包含基本用户信息
- **departments**: 部门表
- **user_roles**: 用户角色表
- **user_login_logs**: 用户登录日志表

### 2. Go 模型 (models.go)
- **User**: 用户实体
- **Department**: 部门实体
- **UserRole**: 用户角色实体
- **UserLoginLog**: 用户登录日志实体
- **UserWithDepartment**: 连接查询结果
- **UserStatistics**: 用户统计信息
- **DepartmentStatistics**: 部门统计信息
- **SearchCondition**: 搜索条件
- **PageRequest/PageResponse**: 分页请求和响应

### 3. 数据访问层 (dao.go)
提供了完整的 DAO 接口，包括：
- 基础 CRUD 操作
- 复杂查询条件
- 连接查询
- 统计查询
- 分页查询
- 动态查询

### 4. Example 功能演示

#### 基础查询
```go
ex := example.NewExample()
criteria := ex.CreateCriteria()
criteria.AndEqualTo("status", "active")
criteria.AndGreaterThan("age", 18)
sql, args := ex.BuildSQL("SELECT * FROM users")
```

#### 复杂查询
```go
ex := example.NewExample()
ex.SetDistinct(true)
criteria := ex.CreateCriteria()
criteria.AndEqualTo("status", "active")
criteria.AndBetween("age", 25, 35)
criteria.AndLike("real_name", "%张%")
criteria.AndIn("department", []interface{}{"IT", "HR"})
ex.SetOrderByClause("salary DESC")
ex.SetLimit(0, 10)
```

#### OR 查询
```go
ex := example.NewExample()
criteria1 := ex.CreateCriteria()
criteria1.AndEqualTo("department", "IT")
criteria1.AndGreaterThan("salary", 10000)

criteria2 := ex.CreateCriteria()
criteria2.AndEqualTo("department", "HR")
criteria2.AndGreaterThan("salary", 8000)
ex.Or(*criteria2)
```

## 支持的查询条件

### 比较操作
- `AndEqualTo(column, value)` - 等于
- `AndNotEqualTo(column, value)` - 不等于
- `AndGreaterThan(column, value)` - 大于
- `AndGreaterThanOrEqualTo(column, value)` - 大于等于
- `AndLessThan(column, value)` - 小于
- `AndLessThanOrEqualTo(column, value)` - 小于等于

### 范围操作
- `AndBetween(column, value1, value2)` - 在范围内
- `AndNotBetween(column, value1, value2)` - 不在范围内
- `AndIn(column, values)` - 在列表中
- `AndNotIn(column, values)` - 不在列表中

### 模糊查询
- `AndLike(column, value)` - 模糊匹配
- `AndNotLike(column, value)` - 不匹配

### 空值检查
- `AndIsNull(column)` - 为空
- `AndIsNotNull(column)` - 不为空

### SQL 子句
- `SetDistinct(true)` - 去重
- `SetOrderByClause("column DESC")` - 排序
- `SetLimit(offset, limit)` - 分页

## 运行示例

### 1. 运行完整示例
```bash
cd /Users/chenjianyu/GolandProjects/gobatis/examples
go run complete_example.go models.go dao.go
```

### 2. 运行测试
```bash
go test -v
```

### 3. 运行性能测试
```bash
go test -bench=.
```

## 使用场景

### 1. 基础查询
适用于简单的条件查询，如根据状态、年龄等字段查询用户。

### 2. 复杂查询
适用于多条件组合查询，支持 AND、OR 逻辑组合。

### 3. 分页查询
支持 LIMIT 和 OFFSET，适用于大数据量的分页显示。

### 4. 统计查询
支持 COUNT、SUM、AVG 等聚合函数查询。

### 5. 动态查询
根据前端传入的条件动态构建查询语句。

### 6. 连接查询
支持多表连接查询，获取关联数据。

## 最佳实践

1. **类型安全**: 使用指针类型处理可选字段
2. **参数化查询**: 所有查询都使用参数化，防止 SQL 注入
3. **链式调用**: 支持方法链式调用，代码更简洁
4. **错误处理**: 完善的错误处理机制
5. **性能优化**: 合理使用索引和分页
6. **代码复用**: DAO 层封装，提高代码复用性

## 注意事项

1. 确保数据库连接配置正确
2. 根据实际数据库类型调整 SQL 语法
3. 注意字段名和数据库表字段的对应关系
4. 合理设置分页大小，避免一次查询过多数据
5. 复杂查询建议添加适当的索引

## 扩展功能

这个示例展示了 GoBatis Example 的核心功能，你可以根据实际需求进行扩展：

1. 添加更多的查询条件
2. 支持更复杂的连接查询
3. 添加缓存机制
4. 集成事务管理
5. 添加数据验证
6. 支持批量操作

通过这个完整的示例，你可以快速了解和使用 GoBatis Example 功能，构建高效、安全的数据访问层。