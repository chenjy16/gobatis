package wrapper

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewQueryWrapper(t *testing.T) {
	wrapper := NewQueryWrapper()

	if wrapper == nil {
		t.Error("Expected non-nil wrapper")
	}

	if wrapper.conditions == nil {
		t.Error("Expected non-nil conditions")
	}

	if wrapper.orderBy == nil {
		t.Error("Expected non-nil orderBy")
	}

	if wrapper.groupBy == nil {
		t.Error("Expected non-nil groupBy")
	}

	if wrapper.having == nil {
		t.Error("Expected non-nil having")
	}

	if len(wrapper.conditions) != 0 {
		t.Error("Expected empty conditions")
	}

	if len(wrapper.orderBy) != 0 {
		t.Error("Expected empty orderBy")
	}

	if len(wrapper.groupBy) != 0 {
		t.Error("Expected empty groupBy")
	}

	if len(wrapper.having) != 0 {
		t.Error("Expected empty having")
	}
}

func TestQueryWrapper_Eq(t *testing.T) {
	wrapper := NewQueryWrapper()
	result := wrapper.Eq("name", "test")

	if result != wrapper {
		t.Error("Expected method chaining")
	}

	if len(wrapper.conditions) != 1 {
		t.Error("Expected one condition")
	}

	condition := wrapper.conditions[0]
	if condition.Column != "name" {
		t.Errorf("Expected column 'name', got: %s", condition.Column)
	}

	if condition.Operator != "=" {
		t.Errorf("Expected operator '=', got: %s", condition.Operator)
	}

	if condition.Value != "test" {
		t.Errorf("Expected value 'test', got: %v", condition.Value)
	}

	if condition.Logic != "AND" {
		t.Errorf("Expected logic 'AND', got: %s", condition.Logic)
	}
}

func TestQueryWrapper_Ne(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.Ne("status", "inactive")

	condition := wrapper.conditions[0]
	if condition.Operator != "!=" {
		t.Errorf("Expected operator '!=', got: %s", condition.Operator)
	}
}

func TestQueryWrapper_Gt(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.Gt("age", 18)

	condition := wrapper.conditions[0]
	if condition.Operator != ">" {
		t.Errorf("Expected operator '>', got: %s", condition.Operator)
	}
}

func TestQueryWrapper_Ge(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.Ge("score", 60)

	condition := wrapper.conditions[0]
	if condition.Operator != ">=" {
		t.Errorf("Expected operator '>=', got: %s", condition.Operator)
	}
}

func TestQueryWrapper_Lt(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.Lt("price", 100)

	condition := wrapper.conditions[0]
	if condition.Operator != "<" {
		t.Errorf("Expected operator '<', got: %s", condition.Operator)
	}
}

func TestQueryWrapper_Le(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.Le("count", 10)

	condition := wrapper.conditions[0]
	if condition.Operator != "<=" {
		t.Errorf("Expected operator '<=', got: %s", condition.Operator)
	}
}

func TestQueryWrapper_Like(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.Like("title", "test")

	condition := wrapper.conditions[0]
	if condition.Operator != "LIKE" {
		t.Errorf("Expected operator 'LIKE', got: %s", condition.Operator)
	}

	if condition.Value != "%test%" {
		t.Errorf("Expected value '%%test%%', got: %v", condition.Value)
	}
}

func TestQueryWrapper_NotLike(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.NotLike("description", "spam")

	condition := wrapper.conditions[0]
	if condition.Operator != "NOT LIKE" {
		t.Errorf("Expected operator 'NOT LIKE', got: %s", condition.Operator)
	}

	if condition.Value != "%spam%" {
		t.Errorf("Expected value '%%spam%%', got: %v", condition.Value)
	}
}

func TestQueryWrapper_In(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.In("id", 1, 2, 3)

	condition := wrapper.conditions[0]
	if condition.Operator != "IN" {
		t.Errorf("Expected operator 'IN', got: %s", condition.Operator)
	}

	values := condition.Value.([]interface{})
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got: %d", len(values))
	}

	if values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Errorf("Expected values [1, 2, 3], got: %v", values)
	}
}

func TestQueryWrapper_NotIn(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.NotIn("status", "deleted", "archived")

	condition := wrapper.conditions[0]
	if condition.Operator != "NOT IN" {
		t.Errorf("Expected operator 'NOT IN', got: %s", condition.Operator)
	}
}

