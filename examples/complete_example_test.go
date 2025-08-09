package examples

import (
	"testing"
	"time"

	"gobatis/core/example"
)

// TestModels 测试模型结构
func TestModels(t *testing.T) {
	// 测试用户模型
	user := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		RealName: stringPtr("测试用户"),
		Age:      intPtr(25),
		Gender:   stringPtr("male"),
		Status:   "active",
	}

	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}

	if *user.RealName != "测试用户" {
		t.Errorf("Expected real_name '测试用户', got %s", *user.RealName)
	}
}

// TestSearchCondition 测试搜索条件
func TestSearchCondition(t *testing.T) {
	condition := &SearchCondition{
		Username:   stringPtr("test"),
		Department: stringPtr("IT"),
		MinAge:     intPtr(25),
		MaxAge:     intPtr(35),
		Status:     stringPtr("active"),
	}

	if *condition.Username != "test" {
		t.Errorf("Expected username 'test', got %s", *condition.Username)
	}

	if *condition.Department != "IT" {
		t.Errorf("Expected department 'IT', got %s", *condition.Department)
	}

	if *condition.MinAge != 25 {
		t.Errorf("Expected min_age 25, got %d", *condition.MinAge)
	}
}

// TestPageRequest 测试分页请求
func TestPageRequest(t *testing.T) {
	pageReq := DefaultPageRequest()

	if pageReq.Page != 1 {
		t.Errorf("Expected page 1, got %d", pageReq.Page)
	}

	if pageReq.PageSize != 10 {
		t.Errorf("Expected page_size 10, got %d", pageReq.PageSize)
	}

	// 测试偏移量计算
	offset := pageReq.GetOffset()
	if offset != 0 {
		t.Errorf("Expected offset 0, got %d", offset)
	}

	// 测试第二页
	pageReq.Page = 2
	offset = pageReq.GetOffset()
	if offset != 10 {
		t.Errorf("Expected offset 10, got %d", offset)
	}

	// 测试总页数计算
	pages := pageReq.CalculatePages(25)
	if pages != 3 {
		t.Errorf("Expected pages 3, got %d", pages)
	}
}

// TestExampleUsage 测试 Example 的使用
func TestExampleUsage(t *testing.T) {
	// 测试基本查询
	ex := example.NewExample()
	criteria := ex.CreateCriteria()
	criteria.AndEqualTo("status", "active")
	criteria.AndGreaterThan("age", 18)

	sql, args := ex.BuildSQL("SELECT * FROM users")
	
	expectedSQL := "SELECT * FROM users WHERE (status = ? AND age > ?)"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}

	if len(args) >= 1 && args[0] != "active" {
		t.Errorf("Expected first arg 'active', got %v", args[0])
	}

	if len(args) >= 2 && args[1] != 18 {
		t.Errorf("Expected second arg 18, got %v", args[1])
	}
}

// TestExampleWithOrder 测试带排序的 Example
func TestExampleWithOrder(t *testing.T) {
	ex := example.NewExample()
	ex.CreateCriteria().AndEqualTo("status", "active")
	ex.SetOrderByClause("created_at DESC")

	sql, args := ex.BuildSQL("SELECT * FROM users")
	
	expectedSQL := "SELECT * FROM users WHERE (status = ?) ORDER BY created_at DESC"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

// TestExampleWithLimit 测试带分页的 Example
func TestExampleWithLimit(t *testing.T) {
	ex := example.NewExample()
	ex.CreateCriteria().AndEqualTo("status", "active")
	ex.SetLimit(0, 10)

	sql, args := ex.BuildSQL("SELECT * FROM users")
	
	expectedSQL := "SELECT * FROM users WHERE (status = ?) LIMIT 0, 10"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

// TestExampleWithDistinct 测试 DISTINCT 查询
func TestExampleWithDistinct(t *testing.T) {
	ex := example.NewExample()
	ex.SetDistinct(true)
	ex.CreateCriteria().AndIsNotNull("department")

	sql, args := ex.BuildSQL("SELECT department FROM users")
	
	expectedSQL := "SELECT DISTINCT department FROM users WHERE (department IS NOT NULL)"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 0 {
		t.Errorf("Expected 0 args, got %d", len(args))
	}
}

// TestExampleWithBetween 测试 BETWEEN 查询
func TestExampleWithBetween(t *testing.T) {
	ex := example.NewExample()
	ex.CreateCriteria().AndBetween("age", 25, 35)

	sql, args := ex.BuildSQL("SELECT * FROM users")
	
	expectedSQL := "SELECT * FROM users WHERE (age BETWEEN ? AND ?)"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(args))
	}

	if args[0] != 25 {
		t.Errorf("Expected first arg 25, got %v", args[0])
	}

	if args[1] != 35 {
		t.Errorf("Expected second arg 35, got %v", args[1])
	}
}

