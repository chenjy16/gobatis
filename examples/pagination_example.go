package examples

import (
	"fmt"
	"log"
	"time"
)

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

// UserSearchCondition 用户搜索条件（包含分页信息）
type UserSearchCondition struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Page     int    `json:"page"` // 分页插件会自动识别
	Size     int    `json:"size"` // 分页插件会自动识别
}

// UserWithAge 扩展的用户结构体（包含年龄字段用于演示）
type UserWithAge struct {
	ID       int64     `db:"id"`
	Username string    `db:"username"`
	Email    string    `db:"email"`
	Age      int       `db:"age"`
	CreateAt time.Time `db:"create_at"`
}

// UserMapperWithPagination 支持分页的用户Mapper接口
type UserMapperWithPagination interface {
	// 普通查询
	FindAll() ([]*UserWithAge, error)

	// 分页查询 - 方式1：直接传入 PageRequest
	FindAllWithPage(pageReq *PageRequest) (*PageResult, error)

	// 分页查询 - 方式2：传入包含分页信息的结构体
	FindByCondition(condition *UserSearchCondition) (*PageResult, error)

	// 分页查询 - 方式3：按年龄范围分页查询
	FindByAgeRange(minAge, maxAge int, pageReq *PageRequest) (*PageResult, error)
}

// PaginationExamples 分页示例演示
func PaginationExamples() {
	fmt.Println("=== 分页功能示例演示 ===")

	// 创建模拟的UserMapper（在实际应用中，这将通过gobatis框架创建）
	userMapper := &MockUserMapperWithPagination{}

	// 示例1: 基本分页查询
	fmt.Println("\n1. 基本分页查询示例:")
	basicPaginationExample(userMapper)

	// 示例2: 条件查询 + 分页
	fmt.Println("\n2. 条件查询 + 分页示例:")
	conditionPaginationExample(userMapper)

	// 示例3: 年龄范围查询 + 分页
	fmt.Println("\n3. 年龄范围查询 + 分页示例:")
	ageRangePaginationExample(userMapper)

	// 示例4: 分页结果处理
	fmt.Println("\n4. 分页结果详细处理示例:")
	paginationResultHandlingExample(userMapper)
}

// basicPaginationExample 基本分页查询示例
func basicPaginationExample(userMapper UserMapperWithPagination) {
	pageRequest := &PageRequest{
		Page:    1,      // 第1页
		Size:    10,     // 每页10条
		SortBy:  "id",   // 按ID排序
		SortDir: "DESC", // 降序
	}

	result, err := userMapper.FindAllWithPage(pageRequest)
	if err != nil {
		log.Printf("分页查询失败: %v", err)
		return
	}

	// 处理分页结果
	fmt.Printf("总记录数: %d\n", result.Total)
	fmt.Printf("当前页: %d/%d\n", result.Page, result.TotalPages)
	fmt.Printf("是否有下一页: %t\n", result.HasNext)
	fmt.Printf("是否有上一页: %t\n", result.HasPrev)

	// 获取数据
	users := result.Data.([]*UserWithAge)
	fmt.Printf("当前页用户数量: %d\n", len(users))
	for i, user := range users {
		if i < 3 { // 只显示前3个用户
			fmt.Printf("  用户%d: %s (邮箱: %s, 年龄: %d)\n", i+1, user.Username, user.Email, user.Age)
		}
	}
	if len(users) > 3 {
		fmt.Printf("  ... 还有 %d 个用户\n", len(users)-3)
	}
}

// conditionPaginationExample 条件查询 + 分页示例
func conditionPaginationExample(userMapper UserMapperWithPagination) {
	condition := &UserSearchCondition{
		Username: "john",
		Page:     2,
		Size:     5,
	}

	result, err := userMapper.FindByCondition(condition)
	if err != nil {
		log.Printf("条件分页查询失败: %v", err)
		return
	}

	fmt.Printf("搜索条件: 用户名包含 '%s'\n", condition.Username)
	fmt.Printf("查询结果: 第%d页，每页%d条，共%d条记录\n", 
		result.Page, result.Size, result.Total)

	users := result.Data.([]*UserWithAge)
	for _, user := range users {
		fmt.Printf("  匹配用户: %s\n", user.Username)
	}
}

// ageRangePaginationExample 年龄范围查询 + 分页示例
func ageRangePaginationExample(userMapper UserMapperWithPagination) {
	pageRequest := &PageRequest{
		Page:    1,
		Size:    8,
		SortBy:  "age",
		SortDir: "ASC",
	}

	result, err := userMapper.FindByAgeRange(25, 35, pageRequest)
	if err != nil {
		log.Printf("年龄范围分页查询失败: %v", err)
		return
	}

	fmt.Printf("年龄范围: 25-35岁\n")
	fmt.Printf("查询结果: 共%d条记录，当前第%d页\n", result.Total, result.Page)

	users := result.Data.([]*UserWithAge)
	for _, user := range users {
		fmt.Printf("  用户: %s (年龄: %d)\n", user.Username, user.Age)
	}
}

