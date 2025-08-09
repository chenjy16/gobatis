package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

// MockPlugin 模拟插件用于测试
type MockPlugin struct {
	properties map[string]string
}

func (p *MockPlugin) Intercept(invocation *Invocation) (interface{}, error) {
	return invocation.Target, nil
}

func (p *MockPlugin) SetProperties(properties map[string]string) {
	p.properties = properties
}

// TestNewConfiguration 测试创建新配置
func TestNewConfiguration(t *testing.T) {
	config := NewConfiguration()
	
	if config == nil {
		t.Fatal("Configuration should not be nil")
	}
	
	if config.MapperConfig == nil {
		t.Fatal("MapperConfig should not be nil")
	}
	
	if config.MapperConfig.Mappers == nil {
		t.Fatal("Mappers should not be nil")
	}
	
	if config.Plugins == nil {
		t.Fatal("Plugins should not be nil")
	}
	
	if len(config.Plugins) != 0 {
		t.Fatal("Plugins should be empty initially")
	}
	
	if len(config.MapperConfig.Mappers) != 0 {
		t.Fatal("Mappers should be empty initially")
	}
}

// TestSetDataSource_InvalidDriver 测试设置无效数据源
func TestSetDataSource_InvalidDriver(t *testing.T) {
	config := NewConfiguration()
	
	err := config.SetDataSource("invalid_driver", "invalid_dsn")
	if err == nil {
		t.Fatal("Expected error for invalid driver")
	}
}

// TestAddPlugin 测试添加插件
func TestAddPlugin(t *testing.T) {
	config := NewConfiguration()
	plugin := &MockPlugin{}
	
	config.AddPlugin(plugin)
	
	if len(config.Plugins) != 1 {
		t.Fatalf("Expected 1 plugin, got %d", len(config.Plugins))
	}
	
	if config.Plugins[0] != plugin {
		t.Fatal("Plugin should be the same instance")
	}
}

// TestAddMultiplePlugins 测试添加多个插件
func TestAddMultiplePlugins(t *testing.T) {
	config := NewConfiguration()
	plugin1 := &MockPlugin{}
	plugin2 := &MockPlugin{}
	
	config.AddPlugin(plugin1)
	config.AddPlugin(plugin2)
	
	if len(config.Plugins) != 2 {
		t.Fatalf("Expected 2 plugins, got %d", len(config.Plugins))
	}
	
	if config.Plugins[0] != plugin1 {
		t.Fatal("First plugin should be plugin1")
	}
	
	if config.Plugins[1] != plugin2 {
		t.Fatal("Second plugin should be plugin2")
	}
}

// TestGetMapperStatement_NotExists 测试获取不存在的Mapper语句
func TestGetMapperStatement_NotExists(t *testing.T) {
	config := NewConfiguration()
	
	stmt, exists := config.GetMapperStatement("NonExistent.Statement")
	if exists {
		t.Fatal("Statement should not exist")
	}
	
	if stmt != nil {
		t.Fatal("Statement should be nil")
	}
}

// TestGetMapperStatement_Exists 测试获取存在的Mapper语句
func TestGetMapperStatement_Exists(t *testing.T) {
	config := NewConfiguration()
	
	// 手动添加一个语句
	expectedStmt := &MapperStatement{
		ID:            "TestMapper.TestStatement",
		SQL:           "SELECT * FROM test",
		StatementType: SELECT,
	}
	config.MapperConfig.Mappers["TestMapper.TestStatement"] = expectedStmt
	
	stmt, exists := config.GetMapperStatement("TestMapper.TestStatement")
	if !exists {
		t.Fatal("Statement should exist")
	}
	
	if stmt != expectedStmt {
		t.Fatal("Statement should be the same instance")
	}
}

// TestAddMapperXML_FileNotExists 测试添加不存在的XML文件
func TestAddMapperXML_FileNotExists(t *testing.T) {
	config := NewConfiguration()
	
	err := config.AddMapperXML("nonexistent.xml")
	if err == nil {
		t.Fatal("Expected error for non-existent file")
	}
}

