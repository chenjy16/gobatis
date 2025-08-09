package examples

import (
	"time"
)

// User 用户实体
type User struct {
	ID          int64      `json:"id" db:"id"`
	Username    string     `json:"username" db:"username"`
	Email       string     `json:"email" db:"email"`
	Password    string     `json:"password" db:"password"`
	RealName    *string    `json:"real_name" db:"real_name"`
	Age         *int       `json:"age" db:"age"`
	Gender      *string    `json:"gender" db:"gender"`
	Phone       *string    `json:"phone" db:"phone"`
	Department  *string    `json:"department" db:"department"`
	Position    *string    `json:"position" db:"position"`
	Salary      *float64   `json:"salary" db:"salary"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
}

// Department 部门实体
type Department struct {
	ID          int64      `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	ManagerID   *int64     `json:"manager_id" db:"manager_id"`
	Budget      *float64   `json:"budget" db:"budget"`
	Status      string     `json:"status" db:"status"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// UserRole 用户角色实体
type UserRole struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	RoleName  string    `json:"role_name" db:"role_name"`
	GrantedAt time.Time `json:"granted_at" db:"granted_at"`
	GrantedBy *int64    `json:"granted_by" db:"granted_by"`
}

// UserLoginLog 用户登录日志实体
type UserLoginLog struct {
	ID            int64     `json:"id" db:"id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	LoginTime     time.Time `json:"login_time" db:"login_time"`
	IPAddress     *string   `json:"ip_address" db:"ip_address"`
	UserAgent     *string   `json:"user_agent" db:"user_agent"`
	LoginResult   string    `json:"login_result" db:"login_result"`
	FailureReason *string   `json:"failure_reason" db:"failure_reason"`
}

// UserWithDepartment 用户和部门的联合查询结果
type UserWithDepartment struct {
	User
	DepartmentName        *string  `json:"department_name" db:"department_name"`
	DepartmentDescription *string  `json:"department_description" db:"department_description"`
	DepartmentBudget      *float64 `json:"department_budget" db:"department_budget"`
}

// UserStatistics 用户统计信息
type UserStatistics struct {
	Department  string  `json:"department" db:"department"`
	UserCount   int     `json:"user_count" db:"user_count"`
	AvgAge      float64 `json:"avg_age" db:"avg_age"`
	AvgSalary   float64 `json:"avg_salary" db:"avg_salary"`
	MinSalary   float64 `json:"min_salary" db:"min_salary"`
	MaxSalary   float64 `json:"max_salary" db:"max_salary"`
	ActiveCount int     `json:"active_count" db:"active_count"`
}

// DepartmentStatistics 部门统计信息
type DepartmentStatistics struct {
	DepartmentName string    `json:"department_name" db:"department_name"`
	TotalUsers     int       `json:"total_users" db:"total_users"`
	ActiveUsers    int       `json:"active_users" db:"active_users"`
	AvgAge         float64   `json:"avg_age" db:"avg_age"`
	TotalSalary    float64   `json:"total_salary" db:"total_salary"`
	AvgSalary      float64   `json:"avg_salary" db:"avg_salary"`
	EarliestUser   time.Time `json:"earliest_user" db:"earliest_user"`
	LatestUser     time.Time `json:"latest_user" db:"latest_user"`
}

// SearchCondition 搜索条件
type SearchCondition struct {
	Username   *string `json:"username"`
	Email      *string `json:"email"`
	RealName   *string `json:"real_name"`
	Department *string `json:"department"`
	Position   *string `json:"position"`
	Status     *string `json:"status"`
	MinAge     *int    `json:"min_age"`
	MaxAge     *int    `json:"max_age"`
	MinSalary  *float64 `json:"min_salary"`
	MaxSalary  *float64 `json:"max_salary"`
	Gender     *string `json:"gender"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
}

// PageRequest 分页请求
type PageRequest struct {
	Page     int    `json:"page"`     // 页码，从1开始
	PageSize int    `json:"page_size"` // 每页大小
	OrderBy  string `json:"order_by"`  // 排序字段
	Order    string `json:"order"`     // 排序方向：ASC/DESC
}

// PageResponse 分页响应
type PageResponse struct {
	Total    int64       `json:"total"`     // 总记录数
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页大小
	Pages    int         `json:"pages"`     // 总页数
	Data     interface{} `json:"data"`      // 数据列表
}

// DefaultPageRequest 创建默认分页请求
func DefaultPageRequest() *PageRequest {
	return &PageRequest{
		Page:     1,
		PageSize: 10,
		OrderBy:  "id",
		Order:    "DESC",
	}
}

// GetOffset 计算偏移量
func (p *PageRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数量
func (p *PageRequest) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

// CalculatePages 计算总页数
func (p *PageRequest) CalculatePages(total int64) int {
	if p.PageSize <= 0 {
		return 0
	}
	pages := int(total) / p.PageSize
	if int(total)%p.PageSize > 0 {
		pages++
	}
	return pages
}