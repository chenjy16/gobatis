package executor

import (
	"errors"
	"reflect"
	"testing"

	"gobatis/core/config"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestNewSimpleExecutor 测试创建简单执行器
func TestNewSimpleExecutor(t *testing.T) {
	configuration := &config.Configuration{}
	executor := NewSimpleExecutor(configuration)

	if executor == nil {
		t.Fatal("Executor should not be nil")
	}

	simpleExec, ok := executor.(*SimpleExecutor)
	if !ok {
		t.Fatal("Executor should be SimpleExecutor type")
	}

	if simpleExec.configuration != configuration {
		t.Fatal("Configuration should be the same instance")
	}

	if simpleExec.parameterBinder == nil {
		t.Fatal("ParameterBinder should not be nil")
	}

	if simpleExec.resultMapper == nil {
		t.Fatal("ResultMapper should not be nil")
	}
}

// TestUser 测试用户结构体
type TestUser struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
}

// TestSimpleExecutor_Query_Success 测试查询成功
func TestSimpleExecutor_Query_Success(t *testing.T) {
	// 创建模拟数据库
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// 创建配置
	configuration := &config.Configuration{
		DataSource: &config.DataSource{
			DB: db,
		},
	}

	// 创建执行器
	executor := NewSimpleExecutor(configuration)

	// 创建语句
	statement := &config.MapperStatement{
		ID:            "TestMapper.GetUser",
		SQL:           "SELECT id, username FROM users WHERE id = #{id}",
		ResultType:    reflect.TypeOf(TestUser{}),
		StatementType: config.SELECT,
	}

	// 设置模拟期望
	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow(1, "john").
		AddRow(2, "jane")
	mock.ExpectQuery("SELECT id, username FROM users WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(rows)

	// 执行查询
	results, err := executor.Query(statement, map[string]interface{}{"id": 1})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// 验证模拟期望
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectations were not met: %v", err)
	}
}

// TestSimpleExecutor_Query_DatabaseError 测试数据库错误
func TestSimpleExecutor_Query_DatabaseError(t *testing.T) {
	// 创建模拟数据库
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// 创建配置
	configuration := &config.Configuration{
		DataSource: &config.DataSource{
			DB: db,
		},
	}

	// 创建执行器
	executor := NewSimpleExecutor(configuration)

	// 创建语句
	statement := &config.MapperStatement{
		ID:            "TestMapper.GetUser",
		SQL:           "SELECT id, username FROM users WHERE id = #{id}",
		ResultType:    reflect.TypeOf(TestUser{}),
		StatementType: config.SELECT,
	}

	// 设置模拟期望返回错误
	mock.ExpectQuery("SELECT id, username FROM users WHERE id = \\?").
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	// 执行查询
	_, err = executor.Query(statement, map[string]interface{}{"id": 1})
	if err == nil {
		t.Fatal("Expected database error")
	}

	// 验证模拟期望
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectations were not met: %v", err)
	}
}

// TestSimpleExecutor_Update_Insert 测试插入操作
func TestSimpleExecutor_Update_Insert(t *testing.T) {
	// 创建模拟数据库
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// 创建配置
	configuration := &config.Configuration{
		DataSource: &config.DataSource{
			DB: db,
		},
	}

	// 创建执行器
	executor := NewSimpleExecutor(configuration)

	// 创建语句
	statement := &config.MapperStatement{
		ID:            "TestMapper.InsertUser",
		SQL:           "INSERT INTO users (username) VALUES (#{username})",
		StatementType: config.INSERT,
	}

	// 设置模拟期望
	mock.ExpectExec("INSERT INTO users \\(username\\) VALUES \\(\\?\\)").
		WithArgs("john").
		WillReturnResult(sqlmock.NewResult(123, 1))

	// 执行更新
	result, err := executor.Update(statement, map[string]interface{}{"username": "john"})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if result != 123 {
		t.Fatalf("Expected insert ID 123, got %d", result)
	}

	// 验证模拟期望
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectations were not met: %v", err)
	}
}