// paginationResultHandlingExample 分页结果详细处理示例
func paginationResultHandlingExample(userMapper UserMapperWithPagination) {
	pageRequest := &PageRequest{
		Page: 1,
		Size: 5,
	}

	result, err := userMapper.FindAllWithPage(pageRequest)
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}

	// 详细的分页信息处理
	fmt.Printf("=== 分页详细信息 ===\n")
	fmt.Printf("总记录数: %d\n", result.Total)
	fmt.Printf("当前页码: %d\n", result.Page)
	fmt.Printf("每页大小: %d\n", result.Size)
	fmt.Printf("总页数: %d\n", result.TotalPages)
	fmt.Printf("是否有上一页: %t\n", result.HasPrev)
	fmt.Printf("是否有下一页: %t\n", result.HasNext)

	// 计算分页导航信息
	fmt.Printf("\n=== 分页导航信息 ===\n")
	if result.HasPrev {
		fmt.Printf("上一页: %d\n", result.Page-1)
	}
	fmt.Printf("当前页: %d\n", result.Page)
	if result.HasNext {
		fmt.Printf("下一页: %d\n", result.Page+1)
	}

	// 计算记录范围
	startRecord := (result.Page-1)*result.Size + 1
	endRecord := startRecord + len(result.Data.([]*UserWithAge)) - 1
	fmt.Printf("当前页记录范围: %d - %d\n", startRecord, endRecord)

	// 数据处理示例
	users := result.Data.([]*UserWithAge)
	fmt.Printf("\n=== 当前页数据 ===\n")
	for i, user := range users {
		fmt.Printf("%d. %s (ID: %d, 邮箱: %s)\n", 
			startRecord+i, user.Username, user.ID, user.Email)
	}
}

// MockUserMapperWithPagination 模拟支持分页的UserMapper实现
type MockUserMapperWithPagination struct{}

func (m *MockUserMapperWithPagination) FindAll() ([]*UserWithAge, error) {
	return generateMockUsers(50), nil
}

func (m *MockUserMapperWithPagination) FindAllWithPage(pageReq *PageRequest) (*PageResult, error) {
	allUsers := generateMockUsers(50)
	return paginateUsers(allUsers, pageReq.Page, pageReq.Size), nil
}

func (m *MockUserMapperWithPagination) FindByCondition(condition *UserSearchCondition) (*PageResult, error) {
	allUsers := generateMockUsers(50)
	
	// 模拟条件过滤
	var filteredUsers []*UserWithAge
	for _, user := range allUsers {
		if condition.Username == "" || 
		   fmt.Sprintf("%s", user.Username) == condition.Username ||
		   fmt.Sprintf("%s", user.Username)[:4] == condition.Username {
			filteredUsers = append(filteredUsers, user)
		}
	}
	
	return paginateUsers(filteredUsers, condition.Page, condition.Size), nil
}

func (m *MockUserMapperWithPagination) FindByAgeRange(minAge, maxAge int, pageReq *PageRequest) (*PageResult, error) {
	allUsers := generateMockUsers(50)
	
	// 模拟年龄范围过滤
	var filteredUsers []*UserWithAge
	for _, user := range allUsers {
		if user.Age >= minAge && user.Age <= maxAge {
			filteredUsers = append(filteredUsers, user)
		}
	}
	
	return paginateUsers(filteredUsers, pageReq.Page, pageReq.Size), nil
}

// generateMockUsers 生成模拟用户数据
func generateMockUsers(count int) []*UserWithAge {
	users := make([]*UserWithAge, count)
	names := []string{"john", "jane", "bob", "alice", "charlie", "diana", "eve", "frank"}
	
	for i := 0; i < count; i++ {
		users[i] = &UserWithAge{
			ID:       int64(i + 1),
			Username: fmt.Sprintf("%s%d", names[i%len(names)], i+1),
			Email:    fmt.Sprintf("%s%d@example.com", names[i%len(names)], i+1),
			Age:      20 + (i % 40), // 年龄在20-59之间
			CreateAt: time.Now(),
		}
	}
	
	return users
}

// paginateUsers 对用户列表进行分页处理
func paginateUsers(users []*UserWithAge, page, size int) *PageResult {
	total := int64(len(users))
	totalPages := int((total + int64(size) - 1) / int64(size))
	
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}
	
	start := (page - 1) * size
	end := start + size
	
	if start > len(users) {
		start = len(users)
	}
	if end > len(users) {
		end = len(users)
	}
	
	pageData := users[start:end]
	
	return &PageResult{
		Data:       pageData,
		Total:      total,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}