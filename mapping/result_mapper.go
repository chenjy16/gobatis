package mapping

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// ResultMapper 结果映射器接口
type ResultMapper interface {
	MapResult(rows *sql.Rows, resultType reflect.Type) (interface{}, error)
	MapResults(rows *sql.Rows, resultType reflect.Type) ([]interface{}, error)
}

// DefaultResultMapper 默认结果映射器
type DefaultResultMapper struct{}

// NewResultMapper 创建新的结果映射器
func NewResultMapper() ResultMapper {
	return &DefaultResultMapper{}
}

// MapResult 映射单个结果
func (m *DefaultResultMapper) MapResult(rows *sql.Rows, resultType reflect.Type) (interface{}, error) {
	results, err := m.MapResults(rows, resultType)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}

// MapResults 映射多个结果
func (m *DefaultResultMapper) MapResults(rows *sql.Rows, resultType reflect.Type) ([]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []interface{}

	for rows.Next() {
		result, err := m.scanRow(rows, columns, resultType)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return results, nil
}

// scanRow 扫描单行数据
func (m *DefaultResultMapper) scanRow(rows *sql.Rows, columns []string, resultType reflect.Type) (interface{}, error) {
	// 创建结果对象
	var result reflect.Value
	var isPtr bool

	if resultType.Kind() == reflect.Ptr {
		isPtr = true
		resultType = resultType.Elem()
		result = reflect.New(resultType)
	} else {
		result = reflect.New(resultType)
	}

	// 如果是基础类型，直接扫描
	if isBasicType(resultType.Kind()) {
		var value interface{}
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("failed to scan basic type: %w", err)
		}

		convertedValue, err := convertToType(value, resultType)
		if err != nil {
			return nil, err
		}

		if isPtr {
			ptrValue := reflect.New(resultType)
			ptrValue.Elem().Set(reflect.ValueOf(convertedValue))
			return ptrValue.Interface(), nil
		}

		return convertedValue, nil
	}

	// 如果是结构体，按字段映射
	if resultType.Kind() == reflect.Struct {
		err := m.scanStruct(rows, columns, result.Elem())
		if err != nil {
			return nil, err
		}

		if isPtr {
			return result.Interface(), nil
		}

		return result.Elem().Interface(), nil
	}

	return nil, fmt.Errorf("unsupported result type: %s", resultType.Kind())
}

// scanStruct 扫描结构体
func (m *DefaultResultMapper) scanStruct(rows *sql.Rows, columns []string, structValue reflect.Value) error {
	structType := structValue.Type()

	// 创建字段映射
	fieldMap := make(map[string]reflect.Value)
	for i := 0; i < structValue.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// 获取字段对应的列名
		columnName := field.Name
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			columnName = dbTag
		} else {
			// 转换为下划线命名
			columnName = camelToSnake(field.Name)
		}

		fieldMap[columnName] = fieldValue
	}

	// 准备扫描目标
	scanTargets := make([]interface{}, len(columns))
	scanValues := make([]reflect.Value, len(columns))

	for i, column := range columns {
		if fieldValue, exists := fieldMap[column]; exists {
			// 创建对应类型的指针用于扫描
			scanValue := reflect.New(fieldValue.Type())
			scanTargets[i] = scanValue.Interface()
			scanValues[i] = scanValue
		} else {
			// 如果没有对应字段，使用 interface{} 接收
			var dummy interface{}
			scanTargets[i] = &dummy
		}
	}

	// 扫描数据
	if err := rows.Scan(scanTargets...); err != nil {
		return fmt.Errorf("failed to scan row: %w", err)
	}

	// 设置字段值
	for i, column := range columns {
		if fieldValue, exists := fieldMap[column]; exists && scanValues[i].IsValid() {
			scannedValue := scanValues[i].Elem()
			if scannedValue.IsValid() {
				convertedValue, err := convertToFieldType(scannedValue.Interface(), fieldValue.Type())
				if err != nil {
					return fmt.Errorf("failed to convert value for field %s: %w", column, err)
				}
				if convertedValue != nil {
					fieldValue.Set(reflect.ValueOf(convertedValue))
				}
			}
		}
	}

	return nil
}

// convertToType 转换到指定类型
func convertToType(value interface{}, targetType reflect.Type) (interface{}, error) {
	if value == nil {
		return reflect.Zero(targetType).Interface(), nil
	}

	sourceValue := reflect.ValueOf(value)
	if sourceValue.Type().AssignableTo(targetType) {
		return value, nil
	}

	if sourceValue.Type().ConvertibleTo(targetType) {
		return sourceValue.Convert(targetType).Interface(), nil
	}

	return value, nil
}

// convertToFieldType 转换到字段类型
func convertToFieldType(value interface{}, fieldType reflect.Type) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	sourceValue := reflect.ValueOf(value)

	// 处理时间类型
	if fieldType == reflect.TypeOf(time.Time{}) {
		switch v := value.(type) {
		case time.Time:
			return v, nil
		case string:
			if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
				return t, nil
			}
			if t, err := time.Parse("2006-01-02", v); err == nil {
				return t, nil
			}
		}
	}

	// 处理指针类型
	if fieldType.Kind() == reflect.Ptr {
		if sourceValue.Type().AssignableTo(fieldType.Elem()) {
			ptr := reflect.New(fieldType.Elem())
			ptr.Elem().Set(sourceValue)
			return ptr.Interface(), nil
		}
	}

	// 直接赋值
	if sourceValue.Type().AssignableTo(fieldType) {
		return value, nil
	}

	// 类型转换
	if sourceValue.Type().ConvertibleTo(fieldType) {
		return sourceValue.Convert(fieldType).Interface(), nil
	}

	return value, nil
}

// isBasicType 判断是否为基础类型
func isBasicType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

// camelToSnake 驼峰转下划线
func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}