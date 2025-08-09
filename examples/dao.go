package examples

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gobatis/core/example"
)

// UserDAO 用户数据访问对象
type UserDAO struct {
	db *sql.DB
}

// NewUserDAO 创建用户DAO
func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{db: db}
}

// SelectByExample 根据Example查询用户列表
func (dao *UserDAO) SelectByExample(ex *example.Example) ([]*User, error) {
	sqlStr, args := ex.BuildSQL("SELECT * FROM users")
	
	// BuildSQL 已经包含了完整的 SQL 语句
	
	rows, err := dao.db.Query(sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	defer rows.Close()
	
	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Password,
			&user.RealName, &user.Age, &user.Gender, &user.Phone,
			&user.Department, &user.Position, &user.Salary, &user.Status,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描用户数据失败: %w", err)
		}
		users = append(users, user)
	}
	
	return users, nil
}

// CountByExample 根据Example统计用户数量
func (dao *UserDAO) CountByExample(ex *example.Example) (int64, error) {
	// 构建COUNT查询
	sqlStr, args := ex.BuildSQL("SELECT COUNT(*) FROM users")
	
	// 移除ORDER BY和LIMIT子句（COUNT查询不需要）
	if idx := strings.Index(strings.ToUpper(sqlStr), "ORDER BY"); idx != -1 {
		sqlStr = sqlStr[:idx]
	}
	if idx := strings.Index(strings.ToUpper(sqlStr), "LIMIT"); idx != -1 {
		sqlStr = sqlStr[:idx]
	}
	
	var count int64
	err := dao.db.QueryRow(sqlStr, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("统计用户数量失败: %w", err)
	}
	
	return count, nil
}

// SelectByPrimaryKey 根据主键查询用户
func (dao *UserDAO) SelectByPrimaryKey(id int64) (*User, error) {
	ex := example.NewExample()
	ex.CreateCriteria().AndEqualTo("id", id)
	
	users, err := dao.SelectByExample(ex)
	if err != nil {
		return nil, err
	}
	
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}
	
	return users[0], nil
}

