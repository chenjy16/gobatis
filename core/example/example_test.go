package example

import (
	"strings"
	"testing"
)

func TestNewExample(t *testing.T) {
	example := NewExample()

	if example == nil {
		t.Error("Expected non-nil example")
	}

	if example.oredCriteria == nil {
		t.Error("Expected non-nil oredCriteria")
	}

	if len(example.oredCriteria) != 0 {
		t.Error("Expected empty oredCriteria")
	}

	if example.distinct {
		t.Error("Expected distinct to be false")
	}
}

func TestExample_CreateCriteria(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()

	if criteria == nil {
		t.Error("Expected non-nil criteria")
	}

	if criteria.criteria == nil {
		t.Error("Expected non-nil criteria.criteria")
	}

	if len(criteria.criteria) != 0 {
		t.Error("Expected empty criteria.criteria")
	}

	if criteria.valid {
		t.Error("Expected criteria.valid to be false")
	}
}

func TestCriteria_AndEqualTo(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	result := criteria.AndEqualTo("name", "test")

	if result != criteria {
		t.Error("Expected method chaining")
	}

	if len(criteria.criteria) != 1 {
		t.Error("Expected one criterion")
	}

	criterion := criteria.criteria[0]
	if criterion.condition != "name =" {
		t.Errorf("Expected condition 'name =', got: %s", criterion.condition)
	}

	if criterion.value != "test" {
		t.Errorf("Expected value 'test', got: %v", criterion.value)
	}

	if !criterion.singleValue {
		t.Error("Expected singleValue to be true")
	}
}

func TestCriteria_AndNotEqualTo(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndNotEqualTo("status", "inactive")

	criterion := criteria.criteria[0]
	if criterion.condition != "status <>" {
		t.Errorf("Expected condition 'status <>', got: %s", criterion.condition)
	}
}

func TestCriteria_AndGreaterThan(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndGreaterThan("age", 18)

	criterion := criteria.criteria[0]
	if criterion.condition != "age >" {
		t.Errorf("Expected condition 'age >', got: %s", criterion.condition)
	}
}

func TestCriteria_AndLessThan(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndLessThan("price", 100)

	criterion := criteria.criteria[0]
	if criterion.condition != "price <" {
		t.Errorf("Expected condition 'price <', got: %s", criterion.condition)
	}
}

func TestCriteria_AndLike(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndLike("title", "%test%")

	criterion := criteria.criteria[0]
	if criterion.condition != "title LIKE" {
		t.Errorf("Expected condition 'title LIKE', got: %s", criterion.condition)
	}

	if criterion.value != "%test%" {
		t.Errorf("Expected value '%%test%%', got: %v", criterion.value)
	}
}

func TestCriteria_AndIn(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	values := []interface{}{1, 2, 3}
	criteria.AndIn("id", values)

	criterion := criteria.criteria[0]
	if criterion.condition != "id IN" {
		t.Errorf("Expected condition 'id IN', got: %s", criterion.condition)
	}

	if !criterion.listValue {
		t.Error("Expected listValue to be true")
	}
}

func TestCriteria_AndIsNull(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndIsNull("deleted_at")

	criterion := criteria.criteria[0]
	if criterion.condition != "deleted_at IS NULL" {
		t.Errorf("Expected condition 'deleted_at IS NULL', got: %s", criterion.condition)
	}

	if !criterion.noValue {
		t.Error("Expected noValue to be true")
	}
}

func TestCriteria_AndIsNotNull(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndIsNotNull("created_at")

	criterion := criteria.criteria[0]
	if criterion.condition != "created_at IS NOT NULL" {
		t.Errorf("Expected condition 'created_at IS NOT NULL', got: %s", criterion.condition)
	}

	if !criterion.noValue {
		t.Error("Expected noValue to be true")
	}
}

func TestCriteria_AndBetween(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndBetween("age", 18, 65)

	criterion := criteria.criteria[0]
	if criterion.condition != "age BETWEEN" {
		t.Errorf("Expected condition 'age BETWEEN', got: %s", criterion.condition)
	}

	if !criterion.betweenValue {
		t.Error("Expected betweenValue to be true")
	}

	if criterion.value != 18 {
		t.Errorf("Expected value 18, got: %v", criterion.value)
	}

	if criterion.secondValue != 65 {
		t.Errorf("Expected secondValue 65, got: %v", criterion.secondValue)
	}
}

func TestExample_BuildSQL_Simple(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndEqualTo("name", "test")

	sql, args := example.BuildSQL("SELECT * FROM users")

	expectedSQL := "SELECT * FROM users WHERE (name = ?)"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 1 || args[0] != "test" {
		t.Errorf("Expected args: [test], got: %v", args)
	}
}