func TestQueryWrapper_IsNull(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.IsNull("deleted_at")

	condition := wrapper.conditions[0]
	if condition.Operator != "IS NULL" {
		t.Errorf("Expected operator 'IS NULL', got: %s", condition.Operator)
	}

	if condition.Value != nil {
		t.Errorf("Expected nil value, got: %v", condition.Value)
	}
}

func TestQueryWrapper_IsNotNull(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.IsNotNull("created_at")

	condition := wrapper.conditions[0]
	if condition.Operator != "IS NOT NULL" {
		t.Errorf("Expected operator 'IS NOT NULL', got: %s", condition.Operator)
	}

	if condition.Value != nil {
		t.Errorf("Expected nil value, got: %v", condition.Value)
	}
}

func TestQueryWrapper_Or(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.Eq("name", "test").Or().Eq("email", "test@example.com")

	if len(wrapper.conditions) != 2 {
		t.Error("Expected two conditions")
	}

	// First condition should be OR (after calling Or())
	if wrapper.conditions[0].Logic != "OR" {
		t.Errorf("Expected first condition logic 'OR', got: %s", wrapper.conditions[0].Logic)
	}

	// Second condition should be AND (default)
	if wrapper.conditions[1].Logic != "AND" {
		t.Errorf("Expected second condition logic 'AND', got: %s", wrapper.conditions[1].Logic)
	}

	// Test Or() with no conditions
	emptyWrapper := NewQueryWrapper()
	emptyWrapper.Or()
	if len(emptyWrapper.conditions) != 0 {
		t.Error("Expected no conditions for empty wrapper")
	}
}

func TestQueryWrapper_OrderByAsc(t *testing.T) {
	wrapper := NewQueryWrapper()
	result := wrapper.OrderByAsc("name", "created_at")

	if result != wrapper {
		t.Error("Expected method chaining")
	}

	if len(wrapper.orderBy) != 2 {
		t.Error("Expected two order by clauses")
	}

	if wrapper.orderBy[0].Column != "name" || wrapper.orderBy[0].Desc != false {
		t.Error("Expected ascending order for name")
	}

	if wrapper.orderBy[1].Column != "created_at" || wrapper.orderBy[1].Desc != false {
		t.Error("Expected ascending order for created_at")
	}
}

func TestQueryWrapper_OrderByDesc(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.OrderByDesc("updated_at", "id")

	if len(wrapper.orderBy) != 2 {
		t.Error("Expected two order by clauses")
	}

	if wrapper.orderBy[0].Column != "updated_at" || wrapper.orderBy[0].Desc != true {
		t.Error("Expected descending order for updated_at")
	}

	if wrapper.orderBy[1].Column != "id" || wrapper.orderBy[1].Desc != true {
		t.Error("Expected descending order for id")
	}
}

func TestQueryWrapper_GroupBy(t *testing.T) {
	wrapper := NewQueryWrapper()
	result := wrapper.GroupBy("department", "status")

	if result != wrapper {
		t.Error("Expected method chaining")
	}

	if len(wrapper.groupBy) != 2 {
		t.Error("Expected two group by columns")
	}

	if wrapper.groupBy[0] != "department" || wrapper.groupBy[1] != "status" {
		t.Errorf("Expected ['department', 'status'], got: %v", wrapper.groupBy)
	}
}

func TestQueryWrapper_Having(t *testing.T) {
	wrapper := NewQueryWrapper()
	result := wrapper.Having("COUNT(*)", ">", 5)

	if result != wrapper {
		t.Error("Expected method chaining")
	}

	if len(wrapper.having) != 1 {
		t.Error("Expected one having condition")
	}

	condition := wrapper.having[0]
	if condition.Column != "COUNT(*)" {
		t.Errorf("Expected column 'COUNT(*)', got: %s", condition.Column)
	}

	if condition.Operator != ">" {
		t.Errorf("Expected operator '>', got: %s", condition.Operator)
	}

	if condition.Value != 5 {
		t.Errorf("Expected value 5, got: %v", condition.Value)
	}
}

func TestQueryWrapper_Limit(t *testing.T) {
	wrapper := NewQueryWrapper()
	result := wrapper.Limit(10)

	if result != wrapper {
		t.Error("Expected method chaining")
	}

	if wrapper.limit == nil {
		t.Error("Expected non-nil limit")
	}

	if *wrapper.limit != 10 {
		t.Errorf("Expected limit 10, got: %d", *wrapper.limit)
	}
}

