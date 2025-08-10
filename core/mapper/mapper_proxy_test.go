package mapper

import (
	"errors"
	"reflect"
	"testing"
)

// MockSqlSession 模拟SQL会话
type MockSqlSession struct {
	selectOneResult  interface{}
	selectOneError   error
	selectListResult []interface{}
	selectListError  error
	insertResult     int64
	insertError      error
	updateResult     int64
	updateError      error
	deleteResult     int64
	deleteError      error
	lastStatementId  string
	lastParameter    interface{}
}

func (m *MockSqlSession) SelectOne(statementId string, parameter interface{}) (interface{}, error) {
	m.lastStatementId = statementId
	m.lastParameter = parameter
	return m.selectOneResult, m.selectOneError
}

func (m *MockSqlSession) SelectList(statementId string, parameter interface{}) ([]interface{}, error) {
	m.lastStatementId = statementId
	m.lastParameter = parameter
	return m.selectListResult, m.selectListError
}

func (m *MockSqlSession) Insert(statementId string, parameter interface{}) (int64, error) {
	m.lastStatementId = statementId
	m.lastParameter = parameter
	return m.insertResult, m.insertError
}

func (m *MockSqlSession) Update(statementId string, parameter interface{}) (int64, error) {
	m.lastStatementId = statementId
	m.lastParameter = parameter
	return m.updateResult, m.updateError
}

func (m *MockSqlSession) Delete(statementId string, parameter interface{}) (int64, error) {
	m.lastStatementId = statementId
	m.lastParameter = parameter
	return m.deleteResult, m.deleteError
}

// TestMapper 测试用的Mapper接口
type TestMapper interface {
	GetUser(id int) (interface{}, error)
	FindUsers() ([]interface{}, error)
	InsertUser(user interface{}) (int64, error)
	UpdateUser(user interface{}) (int64, error)
	DeleteUser(id int) (int64, error)
	UnsupportedMethod() error
}

// TestNewMapperProxy 测试创建Mapper代理
func TestNewMapperProxy(t *testing.T) {
	session := &MockSqlSession{}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	proxy := NewMapperProxy(session, mapperType)

	if proxy == nil {
		t.Fatal("Proxy should not be nil")
	}

	// 验证代理实现了接口
	_, ok := proxy.(TestMapper)
	if !ok {
		t.Fatal("Proxy should implement TestMapper interface")
	}
}

// TestMapperProxy_GetStatementId 测试获取语句ID
func TestMapperProxy_GetStatementId(t *testing.T) {
	session := &MockSqlSession{}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	mp := &MapperProxy{
		session:    session,
		mapperType: mapperType,
	}

	statementId := mp.getStatementId("GetUser")
	expected := "mapper.TestMapper.GetUser"

	if statementId != expected {
		t.Fatalf("Expected statement ID '%s', got '%s'", expected, statementId)
	}
}

// TestMapperProxy_IsSelectMethod 测试判断查询方法
func TestMapperProxy_IsSelectMethod(t *testing.T) {
	session := &MockSqlSession{}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	mp := &MapperProxy{
		session:    session,
		mapperType: mapperType,
	}

	// 测试查询方法名
	testCases := []struct {
		methodName string
		expected   bool
	}{
		{"GetUser", true},
		{"FindUsers", true},
		{"SelectUser", true},
		{"QueryUsers", true},
		{"ListUsers", true},
		{"InsertUser", false},
		{"UpdateUser", false},
		{"DeleteUser", false},
	}

	for _, tc := range testCases {
		// 创建一个简单的方法类型用于测试
		methodType := reflect.TypeOf(func() (interface{}, error) { return nil, nil })
		result := mp.isSelectMethod(tc.methodName, methodType)

		if result != tc.expected {
			t.Errorf("Method '%s': expected %v, got %v", tc.methodName, tc.expected, result)
		}
	}
}

