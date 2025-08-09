# gobatis 分页功能使用指南

本指南详细介绍了如何在 gobatis 框架中使用分页功能，包括配置、使用方法和最佳实践。

## 📋 目录

- [分页结构定义](#分页结构定义)
- [Mapper 接口定义](#mapper-接口定义)
- [基本使用示例](#基本使用示例)
- [高级使用场景](#高级使用场景)
- [最佳实践](#最佳实践)
- [常见问题](#常见问题)

## 🏗️ 分页结构定义

### PageRequest - 分页请求参数

```go
// PageRequest 分页请求参数
type PageRequest struct {
    Page     int    `json:"page"`     // 页码（从1开始）
    Size     int    `json:"size"`     // 每页大小
    Offset   int    `json:"offset"`   // 偏移量（自动计算）
    SortBy   string `json:"sortBy"`   // 排序字段
    SortDir  string `json:"sortDir"`  // 排序方向（ASC/DESC）
}
```

**字段说明：**
- `Page`: 页码，从1开始计数
- `Size`: 每页显示的记录数
- `Offset`: 数据库查询偏移量，通常由框架自动计算
- `SortBy`: 排序字段名
- `SortDir`: 排序方向，支持 "ASC"（升序）和 "DESC"（降序）

### PageResult - 分页结果

```go
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

**字段说明：**
- `Data`: 当前页的数据列表，类型为 `interface{}`，需要进行类型断言
- `Total`: 符合条件的总记录数
- `Page`: 当前页码
- `Size`: 每页大小
- `TotalPages`: 总页数
- `HasNext`: 是否存在下一页
- `HasPrev`: 是否存在上一页

## 🔧 Mapper 接口定义

### 方式1：直接使用 PageRequest

```go
type UserMapper interface {
    // 分页查询所有用户
    FindAllWithPage(pageReq *PageRequest) (*PageResult, error)
    
    // 按条件分页查询
    FindByNameWithPage(name string, pageReq *PageRequest) (*PageResult, error)
}
```

### 方式2：使用包含分页信息的结构体

```go
// 用户搜索条件（包含分页信息）
type UserSearchCondition struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Page     int    `json:"page"`     // 分页插件会自动识别
    Size     int    `json:"size"`     // 分页插件会自动识别
}

type UserMapper interface {
    // 使用条件结构体进行分页查询
    FindByCondition(condition *UserSearchCondition) (*PageResult, error)
}
```

## 📝 基本使用示例

### 示例1：基本分页查询

```go
func basicPaginationExample() {
    // 创建分页请求
    pageRequest := &PageRequest{
        Page:    1,      // 第1页
        Size:    10,     // 每页10条
        SortBy:  "id",   // 按ID排序
        SortDir: "DESC", // 降序
    }

    // 执行分页查询
    result, err := userMapper.FindAllWithPage(pageRequest)
    if err != nil {
        log.Fatal(err)
    }

    // 处理分页结果
    fmt.Printf("总记录数: %d\n", result.Total)
    fmt.Printf("当前页: %d/%d\n", result.Page, result.TotalPages)
    fmt.Printf("是否有下一页: %t\n", result.HasNext)

    // 获取数据（需要类型断言）
    users := result.Data.([]*User)
    for _, user := range users {
        fmt.Printf("用户: %s\n", user.Username)
    }
}
```

### 示例2：条件查询 + 分页

```go
func conditionPaginationExample() {
    // 创建搜索条件
    condition := &UserSearchCondition{
        Username: "john",
        Page:     2,
        Size:     5,
    }

    // 执行条件分页查询
    result, err := userMapper.FindByCondition(condition)
    if err != nil {
        log.Fatal(err)
    }

    // 处理结果
    users := result.Data.([]*User)
    for _, user := range users {
        fmt.Printf("匹配用户: %s\n", user.Username)
    }
}
```

### 示例3：分页结果详细处理

```go
func handlePaginationResult(result *PageResult) {
    // 分页信息
    fmt.Printf("=== 分页信息 ===\n")
    fmt.Printf("总记录数: %d\n", result.Total)
    fmt.Printf("当前页: %d/%d\n", result.Page, result.TotalPages)
    fmt.Printf("每页大小: %d\n", result.Size)
    
    // 导航信息
    fmt.Printf("=== 导航信息 ===\n")
    if result.HasPrev {
        fmt.Printf("上一页: %d\n", result.Page-1)
    }
    if result.HasNext {
        fmt.Printf("下一页: %d\n", result.Page+1)
    }
    
    // 记录范围
    startRecord := (result.Page-1)*result.Size + 1
    endRecord := startRecord + len(result.Data.([]*User)) - 1
    fmt.Printf("当前页记录范围: %d - %d\n", startRecord, endRecord)
    
    // 数据处理
    users := result.Data.([]*User)
    for i, user := range users {
        fmt.Printf("%d. %s\n", startRecord+i, user.Username)
    }
}
```

## 🚀 高级使用场景

### 场景1：多条件复合查询分页

```go
type AdvancedUserQuery struct {
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    MinAge      int       `json:"minAge"`
    MaxAge      int       `json:"maxAge"`
    CreateStart time.Time `json:"createStart"`
    CreateEnd   time.Time `json:"createEnd"`
    Status      string    `json:"status"`
    Page        int       `json:"page"`
    Size        int       `json:"size"`
}

func advancedQueryExample() {
    query := &AdvancedUserQuery{
        Username:    "john%",
        MinAge:      18,
        MaxAge:      65,
        Status:      "active",
        CreateStart: time.Now().AddDate(0, -1, 0), // 一个月前
        CreateEnd:   time.Now(),
        Page:        1,
        Size:        20,
    }

    result, err := userMapper.FindByAdvancedCondition(query)
    // 处理结果...
}
```

### 场景2：动态排序分页

```go
func dynamicSortExample() {
    // 支持多种排序方式
    sortOptions := []struct {
        Field string
        Dir   string
    }{
        {"create_time", "DESC"},  // 按创建时间降序
        {"username", "ASC"},      // 按用户名升序
        {"age", "DESC"},          // 按年龄降序
    }

    for _, sort := range sortOptions {
        pageRequest := &PageRequest{
            Page:    1,
            Size:    10,
            SortBy:  sort.Field,
            SortDir: sort.Dir,
        }

        result, err := userMapper.FindAllWithPage(pageRequest)
        if err != nil {
            continue
        }

        fmt.Printf("按 %s %s 排序的结果:\n", sort.Field, sort.Dir)
        // 处理结果...
    }
}
```

### 场景3：分页数据导出

```go
func exportAllData() {
    const pageSize = 1000
    page := 1
    var allUsers []*User

    for {
        pageRequest := &PageRequest{
            Page: page,
            Size: pageSize,
        }

        result, err := userMapper.FindAllWithPage(pageRequest)
        if err != nil {
            log.Printf("导出第%d页数据失败: %v", page, err)
            break
        }

        users := result.Data.([]*User)
        allUsers = append(allUsers, users...)

        fmt.Printf("已导出第%d页，共%d条记录\n", page, len(users))

        // 检查是否还有下一页
        if !result.HasNext {
            break
        }

        page++
    }

    fmt.Printf("导出完成，总计%d条记录\n", len(allUsers))
    // 处理导出逻辑...
}
```

## 💡 最佳实践

### 1. 分页参数验证

```go
func validatePageRequest(req *PageRequest) error {
    if req.Page < 1 {
        req.Page = 1
    }
    if req.Size < 1 {
        req.Size = 10
    }
    if req.Size > 1000 {
        req.Size = 1000 // 限制最大页面大小
    }
    if req.SortDir != "" && req.SortDir != "ASC" && req.SortDir != "DESC" {
        req.SortDir = "ASC"
    }
    return nil
}
```

### 2. 分页响应包装

```go
type PaginationResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    *PageResult `json:"data"`
}

func buildPaginationResponse(result *PageResult) *PaginationResponse {
    return &PaginationResponse{
        Code:    200,
        Message: "success",
        Data:    result,
    }
}
```

### 3. 分页缓存策略

```go
func getCachedPageData(cacheKey string, pageReq *PageRequest) (*PageResult, bool) {
    // 实现缓存逻辑
    // 注意：缓存键应该包含查询条件和分页参数
    key := fmt.Sprintf("%s:page:%d:size:%d", cacheKey, pageReq.Page, pageReq.Size)
    // 从缓存获取数据...
    return nil, false
}
```

### 4. 性能优化建议

```go
// 对于大数据量的分页查询，建议：
// 1. 使用索引优化排序字段
// 2. 避免使用 OFFSET 进行深度分页
// 3. 考虑使用游标分页（cursor-based pagination）

type CursorPageRequest struct {
    Cursor string `json:"cursor"` // 游标值
    Size   int    `json:"size"`   // 每页大小
}

type CursorPageResult struct {
    Data       interface{} `json:"data"`
    NextCursor string      `json:"nextCursor"`
    HasNext    bool        `json:"hasNext"`
}
```

## ❓ 常见问题

### Q1: 如何处理空的分页结果？

```go
func handleEmptyResult(result *PageResult) {
    if result.Total == 0 {
        fmt.Println("没有找到符合条件的记录")
        return
    }

    users := result.Data.([]*User)
    if len(users) == 0 {
        fmt.Println("当前页没有数据")
        return
    }

    // 正常处理数据...
}
```

### Q2: 如何实现前端分页组件的数据绑定？

```go
type PaginationInfo struct {
    CurrentPage int   `json:"currentPage"`
    PageSize    int   `json:"pageSize"`
    Total       int64 `json:"total"`
    TotalPages  int   `json:"totalPages"`
    HasPrev     bool  `json:"hasPrev"`
    HasNext     bool  `json:"hasNext"`
    StartRecord int   `json:"startRecord"`
    EndRecord   int   `json:"endRecord"`
}

func buildPaginationInfo(result *PageResult) *PaginationInfo {
    startRecord := (result.Page-1)*result.Size + 1
    endRecord := startRecord + len(result.Data.([]*User)) - 1

    return &PaginationInfo{
        CurrentPage: result.Page,
        PageSize:    result.Size,
        Total:       result.Total,
        TotalPages:  result.TotalPages,
        HasPrev:     result.HasPrev,
        HasNext:     result.HasNext,
        StartRecord: startRecord,
        EndRecord:   endRecord,
    }
}
```

### Q3: 如何处理分页查询的异常情况？

```go
func safePaginationQuery(pageReq *PageRequest) (*PageResult, error) {
    // 参数验证
    if err := validatePageRequest(pageReq); err != nil {
        return nil, err
    }

    // 执行查询
    result, err := userMapper.FindAllWithPage(pageReq)
    if err != nil {
        log.Printf("分页查询失败: %v", err)
        return nil, err
    }

    // 结果验证
    if result == nil {
        return &PageResult{
            Data:       []*User{},
            Total:      0,
            Page:       pageReq.Page,
            Size:       pageReq.Size,
            TotalPages: 0,
            HasNext:    false,
            HasPrev:    false,
        }, nil
    }

    return result, nil
}
```

## 🔗 相关链接

- [gobatis 主要文档](../README.md)
- [插件系统文档](PLUGIN_GUIDE.md)
- [配置指南](CONFIGURATION_GUIDE.md)
- [示例代码](pagination_example.go)

---

**注意**: 本指南基于 gobatis 框架的分页插件实现。在实际使用中，请根据具体的数据库类型和业务需求进行相应的调整。