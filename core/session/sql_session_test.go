package session

import (
	"errors"
	"reflect"
	"testing"

	"gobatis/core/config"
)

// MockExecutor 模拟执行器
type MockExecutor struct {
	queryResult []interface{}
	queryError  error
	updateResult int64
	updateError  error
}

func (m *MockExecutor) Query(stmt *config.MapperStatement, parameter interface{}) ([]interface{}, error) {
	if m.queryError != nil {
		return nil, m.queryError
	}
	return m.queryResult, nil
}

func (m *MockExecutor) Update(stmt *config.MapperStatement, parameter interface{}) (int64, error) {
	if m.updateError != nil {
		return 0, m.updateError
	}
	return m.updateResult, nil
}

func TestNewSqlSession(t *testing.T) {
	cfg := config.NewConfiguration()
	session := NewSqlSession(cfg, true)

	if session == nil {
		t.Error("Expected non-nil session")
	}

	defaultSession, ok := session.(*DefaultSqlSession)
	if !ok {
		t.Error("Expected DefaultSqlSession type")
	}

	if defaultSession.configuration != cfg {
		t.Error("Expected configuration to be set")
	}

	if defaultSession.autoCommit != true {
		t.Error("Expected autoCommit to be true")
	}

	if defaultSession.closed != false {
		t.Error("Expected closed to be false")
	}
}