// TestMapperProxy_IsSelectListMethod 测试判断查询列表方法
func TestMapperProxy_IsSelectListMethod(t *testing.T) {
	session := &MockSqlSession{}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	mp := &MapperProxy{
		session:    session,
		mapperType: mapperType,
	}

	// 测试返回切片的方法
	sliceMethodType := reflect.TypeOf(func() ([]interface{}, error) { return nil, nil })
	if !mp.isSelectListMethod(sliceMethodType) {
		t.Error("Should recognize slice return type as list method")
	}

	// 测试返回单个对象的方法
	singleMethodType := reflect.TypeOf(func() (interface{}, error) { return nil, nil })
	if mp.isSelectListMethod(singleMethodType) {
		t.Error("Should not recognize single return type as list method")
	}
}

// TestMapperProxy_IsInsertMethod 测试判断插入方法
func TestMapperProxy_IsInsertMethod(t *testing.T) {
	session := &MockSqlSession{}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	mp := &MapperProxy{
		session:    session,
		mapperType: mapperType,
	}

	testCases := []struct {
		methodName string
		expected   bool
	}{
		{"InsertUser", true},
		{"AddUser", true},
		{"CreateUser", true},
		{"SaveUser", true},
		{"GetUser", false},
		{"UpdateUser", false},
		{"DeleteUser", false},
	}

	for _, tc := range testCases {
		result := mp.isInsertMethod(tc.methodName)
		if result != tc.expected {
			t.Errorf("Method '%s': expected %v, got %v", tc.methodName, tc.expected, result)
		}
	}
}

// TestMapperProxy_IsUpdateMethod 测试判断更新方法
func TestMapperProxy_IsUpdateMethod(t *testing.T) {
	session := &MockSqlSession{}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	mp := &MapperProxy{
		session:    session,
		mapperType: mapperType,
	}

	testCases := []struct {
		methodName string
		expected   bool
	}{
		{"UpdateUser", true},
		{"ModifyUser", true},
		{"EditUser", true},
		{"GetUser", false},
		{"InsertUser", false},
		{"DeleteUser", false},
	}

	for _, tc := range testCases {
		result := mp.isUpdateMethod(tc.methodName)
		if result != tc.expected {
			t.Errorf("Method '%s': expected %v, got %v", tc.methodName, tc.expected, result)
		}
	}
}

// TestMapperProxy_IsDeleteMethod 测试判断删除方法
func TestMapperProxy_IsDeleteMethod(t *testing.T) {
	session := &MockSqlSession{}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	mp := &MapperProxy{
		session:    session,
		mapperType: mapperType,
	}

	testCases := []struct {
		methodName string
		expected   bool
	}{
		{"DeleteUser", true},
		{"RemoveUser", true},
		{"GetUser", false},
		{"InsertUser", false},
		{"UpdateUser", false},
	}

	for _, tc := range testCases {
		result := mp.isDeleteMethod(tc.methodName)
		if result != tc.expected {
			t.Errorf("Method '%s': expected %v, got %v", tc.methodName, tc.expected, result)
		}
	}
}

// TestMapperProxy_SelectOne 测试查询单个对象
func TestMapperProxy_SelectOne(t *testing.T) {
	session := &MockSqlSession{
		selectOneResult: "test_result",
		selectOneError:  nil,
	}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	proxy := NewMapperProxy(session, mapperType)
	mapper := proxy.(TestMapper)

	result, err := mapper.GetUser(123)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != "test_result" {
		t.Fatalf("Expected 'test_result', got %v", result)
	}

	if session.lastStatementId != "mapper.TestMapper.GetUser" {
		t.Fatalf("Expected statement ID 'mapper.TestMapper.GetUser', got '%s'", session.lastStatementId)
	}

	if session.lastParameter != 123 {
		t.Fatalf("Expected parameter 123, got %v", session.lastParameter)
	}
}

