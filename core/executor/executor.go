package executor

import (
	"fmt"
	"reflect"

	"gobatis/binding"
	"gobatis/core/config"
	"gobatis/mapping"
)

// Executor SQL 执行器接口
type Executor interface {
	Query(statement *config.MapperStatement, parameter interface{}) ([]interface{}, error)
	Update(statement *config.MapperStatement, parameter interface{}) (int64, error)
}

// SimpleExecutor 简单执行器
type SimpleExecutor struct {
	configuration   *config.Configuration
	parameterBinder binding.ParameterBinder
	resultMapper    mapping.ResultMapper
}

// NewSimpleExecutor 创建简单执行器
func NewSimpleExecutor(configuration *config.Configuration) Executor {
	return &SimpleExecutor{
		configuration:   configuration,
		parameterBinder: binding.NewParameterBinder(),
		resultMapper:    mapping.NewResultMapper(),
	}
}

// Query 执行查询
func (e *SimpleExecutor) Query(statement *config.MapperStatement, parameter interface{}) ([]interface{}, error) {
	// 绑定参数
	processedSQL, args, err := e.parameterBinder.BindParameters(statement.SQL, parameter)
	if err != nil {
		return nil, fmt.Errorf("failed to bind parameters: %w", err)
	}

	// 执行查询
	rows, err := e.configuration.DataSource.DB.Query(processedSQL, args...)
	if err != nil {
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
	results, err := e.resultMapper.MapResults(rows, resultType)
	if err != nil {
		return nil, fmt.Errorf("failed to map results: %w", err)
	}

	return results, nil
}

// Update 执行更新（包括 INSERT、UPDATE、DELETE）
func (e *SimpleExecutor) Update(statement *config.MapperStatement, parameter interface{}) (int64, error) {
	// 绑定参数
	processedSQL, args, err := e.parameterBinder.BindParameters(statement.SQL, parameter)
	if err != nil {
		return 0, fmt.Errorf("failed to bind parameters: %w", err)
	}

	// 执行更新
	result, err := e.configuration.DataSource.DB.Exec(processedSQL, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute update: %w", err)
	}

	// 根据语句类型返回不同的结果
	switch statement.StatementType {
	case config.INSERT:
		// 对于 INSERT，返回插入的 ID
		if id, err := result.LastInsertId(); err == nil {
			return id, nil
		}
		// 如果获取不到 ID，返回影响的行数
		return result.RowsAffected()
	case config.UPDATE, config.DELETE:
		// 对于 UPDATE 和 DELETE，返回影响的行数
		return result.RowsAffected()
	default:
		return result.RowsAffected()
	}
}

// BatchExecutor 批量执行器
type BatchExecutor struct {
	configuration   *config.Configuration
	parameterBinder binding.ParameterBinder
	statements      []*BatchStatement
}

// BatchStatement 批量语句
type BatchStatement struct {
	Statement *config.MapperStatement
	Parameter interface{}
}

// NewBatchExecutor 创建批量执行器
func NewBatchExecutor(configuration *config.Configuration) *BatchExecutor {
	return &BatchExecutor{
		configuration:   configuration,
		parameterBinder: binding.NewParameterBinder(),
		statements:      make([]*BatchStatement, 0),
	}
}

// AddBatch 添加批量语句
func (e *BatchExecutor) AddBatch(statement *config.MapperStatement, parameter interface{}) {
	e.statements = append(e.statements, &BatchStatement{
		Statement: statement,
		Parameter: parameter,
	})
}

// ExecuteBatch 执行批量操作
func (e *BatchExecutor) ExecuteBatch() ([]int64, error) {
	if len(e.statements) == 0 {
		return nil, nil
	}

	// 开始事务
	tx, err := e.configuration.DataSource.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var results []int64
	for _, batchStmt := range e.statements {
		// 绑定参数
		processedSQL, args, err := e.parameterBinder.BindParameters(batchStmt.Statement.SQL, batchStmt.Parameter)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to bind parameters: %w", err)
		}

		// 执行语句
		result, err := tx.Exec(processedSQL, args...)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to execute batch statement: %w", err)
		}

		// 获取结果
		switch batchStmt.Statement.StatementType {
		case config.INSERT:
			if id, err := result.LastInsertId(); err == nil {
				results = append(results, id)
			} else if affected, err := result.RowsAffected(); err == nil {
				results = append(results, affected)
			} else {
				results = append(results, 0)
			}
		default:
			if affected, err := result.RowsAffected(); err == nil {
				results = append(results, affected)
			} else {
				results = append(results, 0)
			}
		}
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 清空批量语句
	e.statements = e.statements[:0]

	return results, nil
}