// TestExampleWithIn 测试 IN 查询
func TestExampleWithIn(t *testing.T) {
	ex := example.NewExample()
	departments := []interface{}{"IT", "HR", "研发部"}
	ex.CreateCriteria().AndIn("department", departments)

	sql, args := ex.BuildSQL("SELECT * FROM users")
	
	expectedSQL := "SELECT * FROM users WHERE (department IN (?, ?, ?))"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}

	if args[0] != "IT" {
		t.Errorf("Expected first arg 'IT', got %v", args[0])
	}
}

// TestExampleWithLike 测试 LIKE 查询
func TestExampleWithLike(t *testing.T) {
	ex := example.NewExample()
	ex.CreateCriteria().AndLike("real_name", "%张%")

	sql, args := ex.BuildSQL("SELECT * FROM users")
	
	expectedSQL := "SELECT * FROM users WHERE (real_name LIKE ?)"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}

	if args[0] != "%张%" {
		t.Errorf("Expected arg '%%张%%', got %v", args[0])
	}
}

// TestExampleWithOr 测试 OR 查询
func TestExampleWithOr(t *testing.T) {
	ex := example.NewExample()
	
	// 第一组条件
	criteria1 := ex.CreateCriteria()
	criteria1.AndEqualTo("department", "IT")
	criteria1.AndGreaterThan("salary", 10000)
	
	// 第二组条件
	criteria2 := ex.CreateCriteria()
	criteria2.AndEqualTo("department", "HR")
	criteria2.AndGreaterThan("salary", 8000)
	ex.Or(*criteria2)

	sql, args := ex.BuildSQL("SELECT * FROM users")
	
	expectedSQL := "SELECT * FROM users WHERE (department = ? AND salary > ?) OR (department = ? AND salary > ?)"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 4 {
		t.Errorf("Expected 4 args, got %d", len(args))
	}
}

// TestExampleComplex 测试复杂查询
func TestExampleComplex(t *testing.T) {
	ex := example.NewExample()
	ex.SetDistinct(true)
	
	criteria := ex.CreateCriteria()
	criteria.AndEqualTo("status", "active")
	criteria.AndBetween("age", 25, 35)
	criteria.AndLike("real_name", "%张%")
	criteria.AndIsNotNull("phone")
	
	ex.SetOrderByClause("salary DESC, created_at ASC")
	ex.SetLimit(10, 20)

	sql, args := ex.BuildSQL("SELECT * FROM users")
	
	expectedSQL := "SELECT DISTINCT * FROM users WHERE (status = ? AND age BETWEEN ? AND ? AND real_name LIKE ? AND phone IS NOT NULL) ORDER BY salary DESC, created_at ASC LIMIT 10, 20"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 4 {
		t.Errorf("Expected 4 args, got %d", len(args))
	}
}

// TestRunCompleteExample 测试完整示例运行
func TestRunCompleteExample(t *testing.T) {
	// 这个测试只是确保函数可以正常调用，不会 panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RunCompleteExample panicked: %v", r)
		}
	}()

	RunCompleteExample()
}

// TestHelperFunctions 测试辅助函数
func TestHelperFunctions(t *testing.T) {
	// 测试 stringPtr
	str := "test"
	strPtr := stringPtr(str)
	if *strPtr != str {
		t.Errorf("Expected %s, got %s", str, *strPtr)
	}

	// 测试 intPtr
	num := 42
	numPtr := intPtr(num)
	if *numPtr != num {
		t.Errorf("Expected %d, got %d", num, *numPtr)
	}

	// 测试 float64Ptr
	f := 3.14
	fPtr := float64Ptr(f)
	if *fPtr != f {
		t.Errorf("Expected %f, got %f", f, *fPtr)
	}

	// 测试 timePtr
	now := time.Now()
	timePtr := timePtr(now)
	if !timePtr.Equal(now) {
		t.Errorf("Expected %v, got %v", now, *timePtr)
	}
}

// BenchmarkExampleBuild 性能测试
func BenchmarkExampleBuild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ex := example.NewExample()
		criteria := ex.CreateCriteria()
		criteria.AndEqualTo("status", "active")
		criteria.AndGreaterThan("age", 18)
		criteria.AndLike("name", "%test%")
		ex.SetOrderByClause("created_at DESC")
		ex.SetLimit(0, 10)
		
		_, _ = ex.BuildSQL("SELECT * FROM users")
	}
}