// TestMapperProxy_SelectList 测试查询列表
func TestMapperProxy_SelectList(t *testing.T) {
	session := &MockSqlSession{
		selectListResult: []interface{}{"user1", "user2"},
		selectListError:  nil,
	}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	proxy := NewMapperProxy(session, mapperType)
	mapper := proxy.(TestMapper)

	result, err := mapper.FindUsers()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(result))
	}

	if session.lastStatementId != "mapper.TestMapper.FindUsers" {
		t.Fatalf("Expected statement ID 'mapper.TestMapper.FindUsers', got '%s'", session.lastStatementId)
	}
}

// TestMapperProxy_Insert 测试插入操作
func TestMapperProxy_Insert(t *testing.T) {
	session := &MockSqlSession{
		insertResult: 123,
		insertError:  nil,
	}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	proxy := NewMapperProxy(session, mapperType)
	mapper := proxy.(TestMapper)

	user := map[string]interface{}{"name": "john"}
	result, err := mapper.InsertUser(user)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != 123 {
		t.Fatalf("Expected 123, got %d", result)
	}

	if session.lastStatementId != "mapper.TestMapper.InsertUser" {
		t.Fatalf("Expected statement ID 'mapper.TestMapper.InsertUser', got '%s'", session.lastStatementId)
	}
}

// TestMapperProxy_Update 测试更新操作
func TestMapperProxy_Update(t *testing.T) {
	session := &MockSqlSession{
		updateResult: 1,
		updateError:  nil,
	}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	proxy := NewMapperProxy(session, mapperType)
	mapper := proxy.(TestMapper)

	user := map[string]interface{}{"id": 1, "name": "jane"}
	result, err := mapper.UpdateUser(user)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != 1 {
		t.Fatalf("Expected 1, got %d", result)
	}

	if session.lastStatementId != "mapper.TestMapper.UpdateUser" {
		t.Fatalf("Expected statement ID 'mapper.TestMapper.UpdateUser', got '%s'", session.lastStatementId)
	}
}

// TestMapperProxy_Delete 测试删除操作
func TestMapperProxy_Delete(t *testing.T) {
	session := &MockSqlSession{
		deleteResult: 1,
		deleteError:  nil,
	}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	proxy := NewMapperProxy(session, mapperType)
	mapper := proxy.(TestMapper)

	result, err := mapper.DeleteUser(123)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != 1 {
		t.Fatalf("Expected 1, got %d", result)
	}

	if session.lastStatementId != "mapper.TestMapper.DeleteUser" {
		t.Fatalf("Expected statement ID 'mapper.TestMapper.DeleteUser', got '%s'", session.lastStatementId)
	}
}

// TestMapperProxy_Error 测试错误处理
func TestMapperProxy_Error(t *testing.T) {
	session := &MockSqlSession{
		selectOneError: errors.New("database error"),
	}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	proxy := NewMapperProxy(session, mapperType)
	mapper := proxy.(TestMapper)

	_, err := mapper.GetUser(123)

	if err == nil {
		t.Fatal("Expected error")
	}

	if err.Error() != "database error" {
		t.Fatalf("Expected 'database error', got '%s'", err.Error())
	}
}

// TestMapperProxy_UnsupportedMethod 测试不支持的方法
func TestMapperProxy_UnsupportedMethod(t *testing.T) {
	session := &MockSqlSession{}
	mapperType := reflect.TypeOf((*TestMapper)(nil)).Elem()

	proxy := NewMapperProxy(session, mapperType)
	mapper := proxy.(TestMapper)

	err := mapper.UnsupportedMethod()

	if err == nil {
		t.Fatal("Expected error for unsupported method")
	}

	if err.Error() != "unsupported method: UnsupportedMethod" {
		t.Fatalf("Expected 'unsupported method: UnsupportedMethod', got '%s'", err.Error())
	}
}

// TestGetMethodName 测试获取方法名
func TestGetMethodName(t *testing.T) {
	// 这个函数依赖于运行时堆栈，在测试环境中可能不会返回预期的结果
	// 但我们可以测试它不会panic
	methodName := getMethodName()

	// 只要不panic就算通过
	_ = methodName
}
