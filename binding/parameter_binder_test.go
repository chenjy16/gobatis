package binding

import (
	"reflect"
	"testing"
	"time"
)

// TestUser 测试用户结构体
type TestUser struct {
	ID       int64     `db:"id"`
	Username string    `db:"username"`
	Email    string    `db:"email"`
	CreateAt time.Time `db:"create_at"`
}

// TestNewParameterBinder 测试创建参数绑定器
func TestNewParameterBinder(t *testing.T) {
	binder := NewParameterBinder()
	if binder == nil {
		t.Fatal("ParameterBinder should not be nil")
	}

	_, ok := binder.(*DefaultParameterBinder)
	if !ok {
		t.Fatal("Should return DefaultParameterBinder instance")
	}
}

// TestBindParameters_NilParameter 测试空参数绑定
func TestBindParameters_NilParameter(t *testing.T) {
	binder := NewParameterBinder()
	sql := "SELECT * FROM users WHERE id = #{id}"

	processedSQL, args, err := binder.BindParameters(sql, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if processedSQL != sql {
		t.Fatalf("Expected SQL to remain unchanged, got: %s", processedSQL)
	}

	if args != nil {
		t.Fatal("Args should be nil for nil parameter")
	}
}

// TestBindParameters_NoParameters 测试无参数的SQL
func TestBindParameters_NoParameters(t *testing.T) {
	binder := NewParameterBinder()
	sql := "SELECT * FROM users"

	processedSQL, args, err := binder.BindParameters(sql, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if processedSQL != sql {
		t.Fatalf("Expected SQL to remain unchanged, got: %s", processedSQL)
	}

	if args != nil {
		t.Fatal("Args should be nil for SQL without parameters")
	}
}

// TestBindParameters_MapParameter 测试Map参数绑定
func TestBindParameters_MapParameter(t *testing.T) {
	binder := NewParameterBinder()
	sql := "SELECT * FROM users WHERE id = #{id} AND username = #{username}"
	params := map[string]interface{}{
		"id":       int64(1),
		"username": "testuser",
	}

	processedSQL, args, err := binder.BindParameters(sql, params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "SELECT * FROM users WHERE id = ? AND username = ?"
	if processedSQL != expectedSQL {
		t.Fatalf("Expected SQL: %s, got: %s", expectedSQL, processedSQL)
	}

	if len(args) != 2 {
		t.Fatalf("Expected 2 args, got: %d", len(args))
	}

	if args[0] != int64(1) {
		t.Fatalf("Expected first arg to be 1, got: %v", args[0])
	}

	if args[1] != "testuser" {
		t.Fatalf("Expected second arg to be 'testuser', got: %v", args[1])
	}
}

// TestBindParameters_MapParameter_MissingKey 测试Map参数缺少键
func TestBindParameters_MapParameter_MissingKey(t *testing.T) {
	binder := NewParameterBinder()
	sql := "SELECT * FROM users WHERE id = #{id} AND username = #{username}"
	params := map[string]interface{}{
		"id": int64(1),
		// username 缺失
	}

	processedSQL, args, err := binder.BindParameters(sql, params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "SELECT * FROM users WHERE id = ? AND username = ?"
	if processedSQL != expectedSQL {
		t.Fatalf("Expected SQL: %s, got: %s", expectedSQL, processedSQL)
	}

	if len(args) != 2 {
		t.Fatalf("Expected 2 args, got: %d", len(args))
	}

	if args[0] != int64(1) {
		t.Fatalf("Expected first arg to be 1, got: %v", args[0])
	}

	if args[1] != nil {
		t.Fatalf("Expected second arg to be nil, got: %v", args[1])
	}
}

// TestBindParameters_StructParameter 测试结构体参数绑定
func TestBindParameters_StructParameter(t *testing.T) {
	binder := NewParameterBinder()
	sql := "INSERT INTO users (id, username, email) VALUES (#{id}, #{username}, #{email})"
	user := TestUser{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		CreateAt: time.Now(),
	}

	processedSQL, args, err := binder.BindParameters(sql, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "INSERT INTO users (id, username, email) VALUES (?, ?, ?)"
	if processedSQL != expectedSQL {
		t.Fatalf("Expected SQL: %s, got: %s", expectedSQL, processedSQL)
	}

	if len(args) != 3 {
		t.Fatalf("Expected 3 args, got: %d", len(args))
	}

	if args[0] != int64(1) {
		t.Fatalf("Expected first arg to be 1, got: %v", args[0])
	}

	if args[1] != "testuser" {
		t.Fatalf("Expected second arg to be 'testuser', got: %v", args[1])
	}

	if args[2] != "test@example.com" {
		t.Fatalf("Expected third arg to be 'test@example.com', got: %v", args[2])
	}
}

// TestBindParameters_StructPointer 测试结构体指针参数绑定
func TestBindParameters_StructPointer(t *testing.T) {
	binder := NewParameterBinder()
	sql := "UPDATE users SET username = #{username} WHERE id = #{id}"
	user := &TestUser{
		ID:       1,
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	processedSQL, args, err := binder.BindParameters(sql, user)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "UPDATE users SET username = ? WHERE id = ?"
	if processedSQL != expectedSQL {
		t.Fatalf("Expected SQL: %s, got: %s", expectedSQL, processedSQL)
	}

	if len(args) != 2 {
		t.Fatalf("Expected 2 args, got: %d", len(args))
	}

	if args[0] != "updateduser" {
		t.Fatalf("Expected first arg to be 'updateduser', got: %v", args[0])
	}

	if args[1] != int64(1) {
		t.Fatalf("Expected second arg to be 1, got: %v", args[1])
	}
}

// TestBindParameters_NilPointer 测试空指针参数
func TestBindParameters_NilPointer(t *testing.T) {
	binder := NewParameterBinder()
	sql := "SELECT * FROM users WHERE id = #{id}"
	var user *TestUser = nil

	_, _, err := binder.BindParameters(sql, user)
	if err == nil {
		t.Fatal("Expected error for nil pointer parameter")
	}

	expectedError := "parameter is nil pointer"
	if err.Error() != expectedError {
		t.Fatalf("Expected error: %s, got: %s", expectedError, err.Error())
	}
}

// TestBindParameters_BasicType 测试基础类型参数绑定
func TestBindParameters_BasicType(t *testing.T) {
	binder := NewParameterBinder()
	sql := "SELECT * FROM users WHERE id = #{id}"

	processedSQL, args, err := binder.BindParameters(sql, int64(123))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "SELECT * FROM users WHERE id = ?"
	if processedSQL != expectedSQL {
		t.Fatalf("Expected SQL: %s, got: %s", expectedSQL, processedSQL)
	}

	if len(args) != 1 {
		t.Fatalf("Expected 1 arg, got: %d", len(args))
	}

	if args[0] != int64(123) {
		t.Fatalf("Expected arg to be 123, got: %v", args[0])
	}
}

// TestBindParameters_UnsupportedType 测试不支持的参数类型
func TestBindParameters_UnsupportedType(t *testing.T) {
	binder := NewParameterBinder()
	sql := "SELECT * FROM users WHERE id = #{id}"
	unsupported := []int{1, 2, 3} // slice类型

	_, _, err := binder.BindParameters(sql, unsupported)
	if err == nil {
		t.Fatal("Expected error for unsupported parameter type")
	}

	expectedError := "unsupported parameter type: []int"
	if err.Error() != expectedError {
		t.Fatalf("Expected error: %s, got: %s", expectedError, err.Error())
	}
}

// TestIsBasicType 测试基础类型判断
func TestIsBasicType(t *testing.T) {
	testCases := []struct {
		kind     reflect.Kind
		expected bool
	}{
		{reflect.Bool, true},
		{reflect.Int, true},
		{reflect.Int64, true},
		{reflect.String, true},
		{reflect.Float64, true},
		{reflect.Slice, false},
		{reflect.Map, false},
		{reflect.Struct, false},
		{reflect.Ptr, false},
	}

	for _, tc := range testCases {
		result := isBasicType(tc.kind)
		if result != tc.expected {
			t.Errorf("isBasicType(%v) = %v, expected %v", tc.kind, result, tc.expected)
		}
	}
}

// TestConvertValue 测试值类型转换
func TestConvertValue(t *testing.T) {
	// 测试nil值
	result, err := ConvertValue(nil, reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("Expected nil result for nil input")
	}

	// 测试相同类型转换
	result, err = ConvertValue(int64(123), reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != int64(123) {
		t.Fatalf("Expected 123, got: %v", result)
	}

	// 测试字符串到整数转换
	result, err = ConvertValue("456", reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != int64(456) {
		t.Fatalf("Expected 456, got: %v", result)
	}

	// 测试字符串到浮点数转换
	result, err = ConvertValue("123.45", reflect.TypeOf(float64(0)))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != float64(123.45) {
		t.Fatalf("Expected 123.45, got: %v", result)
	}

	// 测试无法转换的情况
	result, err = ConvertValue("invalid", reflect.TypeOf(int64(0)))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != "invalid" {
		t.Fatalf("Expected original value for invalid conversion, got: %v", result)
	}
}

// TestIsNumericType 测试数字类型判断
func TestIsNumericType(t *testing.T) {
	testCases := []struct {
		kind     reflect.Kind
		expected bool
	}{
		{reflect.Int, true},
		{reflect.Int64, true},
		{reflect.Float64, true},
		{reflect.Uint, true},
		{reflect.String, false},
		{reflect.Bool, false},
		{reflect.Slice, false},
	}

	for _, tc := range testCases {
		result := isNumericType(tc.kind)
		if result != tc.expected {
			t.Errorf("isNumericType(%v) = %v, expected %v", tc.kind, result, tc.expected)
		}
	}
}