func TestExample_BuildSQL_Complex(t *testing.T) {
	example := NewExample()

	// 第一个条件组
	criteria1 := example.CreateCriteria()
	criteria1.AndEqualTo("status", "active").
		AndGreaterThan("age", 18).
		AndLike("name", "%john%")

	// 第二个条件组 (OR)
	criteria2 := example.CreateCriteria()
	criteria2.AndEqualTo("type", "vip").
		AndIn("department", []interface{}{"IT", "HR"})
	example.Or(*criteria2)

	// 设置排序
	example.SetOrderByClause("created_at DESC, name ASC")

	sql, args := example.BuildSQL("SELECT * FROM users")

	// 检查 SQL 包含预期的部分
	if !strings.Contains(sql, "WHERE") {
		t.Error("Expected WHERE clause")
	}

	if !strings.Contains(sql, "ORDER BY created_at DESC, name ASC") {
		t.Error("Expected ORDER BY clause")
	}

	if !strings.Contains(sql, "OR") {
		t.Error("Expected OR clause")
	}

	// 检查参数数量 (status, age, name, type, IT, HR)
	expectedArgsCount := 6
	if len(args) != expectedArgsCount {
		t.Errorf("Expected %d args, got: %d", expectedArgsCount, len(args))
	}
}

func TestExample_BuildSQL_WithDistinct(t *testing.T) {
	example := NewExample()
	example.SetDistinct(true)

	criteria := example.CreateCriteria()
	criteria.AndEqualTo("status", "active")

	sql, args := example.BuildSQL("SELECT id, name FROM users")

	if !strings.Contains(sql, "SELECT DISTINCT") {
		t.Error("Expected DISTINCT in SQL")
	}

	if len(args) != 1 || args[0] != "active" {
		t.Errorf("Expected args: [active], got: %v", args)
	}
}

func TestExample_BuildSQL_WithLimit(t *testing.T) {
	example := NewExample()
	example.SetLimit(0, 10)

	criteria := example.CreateCriteria()
	criteria.AndEqualTo("status", "active")

	sql, args := example.BuildSQL("SELECT * FROM users")

	if !strings.Contains(sql, "LIMIT 0, 10") {
		t.Error("Expected LIMIT clause")
	}

	if len(args) != 1 || args[0] != "active" {
		t.Errorf("Expected args: [active], got: %v", args)
	}
}

func TestExample_BuildSQL_NoConditions(t *testing.T) {
	example := NewExample()
	sql, args := example.BuildSQL("SELECT * FROM users")

	expectedSQL := "SELECT * FROM users"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 0 {
		t.Errorf("Expected no args, got: %v", args)
	}
}

func TestCriteria_BuildClause_Multiple(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndEqualTo("name", "test").
		AndGreaterThan("age", 18).
		AndIsNotNull("email")

	clause, args := criteria.buildClause()

	expectedClause := "name = ? AND age > ? AND email IS NOT NULL"
	if clause != expectedClause {
		t.Errorf("Expected clause: %s, got: %s", expectedClause, clause)
	}

	if len(args) != 2 {
		t.Errorf("Expected 2 args, got: %d", len(args))
	}

	if args[0] != "test" || args[1] != 18 {
		t.Errorf("Expected args: [test, 18], got: %v", args)
	}
}

func TestCriteria_ChainedMethods(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	result := criteria.
		AndEqualTo("status", "active").
		AndGreaterThan("age", 18).
		AndLike("name", "%john%").
		AndIsNotNull("email")

	if result != criteria {
		t.Error("Expected method chaining to return same instance")
	}

	if len(criteria.criteria) != 4 {
		t.Error("Expected four criteria")
	}

	if !criteria.valid {
		t.Error("Expected criteria to be valid")
	}
}

func TestExample_Clear(t *testing.T) {
	example := NewExample()
	criteria := example.CreateCriteria()
	criteria.AndEqualTo("name", "test")
	example.Or(*criteria)
	example.SetOrderByClause("name ASC")
	example.SetDistinct(true)
	example.SetLimit(0, 10)

	example.Clear()

	if len(example.oredCriteria) != 0 {
		t.Error("Expected empty oredCriteria after clear")
	}

	if example.orderByClause != "" {
		t.Error("Expected empty orderByClause after clear")
	}

	if example.distinct {
		t.Error("Expected distinct to be false after clear")
	}

	if example.limitStart != nil || example.limitEnd != nil {
		t.Error("Expected limit to be nil after clear")
	}
}

// TestOrderByClauseSecurity 测试ORDER BY子句的安全验证
func TestOrderByClauseSecurity(t *testing.T) {
	example := NewExample()

	// 测试有效的ORDER BY子句
	validOrderByClauses := []string{
		"name ASC",
		"created_at DESC",
		"name ASC, created_at DESC",
		"users.name ASC",
		"u.name ASC, u.created_at DESC",
		"id",
		"",
	}

	for _, orderBy := range validOrderByClauses {
		example.SetOrderByClause(orderBy)
		if example.orderByClause != orderBy {
			t.Errorf("Expected valid ORDER BY clause '%s' to be set", orderBy)
		}
	}

	// 测试无效的ORDER BY子句（可能包含SQL注入）
	invalidOrderByClauses := []string{
		"name; DROP TABLE users;",
		"name UNION SELECT * FROM passwords",
		"name' OR '1'='1",
		"name/*comment*/",
		"name--comment",
		"name()",
		"1=1",
		"name + 1",
	}

	for _, orderBy := range invalidOrderByClauses {
		originalOrderBy := example.orderByClause
		example.SetOrderByClause(orderBy)
		if example.orderByClause == orderBy {
			t.Errorf("Expected invalid ORDER BY clause '%s' to be rejected", orderBy)
		}
		// 恢复原始值用于下次测试
		example.orderByClause = originalOrderBy
	}
}