// TestNewBatchExecutor 测试创建批量执行器
func TestNewBatchExecutor(t *testing.T) {
	configuration := &config.Configuration{}
	executor := NewBatchExecutor(configuration)

	if executor == nil {
		t.Fatal("BatchExecutor should not be nil")
	}

	if executor.configuration != configuration {
		t.Fatal("Configuration should be the same instance")
	}

	if executor.parameterBinder == nil {
		t.Fatal("ParameterBinder should not be nil")
	}

	if executor.statements == nil {
		t.Fatal("Statements should not be nil")
	}

	if len(executor.statements) != 0 {
		t.Fatal("Statements should be empty initially")
	}
}

// TestBatchExecutor_AddBatch 测试添加批量语句
func TestBatchExecutor_AddBatch(t *testing.T) {
	configuration := &config.Configuration{}
	executor := NewBatchExecutor(configuration)

	statement := &config.MapperStatement{
		ID:            "TestMapper.InsertUser",
		SQL:           "INSERT INTO users (username) VALUES (#{username})",
		StatementType: config.INSERT,
	}

	parameter := map[string]interface{}{"username": "john"}

	executor.AddBatch(statement, parameter)

	if len(executor.statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(executor.statements))
	}

	batchStmt := executor.statements[0]
	if batchStmt.Statement != statement {
		t.Fatal("Statement should be the same instance")
	}

	// 不能直接比较map，只检查类型
	if batchStmt.Parameter == nil {
		t.Fatal("Parameter should not be nil")
	}
}

// TestBatchExecutor_ExecuteBatch_Empty 测试执行空批量
func TestBatchExecutor_ExecuteBatch_Empty(t *testing.T) {
	configuration := &config.Configuration{}
	executor := NewBatchExecutor(configuration)

	results, err := executor.ExecuteBatch()
	if err != nil {
		t.Fatalf("ExecuteBatch failed: %v", err)
	}

	if results != nil {
		t.Fatal("Results should be nil for empty batch")
	}
}

// TestBatchExecutor_ExecuteBatch_Success 测试批量执行成功
func TestBatchExecutor_ExecuteBatch_Success(t *testing.T) {
	// 创建模拟数据库
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// 创建配置
	configuration := &config.Configuration{
		DataSource: &config.DataSource{
			DB: db,
		},
	}

	// 创建执行器
	executor := NewBatchExecutor(configuration)

	// 创建语句
	insertStmt := &config.MapperStatement{
		ID:            "TestMapper.InsertUser",
		SQL:           "INSERT INTO users (username) VALUES (#{username})",
		StatementType: config.INSERT,
	}

	updateStmt := &config.MapperStatement{
		ID:            "TestMapper.UpdateUser",
		SQL:           "UPDATE users SET username = #{username} WHERE id = #{id}",
		StatementType: config.UPDATE,
	}

	// 添加批量语句
	executor.AddBatch(insertStmt, map[string]interface{}{"username": "john"})
	executor.AddBatch(updateStmt, map[string]interface{}{"username": "jane", "id": 1})

	// 设置模拟期望
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO users \\(username\\) VALUES \\(\\?\\)").
		WithArgs("john").
		WillReturnResult(sqlmock.NewResult(123, 1))
	mock.ExpectExec("UPDATE users SET username = \\? WHERE id = \\?").
		WithArgs("jane", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// 执行批量操作
	results, err := executor.ExecuteBatch()
	if err != nil {
		t.Fatalf("ExecuteBatch failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0] != 123 {
		t.Fatalf("Expected insert ID 123, got %d", results[0])
	}

	if results[1] != 1 {
		t.Fatalf("Expected affected rows 1, got %d", results[1])
	}

	// 验证批量语句已清空
	if len(executor.statements) != 0 {
		t.Fatal("Statements should be cleared after execution")
	}

	// 验证模拟期望
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectations were not met: %v", err)
	}
}

// TestBatchStatement 测试BatchStatement结构
func TestBatchStatement(t *testing.T) {
	statement := &config.MapperStatement{
		ID:            "TestMapper.InsertUser",
		SQL:           "INSERT INTO users (username) VALUES (#{username})",
		StatementType: config.INSERT,
	}

	parameter := map[string]interface{}{"username": "john"}

	batchStmt := &BatchStatement{
		Statement: statement,
		Parameter: parameter,
	}

	if batchStmt.Statement != statement {
		t.Fatal("Statement should be the same instance")
	}

	if batchStmt.Parameter == nil {
		t.Fatal("Parameter should not be nil")
	}
}
