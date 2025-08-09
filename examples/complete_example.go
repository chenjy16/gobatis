package examples

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gobatis/core/example"
	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
)

// CompleteExample 完整的 Example 使用示例
type CompleteExample struct {
	db      *sql.DB
	userDAO *UserDAO
}

// NewCompleteExample 创建完整示例实例
func NewCompleteExample(dsn string) (*CompleteExample, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return &CompleteExample{
		db:      db,
		userDAO: NewUserDAO(db),
	}, nil
}

// Close 关闭数据库连接
func (ce *CompleteExample) Close() error {
	return ce.db.Close()
}

// RunAllExamples 运行所有示例
func (ce *CompleteExample) RunAllExamples() {
	fmt.Println("=== 完整的 Example 使用示例 ===")

	// 1. 基本查询示例
	ce.BasicQueryExamples()

	// 2. 复杂查询示例
	ce.ComplexQueryExamples()

	// 3. 分页查询示例
	ce.PaginationExamples()

	// 4. 统计查询示例
	ce.StatisticsExamples()

	// 5. 动态查询示例
	ce.DynamicQueryExamples()

	// 6. 联合查询示例
	ce.JoinQueryExamples()

	// 7. 数据操作示例
	ce.DataOperationExamples()
}

// BasicQueryExamples 基本查询示例
func (ce *CompleteExample) BasicQueryExamples() {
	fmt.Println("1. 基本查询示例")
	fmt.Println("================")

	// 1.1 查询所有活跃用户
	fmt.Println("1.1 查询所有活跃用户:")
	ex1 := example.NewExample()
	ex1.CreateCriteria().AndEqualTo("status", "active")
	ex1.SetOrderByClause("created_at DESC")

	users, err := ce.userDAO.SelectByExample(ex1)
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}
	fmt.Printf("找到 %d 个活跃用户\n", len(users))
	for _, user := range users {
		fmt.Printf("- %s (%s) - %s\n", *user.RealName, user.Username, *user.Department)
	}

	// 1.2 根据主键查询
	fmt.Println("\n1.2 根据主键查询用户:")
	user, err := ce.userDAO.SelectByPrimaryKey(1)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("用户: %s (%s)\n", *user.RealName, user.Username)
	}

	// 1.3 查询指定部门的用户
	fmt.Println("\n1.3 查询IT部门的用户:")
	ex2 := example.NewExample()
	ex2.CreateCriteria().
		AndEqualTo("department", "IT").
		AndEqualTo("status", "active")
	ex2.SetOrderByClause("salary DESC")

	itUsers, err := ce.userDAO.SelectByExample(ex2)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("IT部门有 %d 个活跃用户\n", len(itUsers))
		for _, user := range itUsers {
			fmt.Printf("- %s: %.2f\n", *user.RealName, *user.Salary)
		}
	}

	fmt.Println()
}

// ComplexQueryExamples 复杂查询示例
func (ce *CompleteExample) ComplexQueryExamples() {
	fmt.Println("2. 复杂查询示例")
	fmt.Println("================")

	// 2.1 年龄范围查询
	fmt.Println("2.1 查询25-35岁的用户:")
	ex1 := example.NewExample()
	ex1.CreateCriteria().
		AndBetween("age", 25, 35).
		AndEqualTo("status", "active")
	ex1.SetOrderByClause("age ASC")

	users, err := ce.userDAO.SelectByExample(ex1)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个25-35岁的用户\n", len(users))
		for _, user := range users {
			fmt.Printf("- %s: %d岁\n", *user.RealName, *user.Age)
		}
	}

	// 2.2 模糊查询
	fmt.Println("\n2.2 查询姓名包含'张'的用户:")
	ex2 := example.NewExample()
	ex2.CreateCriteria().
		AndLike("real_name", "%张%").
		AndEqualTo("status", "active")

	users, err = ce.userDAO.SelectByExample(ex2)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个姓名包含'张'的用户\n", len(users))
		for _, user := range users {
			fmt.Printf("- %s (%s)\n", *user.RealName, user.Username)
		}
	}

	// 2.3 IN 查询
	fmt.Println("\n2.3 查询多个部门的用户:")
	ex3 := example.NewExample()
	departments := []interface{}{"IT", "HR", "研发部"}
	ex3.CreateCriteria().
		AndIn("department", departments).
		AndEqualTo("status", "active")
	ex3.SetOrderByClause("department ASC, salary DESC")

	users, err = ce.userDAO.SelectByExample(ex3)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个指定部门的用户\n", len(users))
		for _, user := range users {
			fmt.Printf("- %s (%s): %.2f\n", *user.RealName, *user.Department, *user.Salary)
		}
	}

	// 2.4 NULL 值查询
	fmt.Println("\n2.4 查询没有电话号码的用户:")
	ex4 := example.NewExample()
	ex4.CreateCriteria().
		AndIsNull("phone").
		AndEqualTo("status", "active")

	users, err = ce.userDAO.SelectByExample(ex4)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个没有电话号码的用户\n", len(users))
		for _, user := range users {
			fmt.Printf("- %s (%s)\n", *user.RealName, user.Username)
		}
	}

	fmt.Println()
}