func TestQueryWrapper_Offset(t *testing.T) {
	wrapper := NewQueryWrapper()
	result := wrapper.Offset(20)

	if result != wrapper {
		t.Error("Expected method chaining")
	}

	if wrapper.offset == nil {
		t.Error("Expected non-nil offset")
	}

	if *wrapper.offset != 20 {
		t.Errorf("Expected offset 20, got: %d", *wrapper.offset)
	}
}

func TestQueryWrapper_BuildSQL_Simple(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.Eq("name", "test")

	sql, args := wrapper.BuildSQL("SELECT * FROM users")

	expectedSQL := "SELECT * FROM users WHERE name = ?"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 1 || args[0] != "test" {
		t.Errorf("Expected args: [test], got: %v", args)
	}
}

func TestQueryWrapper_BuildSQL_Complex(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.Eq("status", "active").
		Gt("age", 18).
		Like("name", "john").
		In("department", "IT", "HR").
		OrderByDesc("created_at").
		OrderByAsc("name").
		GroupBy("department").
		Having("COUNT(*)", ">", 5).
		Limit(10).
		Offset(20)

	sql, args := wrapper.BuildSQL("SELECT * FROM users")

	// Check if SQL contains expected parts
	if !strings.Contains(sql, "WHERE") {
		t.Error("Expected WHERE clause")
	}

	if !strings.Contains(sql, "GROUP BY department") {
		t.Error("Expected GROUP BY clause")
	}

	if !strings.Contains(sql, "HAVING") {
		t.Error("Expected HAVING clause")
	}

	if !strings.Contains(sql, "ORDER BY") {
		t.Error("Expected ORDER BY clause")
	}

	if !strings.Contains(sql, "LIMIT 10") {
		t.Error("Expected LIMIT clause")
	}

	if !strings.Contains(sql, "OFFSET 20") {
		t.Error("Expected OFFSET clause")
	}

	// Check args count (status, age, name, IT, HR, having count)
	expectedArgsCount := 6
	if len(args) != expectedArgsCount {
		t.Errorf("Expected %d args, got: %d", expectedArgsCount, len(args))
	}
}

func TestQueryWrapper_BuildSQL_NoConditions(t *testing.T) {
	wrapper := NewQueryWrapper()
	sql, args := wrapper.BuildSQL("SELECT * FROM users")

	expectedSQL := "SELECT * FROM users"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(args) != 0 {
		t.Errorf("Expected no args, got: %v", args)
	}
}

func TestQueryWrapper_BuildCondition_IN(t *testing.T) {
	wrapper := NewQueryWrapper()
	condition := Condition{
		Column:   "id",
		Operator: "IN",
		Value:    []interface{}{1, 2, 3},
	}

	clause, args := wrapper.buildCondition(condition)

	expectedClause := "id IN (?, ?, ?)"
	if clause != expectedClause {
		t.Errorf("Expected clause: %s, got: %s", expectedClause, clause)
	}

	if !reflect.DeepEqual(args, []interface{}{1, 2, 3}) {
		t.Errorf("Expected args: [1, 2, 3], got: %v", args)
	}
}

func TestQueryWrapper_BuildCondition_IsNull(t *testing.T) {
	wrapper := NewQueryWrapper()
	condition := Condition{
		Column:   "deleted_at",
		Operator: "IS NULL",
		Value:    nil,
	}

	clause, args := wrapper.buildCondition(condition)

	expectedClause := "deleted_at IS NULL"
	if clause != expectedClause {
		t.Errorf("Expected clause: %s, got: %s", expectedClause, clause)
	}

	if len(args) != 0 {
		t.Errorf("Expected no args, got: %v", args)
	}
}

func TestQueryWrapper_BuildOrderByClause(t *testing.T) {
	wrapper := NewQueryWrapper()
	wrapper.OrderByAsc("name").OrderByDesc("created_at")

	clause := wrapper.buildOrderByClause()

	expectedClause := "name ASC, created_at DESC"
	if clause != expectedClause {
		t.Errorf("Expected clause: %s, got: %s", expectedClause, clause)
	}
}

func TestQueryWrapper_ChainedMethods(t *testing.T) {
	wrapper := NewQueryWrapper()
	result := wrapper.
		Eq("status", "active").
		Gt("age", 18).
		OrderByDesc("created_at").
		Limit(10)

	if result != wrapper {
		t.Error("Expected method chaining to return same instance")
	}

	if len(wrapper.conditions) != 2 {
		t.Error("Expected two conditions")
	}

	if len(wrapper.orderBy) != 1 {
		t.Error("Expected one order by clause")
	}

	if wrapper.limit == nil || *wrapper.limit != 10 {
		t.Error("Expected limit to be set to 10")
	}
}