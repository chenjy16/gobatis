package session

import (
	"database/sql"
	"fmt"
	"reflect"

	"gobatis/core/config"
	"gobatis/core/executor"
	"gobatis/core/mapper"
)

// SqlSession SQL 会话接口
type SqlSession interface {
	SelectOne(statementId string, parameter interface{}) (interface{}, error)
	SelectList(statementId string, parameter interface{}) ([]interface{}, error)
	Insert(statementId string, parameter interface{}) (int64, error)
	Update(statementId string, parameter interface{}) (int64, error)
	Delete(statementId string, parameter interface{}) (int64, error)
	GetMapper(mapperType interface{}) interface{}
	Commit() error
	Rollback() error
	Close() error
}

// DefaultSqlSession 默认 SQL 会话实现
type DefaultSqlSession struct {
	configuration *config.Configuration
	executor      executor.Executor
	tx            *sql.Tx
	autoCommit    bool
	closed        bool
}

// NewSqlSession 创建新的 SQL 会话
func NewSqlSession(configuration *config.Configuration, autoCommit bool) SqlSession {
	exec := executor.NewSimpleExecutor(configuration)
	return &DefaultSqlSession{
		configuration: configuration,
		executor:      exec,
		autoCommit:    autoCommit,
		closed:        false,
	}
}

// SelectOne 查询单个结果
func (s *DefaultSqlSession) SelectOne(statementId string, parameter interface{}) (interface{}, error) {
	if s.closed {
		return nil, fmt.Errorf("session is closed")
	}

	stmt, exists := s.configuration.GetMapperStatement(statementId)
	if !exists {
		return nil, fmt.Errorf("statement not found: %s", statementId)
	}

	if stmt.StatementType != config.SELECT {
		return nil, fmt.Errorf("statement %s is not a select statement", statementId)
	}

	results, err := s.executor.Query(stmt, parameter)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}

// SelectList 查询多个结果
func (s *DefaultSqlSession) SelectList(statementId string, parameter interface{}) ([]interface{}, error) {
	if s.closed {
		return nil, fmt.Errorf("session is closed")
	}

	stmt, exists := s.configuration.GetMapperStatement(statementId)
	if !exists {
		return nil, fmt.Errorf("statement not found: %s", statementId)
	}

	if stmt.StatementType != config.SELECT {
		return nil, fmt.Errorf("statement %s is not a select statement", statementId)
	}

	return s.executor.Query(stmt, parameter)
}

// Insert 插入数据
func (s *DefaultSqlSession) Insert(statementId string, parameter interface{}) (int64, error) {
	if s.closed {
		return 0, fmt.Errorf("session is closed")
	}

	stmt, exists := s.configuration.GetMapperStatement(statementId)
	if !exists {
		return 0, fmt.Errorf("statement not found: %s", statementId)
	}

	if stmt.StatementType != config.INSERT {
		return 0, fmt.Errorf("statement %s is not an insert statement", statementId)
	}

	return s.executor.Update(stmt, parameter)
}

// Update 更新数据
func (s *DefaultSqlSession) Update(statementId string, parameter interface{}) (int64, error) {
	if s.closed {
		return 0, fmt.Errorf("session is closed")
	}

	stmt, exists := s.configuration.GetMapperStatement(statementId)
	if !exists {
		return 0, fmt.Errorf("statement not found: %s", statementId)
	}

	if stmt.StatementType != config.UPDATE {
		return 0, fmt.Errorf("statement %s is not an update statement", statementId)
	}

	return s.executor.Update(stmt, parameter)
}

// Delete 删除数据
func (s *DefaultSqlSession) Delete(statementId string, parameter interface{}) (int64, error) {
	if s.closed {
		return 0, fmt.Errorf("session is closed")
	}

	stmt, exists := s.configuration.GetMapperStatement(statementId)
	if !exists {
		return 0, fmt.Errorf("statement not found: %s", statementId)
	}

	if stmt.StatementType != config.DELETE {
		return 0, fmt.Errorf("statement %s is not a delete statement", statementId)
	}

	return s.executor.Update(stmt, parameter)
}

// GetMapper 获取 Mapper 代理
func (s *DefaultSqlSession) GetMapper(mapperType interface{}) interface{} {
	if s.closed {
		return nil
	}

	var t reflect.Type
	switch v := mapperType.(type) {
	case reflect.Type:
		t = v
	default:
		t = reflect.TypeOf(mapperType)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	return mapper.NewMapperProxy(s, t)
}

// Commit 提交事务
func (s *DefaultSqlSession) Commit() error {
	if s.closed {
		return fmt.Errorf("session is closed")
	}

	if s.tx != nil {
		err := s.tx.Commit()
		s.tx = nil
		return err
	}

	return nil
}

// Rollback 回滚事务
func (s *DefaultSqlSession) Rollback() error {
	if s.closed {
		return fmt.Errorf("session is closed")
	}

	if s.tx != nil {
		err := s.tx.Rollback()
		s.tx = nil
		return err
	}

	return nil
}

// Close 关闭会话
func (s *DefaultSqlSession) Close() error {
	if s.closed {
		return nil
	}

	if s.tx != nil {
		s.tx.Rollback()
		s.tx = nil
	}

	s.closed = true
	return nil
}