// TestAddMapperXML_InvalidXML 测试添加无效的XML文件
func TestAddMapperXML_InvalidXML(t *testing.T) {
	config := NewConfiguration()
	
	// 创建临时的无效XML文件
	tempFile, err := ioutil.TempFile("", "invalid_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// 写入无效XML内容
	_, err = tempFile.WriteString("invalid xml content")
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	err = config.AddMapperXML(tempFile.Name())
	if err == nil {
		t.Fatal("Expected error for invalid XML")
	}
}

// TestAddMapperXML_ValidXML 测试添加有效的XML文件
func TestAddMapperXML_ValidXML(t *testing.T) {
	config := NewConfiguration()
	
	// 创建临时的有效XML文件
	tempFile, err := ioutil.TempFile("", "valid_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	// 写入有效XML内容
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<mapper namespace="TestMapper">
    <select id="GetUser" resultType="User">
        SELECT id, username FROM users WHERE id = #{id}
    </select>
    <insert id="InsertUser">
        INSERT INTO users (username) VALUES (#{username})
    </insert>
    <update id="UpdateUser">
        UPDATE users SET username = #{username} WHERE id = #{id}
    </update>
    <delete id="DeleteUser">
        DELETE FROM users WHERE id = #{id}
    </delete>
</mapper>`
	
	_, err = tempFile.WriteString(xmlContent)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	err = config.AddMapperXML(tempFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// 验证SELECT语句
	stmt, exists := config.GetMapperStatement("TestMapper.GetUser")
	if !exists {
		t.Fatal("GetUser statement should exist")
	}
	if stmt.StatementType != SELECT {
		t.Fatalf("Expected SELECT statement type, got %v", stmt.StatementType)
	}
	if stmt.SQL != "SELECT id, username FROM users WHERE id = #{id}" {
		t.Fatalf("Unexpected SQL: %s", stmt.SQL)
	}
	
	// 验证INSERT语句
	stmt, exists = config.GetMapperStatement("TestMapper.InsertUser")
	if !exists {
		t.Fatal("InsertUser statement should exist")
	}
	if stmt.StatementType != INSERT {
		t.Fatalf("Expected INSERT statement type, got %v", stmt.StatementType)
	}
	
	// 验证UPDATE语句
	stmt, exists = config.GetMapperStatement("TestMapper.UpdateUser")
	if !exists {
		t.Fatal("UpdateUser statement should exist")
	}
	if stmt.StatementType != UPDATE {
		t.Fatalf("Expected UPDATE statement type, got %v", stmt.StatementType)
	}
	
	// 验证DELETE语句
	stmt, exists = config.GetMapperStatement("TestMapper.DeleteUser")
	if !exists {
		t.Fatal("DeleteUser statement should exist")
	}
	if stmt.StatementType != DELETE {
		t.Fatalf("Expected DELETE statement type, got %v", stmt.StatementType)
	}
}

// TestStatementType 测试语句类型常量
func TestStatementType(t *testing.T) {
	if SELECT != 0 {
		t.Fatalf("Expected SELECT to be 0, got %d", SELECT)
	}
	if INSERT != 1 {
		t.Fatalf("Expected INSERT to be 1, got %d", INSERT)
	}
	if UPDATE != 2 {
		t.Fatalf("Expected UPDATE to be 2, got %d", UPDATE)
	}
	if DELETE != 3 {
		t.Fatalf("Expected DELETE to be 3, got %d", DELETE)
	}
}

// TestMapperStatement 测试MapperStatement结构
func TestMapperStatement(t *testing.T) {
	stmt := &MapperStatement{
		ID:            "TestMapper.TestStatement",
		SQL:           "SELECT * FROM test",
		ResultType:    reflect.TypeOf(""),
		StatementType: SELECT,
	}
	
	if stmt.ID != "TestMapper.TestStatement" {
		t.Fatalf("Expected ID to be 'TestMapper.TestStatement', got %s", stmt.ID)
	}
	
	if stmt.SQL != "SELECT * FROM test" {
		t.Fatalf("Expected SQL to be 'SELECT * FROM test', got %s", stmt.SQL)
	}
	
	if stmt.ResultType != reflect.TypeOf("") {
		t.Fatalf("Expected ResultType to be string type, got %v", stmt.ResultType)
	}
	
	if stmt.StatementType != SELECT {
		t.Fatalf("Expected StatementType to be SELECT, got %v", stmt.StatementType)
	}
}

// TestDataSource 测试DataSource结构
func TestDataSource(t *testing.T) {
	ds := &DataSource{
		DriverName:     "mysql",
		DataSourceName: "user:pass@tcp(localhost:3306)/db",
		DB:             nil,
	}
	
	if ds.DriverName != "mysql" {
		t.Fatalf("Expected DriverName to be 'mysql', got %s", ds.DriverName)
	}
	
	if ds.DataSourceName != "user:pass@tcp(localhost:3306)/db" {
		t.Fatalf("Expected DataSourceName to be 'user:pass@tcp(localhost:3306)/db', got %s", ds.DataSourceName)
	}
	
	if ds.DB != nil {
		t.Fatal("Expected DB to be nil")
	}
}

// TestMapperConfig 测试MapperConfig结构
func TestMapperConfig(t *testing.T) {
	mc := &MapperConfig{
		Mappers: make(map[string]*MapperStatement),
	}
	
	if mc.Mappers == nil {
		t.Fatal("Mappers should not be nil")
	}
	
	if len(mc.Mappers) != 0 {
		t.Fatal("Mappers should be empty initially")
	}
	
	// 添加一个语句
	stmt := &MapperStatement{
		ID:            "Test.Statement",
		SQL:           "SELECT 1",
		StatementType: SELECT,
	}
	mc.Mappers["Test.Statement"] = stmt
	
	if len(mc.Mappers) != 1 {
		t.Fatalf("Expected 1 mapper, got %d", len(mc.Mappers))
	}
	
	retrievedStmt := mc.Mappers["Test.Statement"]
	if retrievedStmt != stmt {
		t.Fatal("Retrieved statement should be the same instance")
	}
}

// TestInvocation 测试Invocation结构
func TestInvocation(t *testing.T) {
	target := "test_target"
	method := reflect.Method{}
	args := []interface{}{1, "test"}
	
	invocation := &Invocation{
		Target: target,
		Method: method,
		Args:   args,
	}
	
	if invocation.Target != target {
		t.Fatal("Target should be the same")
	}
	
	if len(invocation.Args) != 2 {
		t.Fatalf("Expected 2 args, got %d", len(invocation.Args))
	}
	
	if invocation.Args[0] != 1 {
		t.Fatalf("Expected first arg to be 1, got %v", invocation.Args[0])
	}
	
	if invocation.Args[1] != "test" {
		t.Fatalf("Expected second arg to be 'test', got %v", invocation.Args[1])
	}
}