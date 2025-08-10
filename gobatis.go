package gobatis

import (
	"context"
	"database/sql"
	"fmt"
	"gobatis/binding"
	"gobatis/core/config"
	"gobatis/core/mapper"
	"gobatis/mapping"
	"gobatis/plugins"
	"reflect"
	"time"
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

// SqlSessionFactory SQL 会话工厂接口
type SqlSessionFactory interface {
	OpenSession() SqlSession
	OpenSessionWithAutoCommit(autoCommit bool) SqlSession
}

// DefaultSqlSession 默认 SQL 会话实现
type DefaultSqlSession struct {
	configuration   *config.Configuration
	parameterBinder binding.ParameterBinder
	resultMapper    mapping.ResultMapper
	pluginManager   *plugins.PluginManager
	tx              *sql.Tx
	autoCommit      bool
	closed          bool
}

// DefaultSqlSessionFactory 默认 SQL 会话工厂
type DefaultSqlSessionFactory struct {
	configuration *config.Configuration
	pluginManager *plugins.PluginManager
}

// NewConfiguration 创建新的配置
func NewConfiguration() *config.Configuration {
	return config.NewConfiguration()
}

// NewSqlSessionFactory 创建 SQL 会话工厂
func NewSqlSessionFactory(configuration *config.Configuration) SqlSessionFactory {
	return &DefaultSqlSessionFactory{
		configuration: configuration,
		pluginManager: plugins.NewPluginManager(),
	}
}

// NewSqlSessionFactoryWithPlugins 创建带插件管理器的 SQL 会话工厂
func NewSqlSessionFactoryWithPlugins(configuration *config.Configuration, pluginManager *plugins.PluginManager) SqlSessionFactory {
	return &DefaultSqlSessionFactory{
		configuration: configuration,
		pluginManager: pluginManager,
	}
}

// OpenSession 打开会话（自动提交）
func (f *DefaultSqlSessionFactory) OpenSession() SqlSession {
	return f.OpenSessionWithAutoCommit(true)
}

// OpenSessionWithAutoCommit 打开会话（指定是否自动提交）
func (f *DefaultSqlSessionFactory) OpenSessionWithAutoCommit(autoCommit bool) SqlSession {
	return &DefaultSqlSession{
		configuration:   f.configuration,
		parameterBinder: binding.NewParameterBinder(),
		resultMapper:    mapping.NewResultMapper(),
		pluginManager:   f.pluginManager,
		autoCommit:      autoCommit,
		closed:          false,
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

	// 创建拦截调用
	invocation := &plugins.Invocation{
		Target:      s,
		Method:      reflect.Method{Name: "SelectOne"},
		Args:        []interface{}{statementId, parameter},
		StatementId: statementId,
		Properties:  map[string]interface{}{"sql": stmt.SQL},
		Proceed: func() (interface{}, error) {
			results, err := s.query(stmt, parameter)
			if err != nil {
				return nil, err
			}
			if len(results) == 0 {
				return nil, nil
			}
			return results[0], nil
		},
	}

	// 如果有插件管理器，使用插件拦截
	if s.pluginManager != nil && s.pluginManager.Size() > 0 {
		method := reflect.Method{Name: "SelectOne"}
		return s.pluginManager.InterceptMethod(s, method, []interface{}{statementId, parameter}, statementId, invocation.Proceed)
	}

	// 否则直接执行
	return invocation.Proceed()
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

	// 创建拦截调用
	invocation := &plugins.Invocation{
		Target:      s,
		Method:      reflect.Method{Name: "SelectList"},
		Args:        []interface{}{statementId, parameter},
		StatementId: statementId,
		Properties:  map[string]interface{}{"sql": stmt.SQL},
		Proceed: func() (interface{}, error) {
			return s.query(stmt, parameter)
		},
	}

	// 如果有插件管理器，使用插件拦截
	if s.pluginManager != nil && s.pluginManager.Size() > 0 {
		method := reflect.Method{Name: "SelectList"}
		result, err := s.pluginManager.InterceptMethod(s, method, []interface{}{statementId, parameter}, statementId, invocation.Proceed)
		if err != nil {
			return nil, err
		}
		// 确保返回类型是 []interface{}
		if results, ok := result.([]interface{}); ok {
			return results, nil
		}
		return nil, fmt.Errorf("unexpected result type from plugin: %T", result)
	}

	// 否则直接执行
	result, err := invocation.Proceed()
	if err != nil {
		return nil, err
	}
	return result.([]interface{}), nil
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

	return s.update(stmt, parameter)
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

	return s.update(stmt, parameter)
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

	return s.update(stmt, parameter)
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

// query 执行查询
func (s *DefaultSqlSession) query(statement *config.MapperStatement, parameter interface{}) ([]interface{}, error) {
	// 开始计时
	begin := time.Now()
	ctx := context.Background()

	// 绑定参数
	processedSQL, args, err := s.parameterBinder.BindParameters(statement.SQL, parameter)
	if err != nil {
		// 记录参数绑定错误
		s.configuration.Logger.Trace(ctx, begin, func() (string, int64) {
			return fmt.Sprintf("%s [PARAMS: %v]", statement.SQL, parameter), -1
		}, err)
		return nil, fmt.Errorf("failed to bind parameters: %w", err)
	}

	// 执行查询
	rows, err := s.configuration.DataSource.DB.Query(processedSQL, args...)
	if err != nil {
		// 记录查询执行错误
		s.configuration.Logger.Trace(ctx, begin, func() (string, int64) {
			return fmt.Sprintf("%s [ARGS: %v]", processedSQL, args), -1
		}, err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// 确定结果类型
	resultType := statement.ResultType
	if resultType == nil {
		// 如果没有指定结果类型，使用 map[string]interface{}
		resultType = reflect.TypeOf(map[string]interface{}{})
	}

	// 映射结果
	results, err := s.resultMapper.MapResults(rows, resultType)
	if err != nil {
		// 记录结果映射错误
		s.configuration.Logger.Trace(ctx, begin, func() (string, int64) {
			return fmt.Sprintf("%s [ARGS: %v]", processedSQL, args), -1
		}, err)
		return nil, fmt.Errorf("failed to map results: %w", err)
	}

	// 记录成功的查询
	s.configuration.Logger.Trace(ctx, begin, func() (string, int64) {
		return fmt.Sprintf("%s [ARGS: %v]", processedSQL, args), int64(len(results))
	}, nil)

	return results, nil
}

// update 执行更新（包括 INSERT、UPDATE、DELETE）
func (s *DefaultSqlSession) update(statement *config.MapperStatement, parameter interface{}) (int64, error) {
	// 开始计时
	begin := time.Now()
	ctx := context.Background()

	// 绑定参数
	processedSQL, args, err := s.parameterBinder.BindParameters(statement.SQL, parameter)
	if err != nil {
		// 记录参数绑定错误
		s.configuration.Logger.Trace(ctx, begin, func() (string, int64) {
			return fmt.Sprintf("%s [PARAMS: %v]", statement.SQL, parameter), -1
		}, err)
		return 0, fmt.Errorf("failed to bind parameters: %w", err)
	}

	// 执行更新
	result, err := s.configuration.DataSource.DB.Exec(processedSQL, args...)
	if err != nil {
		// 记录执行错误
		s.configuration.Logger.Trace(ctx, begin, func() (string, int64) {
			return fmt.Sprintf("%s [ARGS: %v]", processedSQL, args), -1
		}, err)
		return 0, fmt.Errorf("failed to execute update: %w", err)
	}

	// 根据语句类型返回不同的结果
	var affectedRows int64
	switch statement.StatementType {
	case config.INSERT:
		// 对于 INSERT，返回插入的 ID
		if id, err := result.LastInsertId(); err == nil {
			affectedRows = id
		} else if affected, err := result.RowsAffected(); err == nil {
			affectedRows = affected
		}
	case config.UPDATE, config.DELETE:
		// 对于 UPDATE 和 DELETE，返回影响的行数
		if affected, err := result.RowsAffected(); err == nil {
			affectedRows = affected
		}
	default:
		if affected, err := result.RowsAffected(); err == nil {
			affectedRows = affected
		}
	}

	// 记录成功的更新
	s.configuration.Logger.Trace(ctx, begin, func() (string, int64) {
		return fmt.Sprintf("%s [ARGS: %v]", processedSQL, args), affectedRows
	}, nil)

	return affectedRows, nil
}