// PaginationExamples 分页查询示例
func (ce *CompleteExample) PaginationExamples() {
	fmt.Println("3. 分页查询示例")
	fmt.Println("================")

	// 3.1 基本分页
	fmt.Println("3.1 分页查询用户 (第1页，每页3条):")
	ex1 := example.NewExample()
	ex1.CreateCriteria().AndEqualTo("status", "active")
	ex1.SetOrderByClause("id ASC")
	ex1.SetLimit(0, 3) // 第1页，每页3条

	users, err := ce.userDAO.SelectByExample(ex1)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("第1页结果 (%d条):\n", len(users))
		for _, user := range users {
			fmt.Printf("- ID:%d %s (%s)\n", user.ID, *user.RealName, user.Username)
		}
	}

	// 3.2 第二页
	fmt.Println("\n3.2 分页查询用户 (第2页，每页3条):")
	ex2 := example.NewExample()
	ex2.CreateCriteria().AndEqualTo("status", "active")
	ex2.SetOrderByClause("id ASC")
	ex2.SetLimit(3, 3) // 第2页，每页3条

	users, err = ce.userDAO.SelectByExample(ex2)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("第2页结果 (%d条):\n", len(users))
		for _, user := range users {
			fmt.Printf("- ID:%d %s (%s)\n", user.ID, *user.RealName, user.Username)
		}
	}

	// 3.3 统计总数
	fmt.Println("\n3.3 统计活跃用户总数:")
	ex3 := example.NewExample()
	ex3.CreateCriteria().AndEqualTo("status", "active")

	count, err := ce.userDAO.CountByExample(ex3)
	if err != nil {
		log.Printf("统计失败: %v", err)
	} else {
		fmt.Printf("活跃用户总数: %d\n", count)
	}

	fmt.Println()
}

// StatisticsExamples 统计查询示例
func (ce *CompleteExample) StatisticsExamples() {
	fmt.Println("4. 统计查询示例")
	fmt.Println("================")

	// 4.1 用户统计
	fmt.Println("4.1 用户统计信息:")
	stats, err := ce.userDAO.GetUserStatistics()
	if err != nil {
		log.Printf("查询统计信息失败: %v", err)
	} else {
		fmt.Println("部门统计:")
		for _, stat := range stats {
			fmt.Printf("- %s: %d人, 平均年龄%.1f岁, 平均薪资%.2f, 活跃用户%d人\n",
				stat.Department, stat.UserCount, stat.AvgAge, stat.AvgSalary, stat.ActiveCount)
		}
	}

	// 4.2 高薪用户
	fmt.Println("\n4.2 查询高薪用户 (薪资>10000):")
	highSalaryUsers, err := ce.userDAO.GetHighSalaryUsers(10000)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个高薪用户:\n", len(highSalaryUsers))
		for _, user := range highSalaryUsers {
			fmt.Printf("- %s (%s): %.2f\n", *user.RealName, *user.Department, *user.Salary)
		}
	}

	// 4.3 最近注册用户
	fmt.Println("\n4.3 查询最近30天注册的用户:")
	recentUsers, err := ce.userDAO.GetRecentUsers(30)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 个最近注册的用户:\n", len(recentUsers))
		for _, user := range recentUsers {
			fmt.Printf("- %s: %s\n", *user.RealName, user.CreatedAt.Format("2006-01-02"))
		}
	}

	fmt.Println()
}

// DynamicQueryExamples 动态查询示例
func (ce *CompleteExample) DynamicQueryExamples() {
	fmt.Println("5. 动态查询示例")
	fmt.Println("================")

	// 5.1 动态搜索条件
	fmt.Println("5.1 动态用户搜索:")
	condition := &SearchCondition{
		Department: stringPtr("IT"),
		MinAge:     intPtr(25),
		MaxAge:     intPtr(40),
		Status:     stringPtr("active"),
	}

	pageReq := &PageRequest{
		Page:     1,
		PageSize: 5,
		OrderBy:  "salary",
		Order:    "DESC",
	}

	result, err := ce.userDAO.SearchUsers(condition, pageReq)
	if err != nil {
		log.Printf("搜索失败: %v", err)
	} else {
		fmt.Printf("搜索结果: 总共%d条, 第%d页, 共%d页\n", result.Total, result.Page, result.Pages)
		users := result.Data.([]*User)
		for _, user := range users {
			fmt.Printf("- %s (%d岁): %.2f\n", *user.RealName, *user.Age, *user.Salary)
		}
	}

	// 5.2 复杂动态条件
	fmt.Println("\n5.2 复杂动态搜索:")
	condition2 := &SearchCondition{
		RealName:  stringPtr("张"),
		MinSalary: float64Ptr(8000),
	}

	pageReq2 := DefaultPageRequest()
	pageReq2.OrderBy = "created_at"
	pageReq2.Order = "DESC"

	result2, err := ce.userDAO.SearchUsers(condition2, pageReq2)
	if err != nil {
		log.Printf("搜索失败: %v", err)
	} else {
		fmt.Printf("搜索结果: 总共%d条记录\n", result2.Total)
		users := result2.Data.([]*User)
		for _, user := range users {
			fmt.Printf("- %s: %.2f (%s)\n", *user.RealName, *user.Salary, user.CreatedAt.Format("2006-01-02"))
		}
	}

	fmt.Println()
}