// Insert 插入用户
func (dao *UserDAO) Insert(user *User) error {
	sqlStr := `
		INSERT INTO users (username, email, password, real_name, age, gender, 
		                  phone, department, position, salary, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := dao.db.Exec(sqlStr,
		user.Username, user.Email, user.Password, user.RealName,
		user.Age, user.Gender, user.Phone, user.Department,
		user.Position, user.Salary, user.Status,
	)
	if err != nil {
		return fmt.Errorf("插入用户失败: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取插入ID失败: %w", err)
	}
	
	user.ID = id
	return nil
}

// UpdateByExample 根据Example更新用户
func (dao *UserDAO) UpdateByExample(user *User, ex *example.Example) (int64, error) {
	// 构建UPDATE语句
	updateSQL := `
		UPDATE users SET 
			username = ?, email = ?, real_name = ?, age = ?, gender = ?,
			phone = ?, department = ?, position = ?, salary = ?, status = ?,
			updated_at = CURRENT_TIMESTAMP
	`
	
	updateArgs := []interface{}{
		user.Username, user.Email, user.RealName, user.Age, user.Gender,
		user.Phone, user.Department, user.Position, user.Salary, user.Status,
	}
	
	// 获取WHERE条件
	whereSQL, args := ex.BuildSQL("")
	if strings.Contains(strings.ToUpper(whereSQL), "WHERE") {
		updateSQL += " " + whereSQL
		updateArgs = append(updateArgs, args...)
	}
	
	result, err := dao.db.Exec(updateSQL, updateArgs...)
	if err != nil {
		return 0, fmt.Errorf("更新用户失败: %w", err)
	}
	
	return result.RowsAffected()
}

// DeleteByExample 根据Example删除用户
func (dao *UserDAO) DeleteByExample(ex *example.Example) (int64, error) {
	whereSQL, args := ex.BuildSQL("")
	
	deleteSQL := "DELETE FROM users"
	if strings.Contains(strings.ToUpper(whereSQL), "WHERE") {
		deleteSQL += " " + whereSQL
	}
	
	result, err := dao.db.Exec(deleteSQL, args...)
	if err != nil {
		return 0, fmt.Errorf("删除用户失败: %w", err)
	}
	
	return result.RowsAffected()
}

// SelectWithDepartment 查询用户及其部门信息
func (dao *UserDAO) SelectWithDepartment(ex *example.Example) ([]*UserWithDepartment, error) {
	baseSQL := `
		SELECT u.*, d.name as department_name, d.description as department_description, d.budget as department_budget
		FROM users u
		LEFT JOIN departments d ON u.department = d.name
	`
	
	joinSQL, args := ex.BuildSQL(baseSQL)
	
	rows, err := dao.db.Query(joinSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("查询用户部门信息失败: %w", err)
	}
	defer rows.Close()
	
	var results []*UserWithDepartment
	for rows.Next() {
		result := &UserWithDepartment{}
		err := rows.Scan(
			&result.ID, &result.Username, &result.Email, &result.Password,
			&result.RealName, &result.Age, &result.Gender, &result.Phone,
			&result.Department, &result.Position, &result.Salary, &result.Status,
			&result.CreatedAt, &result.UpdatedAt, &result.LastLoginAt,
			&result.DepartmentName, &result.DepartmentDescription, &result.DepartmentBudget,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描用户部门数据失败: %w", err)
		}
		results = append(results, result)
	}
	
	return results, nil
}

// GetUserStatistics 获取用户统计信息
func (dao *UserDAO) GetUserStatistics() ([]*UserStatistics, error) {
	sqlStr := `
		SELECT 
			COALESCE(department, '未分配') as department,
			COUNT(*) as user_count,
			COALESCE(AVG(age), 0) as avg_age,
			COALESCE(AVG(salary), 0) as avg_salary,
			COALESCE(MIN(salary), 0) as min_salary,
			COALESCE(MAX(salary), 0) as max_salary,
			SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END) as active_count
		FROM users
		GROUP BY department
		ORDER BY user_count DESC
	`
	
	rows, err := dao.db.Query(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("查询用户统计信息失败: %w", err)
	}
	defer rows.Close()
	
	var stats []*UserStatistics
	for rows.Next() {
		stat := &UserStatistics{}
		err := rows.Scan(
			&stat.Department, &stat.UserCount, &stat.AvgAge,
			&stat.AvgSalary, &stat.MinSalary, &stat.MaxSalary, &stat.ActiveCount,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描统计数据失败: %w", err)
		}
		stats = append(stats, stat)
	}
	
	return stats, nil
}

// SearchUsers 动态搜索用户
func (dao *UserDAO) SearchUsers(condition *SearchCondition, pageReq *PageRequest) (*PageResponse, error) {
	ex := example.NewExample()
	criteria := ex.CreateCriteria()
	
	// 动态构建查询条件
	if condition.Username != nil && *condition.Username != "" {
		criteria.AndLike("username", "%"+*condition.Username+"%")
	}
	if condition.Email != nil && *condition.Email != "" {
		criteria.AndLike("email", "%"+*condition.Email+"%")
	}
	if condition.RealName != nil && *condition.RealName != "" {
		criteria.AndLike("real_name", "%"+*condition.RealName+"%")
	}
	if condition.Department != nil && *condition.Department != "" {
		criteria.AndEqualTo("department", *condition.Department)
	}
	if condition.Position != nil && *condition.Position != "" {
		criteria.AndLike("position", "%"+*condition.Position+"%")
	}
	if condition.Status != nil && *condition.Status != "" {
		criteria.AndEqualTo("status", *condition.Status)
	}
	if condition.Gender != nil && *condition.Gender != "" {
		criteria.AndEqualTo("gender", *condition.Gender)
	}
	if condition.MinAge != nil {
		criteria.AndGreaterThanOrEqualTo("age", *condition.MinAge)
	}
	if condition.MaxAge != nil {
		criteria.AndLessThanOrEqualTo("age", *condition.MaxAge)
	}
	if condition.MinSalary != nil {
		criteria.AndGreaterThanOrEqualTo("salary", *condition.MinSalary)
	}
	if condition.MaxSalary != nil {
		criteria.AndLessThanOrEqualTo("salary", *condition.MaxSalary)
	}
	if condition.StartDate != nil {
		criteria.AndGreaterThanOrEqualTo("created_at", *condition.StartDate)
	}
	if condition.EndDate != nil {
		criteria.AndLessThanOrEqualTo("created_at", *condition.EndDate)
	}
	
	// 设置排序
	if pageReq.OrderBy != "" {
		if strings.ToUpper(pageReq.Order) == "DESC" {
			ex.SetOrderByClause(pageReq.OrderBy + " DESC")
		} else {
			ex.SetOrderByClause(pageReq.OrderBy + " ASC")
		}
	}
	
	// 先查询总数
	total, err := dao.CountByExample(ex)
	if err != nil {
		return nil, err
	}
	
	// 设置分页
	ex.SetLimit(pageReq.GetOffset(), pageReq.GetLimit())
	
	// 查询数据
	users, err := dao.SelectByExample(ex)
	if err != nil {
		return nil, err
	}
	
	return &PageResponse{
		Total:    total,
		Page:     pageReq.Page,
		PageSize: pageReq.PageSize,
		Pages:    pageReq.CalculatePages(total),
		Data:     users,
	}, nil
}

// GetActiveUsersInDepartments 获取指定部门的活跃用户
func (dao *UserDAO) GetActiveUsersInDepartments(departments []string) ([]*User, error) {
	ex := example.NewExample()
	criteria := ex.CreateCriteria()
	
	criteria.AndEqualTo("status", "active")
	if len(departments) > 0 {
		// 转换 []string 为 []interface{}
		deptInterfaces := make([]interface{}, len(departments))
		for i, dept := range departments {
			deptInterfaces[i] = dept
		}
		criteria.AndIn("department", deptInterfaces)
	}
	criteria.AndIsNotNull("last_login_at")
	
	// 按最后登录时间倒序
	ex.SetOrderByClause("last_login_at DESC")
	
	return dao.SelectByExample(ex)
}

// GetUsersByAgeRange 获取指定年龄范围的用户
func (dao *UserDAO) GetUsersByAgeRange(minAge, maxAge int) ([]*User, error) {
	ex := example.NewExample()
	criteria := ex.CreateCriteria()
	
	criteria.AndBetween("age", minAge, maxAge)
	criteria.AndEqualTo("status", "active")
	
	// 按年龄升序
	ex.SetOrderByClause("age ASC, real_name ASC")
	
	return dao.SelectByExample(ex)
}

// GetHighSalaryUsers 获取高薪用户（薪资大于指定值）
func (dao *UserDAO) GetHighSalaryUsers(minSalary float64) ([]*User, error) {
	ex := example.NewExample()
	criteria := ex.CreateCriteria()
	
	criteria.AndGreaterThan("salary", minSalary)
	criteria.AndEqualTo("status", "active")
	
	// 按薪资倒序
	ex.SetOrderByClause("salary DESC")
	
	return dao.SelectByExample(ex)
}

// GetRecentUsers 获取最近注册的用户
func (dao *UserDAO) GetRecentUsers(days int) ([]*User, error) {
	ex := example.NewExample()
	criteria := ex.CreateCriteria()
	
	// 计算日期
	since := time.Now().AddDate(0, 0, -days)
	criteria.AndGreaterThanOrEqualTo("created_at", since)
	
	// 按创建时间倒序
	ex.SetOrderByClause("created_at DESC")
	ex.SetLimit(0, 50) // 限制50条
	
	return dao.SelectByExample(ex)
}

// GetUsersWithComplexCondition 复杂条件查询示例
func (dao *UserDAO) GetUsersWithComplexCondition() ([]*User, error) {
	ex := example.NewExample()
	
	// 第一组条件：高级开发人员
	criteria1 := ex.CreateCriteria()
	criteria1.AndEqualTo("department", "IT")
	criteria1.AndLike("position", "%高级%")
	criteria1.AndGreaterThan("salary", 10000)
	
	// 第二组条件：管理人员
	criteria2 := ex.CreateCriteria()
	criteria2.AndLike("position", "%经理%")
	criteria2.AndLike("position", "%主管%")
	criteria2.AndEqualTo("status", "active")
	ex.Or(*criteria2)
	
	// 第三组条件：新员工但薪资较高
	criteria3 := ex.CreateCriteria()
	criteria3.AndGreaterThanOrEqualTo("created_at", time.Now().AddDate(0, -6, 0)) // 6个月内
	criteria3.AndGreaterThan("salary", 8000)
	criteria3.AndEqualTo("status", "active")
	ex.Or(*criteria3)
	
	ex.SetOrderByClause("salary DESC, created_at DESC")
	
	return dao.SelectByExample(ex)
}