func TestDefaultSqlSession_SelectOne(t *testing.T) {
	cfg := config.NewConfiguration()
	session := &DefaultSqlSession{
		configuration: cfg,
		executor:      &MockExecutor{},
		autoCommit:    true,
		closed:        false,
	}

	// Test with closed session
	session.closed = true
	result, err := session.SelectOne("test.select", nil)
	if err == nil || err.Error() != "session is closed" {
		t.Error("Expected session closed error")
	}
	if result != nil {
		t.Error("Expected nil result")
	}

	// Reset session
	session.closed = false

	// Test with non-existent statement
	result, err = session.SelectOne("non.existent", nil)
	if err == nil || err.Error() != "statement not found: non.existent" {
		t.Error("Expected statement not found error")
	}

	// Add a SELECT statement
	stmt := &config.MapperStatement{
		ID:            "test.select",
		SQL:           "SELECT * FROM users",
		ResultType:    reflect.TypeOf(""),
		StatementType: config.SELECT,
	}
	cfg.MapperConfig.Mappers["test.select"] = stmt

	// Test with wrong statement type
	insertStmt := &config.MapperStatement{
		ID:            "test.insert",
		SQL:           "INSERT INTO users",
		ResultType:    reflect.TypeOf(""),
		StatementType: config.INSERT,
	}
	cfg.MapperConfig.Mappers["test.insert"] = insertStmt

	result, err = session.SelectOne("test.insert", nil)
	if err == nil || err.Error() != "statement test.insert is not a select statement" {
		t.Error("Expected wrong statement type error")
	}

	// Test successful query with results
	mockExecutor := &MockExecutor{
		queryResult: []interface{}{"result1", "result2"},
	}
	session.executor = mockExecutor

	result, err = session.SelectOne("test.select", nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != "result1" {
		t.Errorf("Expected 'result1', got: %v", result)
	}

	// Test with empty results
	mockExecutor.queryResult = []interface{}{}
	result, err = session.SelectOne("test.select", nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}

	// Test with executor error
	mockExecutor.queryError = errors.New("database error")
	result, err = session.SelectOne("test.select", nil)
	if err == nil || err.Error() != "database error" {
		t.Error("Expected database error")
	}
}

func TestDefaultSqlSession_SelectList(t *testing.T) {
	cfg := config.NewConfiguration()
	session := &DefaultSqlSession{
		configuration: cfg,
		executor:      &MockExecutor{},
		autoCommit:    true,
		closed:        false,
	}

	// Test with closed session
	session.closed = true
	result, err := session.SelectList("test.select", nil)
	if err == nil || err.Error() != "session is closed" {
		t.Error("Expected session closed error")
	}
	if result != nil {
		t.Error("Expected nil result")
	}

	// Reset session
	session.closed = false

	// Add a SELECT statement
	stmt := &config.MapperStatement{
		ID:            "test.select",
		SQL:           "SELECT * FROM users",
		ResultType:    reflect.TypeOf(""),
		StatementType: config.SELECT,
	}
	cfg.MapperConfig.Mappers["test.select"] = stmt

	// Test successful query
	mockExecutor := &MockExecutor{
		queryResult: []interface{}{"result1", "result2"},
	}
	session.executor = mockExecutor

	result, err = session.SelectList("test.select", nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 results, got: %d", len(result))
	}
}

func TestDefaultSqlSession_Insert(t *testing.T) {
	cfg := config.NewConfiguration()
	session := &DefaultSqlSession{
		configuration: cfg,
		executor:      &MockExecutor{},
		autoCommit:    true,
		closed:        false,
	}

	// Test with closed session
	session.closed = true
	result, err := session.Insert("test.insert", nil)
	if err == nil || err.Error() != "session is closed" {
		t.Error("Expected session closed error")
	}
	if result != 0 {
		t.Error("Expected 0 result")
	}

	// Reset session
	session.closed = false

	// Add an INSERT statement
	stmt := &config.MapperStatement{
		ID:            "test.insert",
		SQL:           "INSERT INTO users",
		ResultType:    reflect.TypeOf(""),
		StatementType: config.INSERT,
	}
	cfg.MapperConfig.Mappers["test.insert"] = stmt

	// Test successful insert
	mockExecutor := &MockExecutor{
		updateResult: 1,
	}
	session.executor = mockExecutor

	result, err = session.Insert("test.insert", nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected 1, got: %d", result)
	}
}

func TestDefaultSqlSession_Update(t *testing.T) {
	cfg := config.NewConfiguration()
	session := &DefaultSqlSession{
		configuration: cfg,
		executor:      &MockExecutor{},
		autoCommit:    true,
		closed:        false,
	}

	// Add an UPDATE statement
	stmt := &config.MapperStatement{
		ID:            "test.update",
		SQL:           "UPDATE users SET name = ?",
		ResultType:    reflect.TypeOf(""),
		StatementType: config.UPDATE,
	}
	cfg.MapperConfig.Mappers["test.update"] = stmt

	// Test successful update
	mockExecutor := &MockExecutor{
		updateResult: 2,
	}
	session.executor = mockExecutor

	result, err := session.Update("test.update", nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != 2 {
		t.Errorf("Expected 2, got: %d", result)
	}
}

func TestDefaultSqlSession_Delete(t *testing.T) {
	cfg := config.NewConfiguration()
	session := &DefaultSqlSession{
		configuration: cfg,
		executor:      &MockExecutor{},
		autoCommit:    true,
		closed:        false,
	}

	// Add a DELETE statement
	stmt := &config.MapperStatement{
		ID:            "test.delete",
		SQL:           "DELETE FROM users WHERE id = ?",
		ResultType:    reflect.TypeOf(""),
		StatementType: config.DELETE,
	}
	cfg.MapperConfig.Mappers["test.delete"] = stmt

	// Test successful delete
	mockExecutor := &MockExecutor{
		updateResult: 1,
	}
	session.executor = mockExecutor

	result, err := session.Delete("test.delete", nil)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected 1, got: %d", result)
	}
}

func TestDefaultSqlSession_GetMapper(t *testing.T) {
	cfg := config.NewConfiguration()
	session := &DefaultSqlSession{
		configuration: cfg,
		executor:      &MockExecutor{},
		autoCommit:    true,
		closed:        false,
	}

	// Test with closed session
	session.closed = true
	result := session.GetMapper(reflect.TypeOf((*interface{})(nil)).Elem())
	if result != nil {
		t.Error("Expected nil result for closed session")
	}

	// Reset session
	session.closed = false

	// Test with interface type
	interfaceType := reflect.TypeOf((*interface{})(nil)).Elem()
	result = session.GetMapper(interfaceType)
	if result == nil {
		t.Error("Expected non-nil result")
	}

	// Test with pointer type
	stringType := reflect.TypeOf("")
	result = session.GetMapper(&stringType)
	if result == nil {
		t.Error("Expected non-nil result")
	}
}

func TestDefaultSqlSession_Commit(t *testing.T) {
	cfg := config.NewConfiguration()
	session := &DefaultSqlSession{
		configuration: cfg,
		executor:      &MockExecutor{},
		autoCommit:    true,
		closed:        false,
	}

	// Test with closed session
	session.closed = true
	err := session.Commit()
	if err == nil || err.Error() != "session is closed" {
		t.Error("Expected session closed error")
	}

	// Reset session
	session.closed = false

	// Test with no transaction
	err = session.Commit()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDefaultSqlSession_Rollback(t *testing.T) {
	cfg := config.NewConfiguration()
	session := &DefaultSqlSession{
		configuration: cfg,
		executor:      &MockExecutor{},
		autoCommit:    true,
		closed:        false,
	}

	// Test with closed session
	session.closed = true
	err := session.Rollback()
	if err == nil || err.Error() != "session is closed" {
		t.Error("Expected session closed error")
	}

	// Reset session
	session.closed = false

	// Test with no transaction
	err = session.Rollback()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestDefaultSqlSession_Close(t *testing.T) {
	cfg := config.NewConfiguration()
	session := &DefaultSqlSession{
		configuration: cfg,
		executor:      &MockExecutor{},
		autoCommit:    true,
		closed:        false,
	}

	// Test closing session
	err := session.Close()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !session.closed {
		t.Error("Expected session to be closed")
	}

	// Test closing already closed session
	err = session.Close()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}