// JoinQueryExamples 联合查询示例
func (ce *CompleteExample) JoinQueryExamples() {
	fmt.Println("6. 联合查询示例")
	fmt.Println("================")

	// 6.1 用户和部门联合查询
	fmt.Println("6.1 查询用户及其部门信息:")
	ex := example.NewExample()
	ex.CreateCriteria().
		AndEqualTo("u.status", "active").
		AndIsNotNull("u.department")
	ex.SetOrderByClause("u.department ASC, u.salary DESC")
	ex.SetLimit(0, 5)

	userDepts, err := ce.userDAO.SelectWithDepartment(ex)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 条用户部门记录:\n", len(userDepts))
		for _, ud := range userDepts {
			deptName := "未知部门"
			if ud.DepartmentName != nil {
				deptName = *ud.DepartmentName
			}
			fmt.Printf("- %s (%s): %.2f - %s\n", *ud.RealName, *ud.Department, *ud.Salary, deptName)
		}
	}

	fmt.Println()
}

// DataOperationExamples 数据操作示例
func (ce *CompleteExample) DataOperationExamples() {
	fmt.Println("7. 数据操作示例")
	fmt.Println("================")

	// 7.1 插入新用户
	fmt.Println("7.1 插入新用户:")
	newUser := &User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "hashed_password",
		RealName: stringPtr("测试用户"),
		Age:      intPtr(30),
		Gender:   stringPtr("male"),
		Phone:    stringPtr("13800000000"),
		Department: stringPtr("测试部"),
		Position:   stringPtr("测试工程师"),
		Salary:     float64Ptr(9000.00),
		Status:     "active",
	}

	err := ce.userDAO.Insert(newUser)
	if err != nil {
		log.Printf("插入用户失败: %v", err)
	} else {
		fmt.Printf("成功插入用户，ID: %d\n", newUser.ID)
	}

	// 7.2 更新用户
	fmt.Println("\n7.2 更新用户信息:")
	if newUser.ID > 0 {
		newUser.Salary = float64Ptr(9500.00)
		newUser.Position = stringPtr("高级测试工程师")

		ex := example.NewExample()
		ex.CreateCriteria().AndEqualTo("id", newUser.ID)

		affected, err := ce.userDAO.UpdateByExample(newUser, ex)
		if err != nil {
			log.Printf("更新用户失败: %v", err)
		} else {
			fmt.Printf("成功更新 %d 条记录\n", affected)
		}
	}

	// 7.3 删除测试用户
	fmt.Println("\n7.3 删除测试用户:")
	if newUser.ID > 0 {
		ex := example.NewExample()
		ex.CreateCriteria().AndEqualTo("id", newUser.ID)

		affected, err := ce.userDAO.DeleteByExample(ex)
		if err != nil {
			log.Printf("删除用户失败: %v", err)
		} else {
			fmt.Printf("成功删除 %d 条记录\n", affected)
		}
	}

	fmt.Println()
}

// 辅助函数
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// RunCompleteExample 运行完整示例的入口函数
func RunCompleteExample() {
	fmt.Println("注意：这是一个完整的 Example 使用示例")
	fmt.Println("要运行此示例，请：")
	fmt.Println("1. 安装 MySQL 数据库")
	fmt.Println("2. 创建数据库：CREATE DATABASE gobatis_example;")
	fmt.Println("3. 执行 schema.sql 中的建表语句")
	fmt.Println("4. 修改 DSN 连接字符串")
	fmt.Println("5. 安装 MySQL 驱动：go get github.com/go-sql-driver/mysql")
	fmt.Println()

	// 示例 DSN，实际使用时需要替换为真实的数据库连接信息
	exampleDSN := "user:password@tcp(localhost:3306)/gobatis_example?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Printf("示例 DSN: %s\n", exampleDSN)
	fmt.Println()

	// 如果要实际运行，取消注释以下代码：
	/*
	ce, err := NewCompleteExample(exampleDSN)
	if err != nil {
		log.Fatalf("创建示例失败: %v", err)
	}
	defer ce.Close()

	ce.RunAllExamples()
	*/

	fmt.Println("示例代码已准备就绪，请按照上述步骤配置数据库后运行。")
}