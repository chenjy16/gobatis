package binding

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ParameterBinder 参数绑定器接口
type ParameterBinder interface {
	BindParameters(sql string, parameter interface{}) (string, []interface{}, error)
}

// DefaultParameterBinder 默认参数绑定器
type DefaultParameterBinder struct{}

// NewParameterBinder 创建新的参数绑定器
func NewParameterBinder() ParameterBinder {
	return &DefaultParameterBinder{}
}

// BindParameters 绑定参数
func (b *DefaultParameterBinder) BindParameters(sql string, parameter interface{}) (string, []interface{}, error) {
	if parameter == nil {
		return sql, nil, nil
	}

	// 查找所有的具名参数 #{paramName}
	re := regexp.MustCompile(`#\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(sql, -1)

	if len(matches) == 0 {
		return sql, nil, nil
	}

	var args []interface{}
	processedSQL := sql

	// 根据参数类型处理
	switch v := parameter.(type) {
	case map[string]interface{}:
		args, processedSQL = b.bindMapParameters(sql, matches, v)
	default:
		var err error
		args, processedSQL, err = b.bindStructParameters(sql, matches, parameter)
		if err != nil {
			return "", nil, err
		}
	}

	return processedSQL, args, nil
}

// bindMapParameters 绑定 Map 参数
func (b *DefaultParameterBinder) bindMapParameters(sql string, matches [][]string, params map[string]interface{}) ([]interface{}, string) {
	var args []interface{}
	processedSQL := sql

	for _, match := range matches {
		paramName := strings.TrimSpace(match[1])
		if value, exists := params[paramName]; exists {
			args = append(args, value)
			processedSQL = strings.Replace(processedSQL, match[0], "?", 1)
		} else {
			args = append(args, nil)
			processedSQL = strings.Replace(processedSQL, match[0], "?", 1)
		}
	}

	return args, processedSQL
}

// bindStructParameters 绑定结构体参数
func (b *DefaultParameterBinder) bindStructParameters(sql string, matches [][]string, parameter interface{}) ([]interface{}, string, error) {
	var args []interface{}
	processedSQL := sql

	v := reflect.ValueOf(parameter)
	t := reflect.TypeOf(parameter)

	// 如果是指针，获取实际值
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, "", fmt.Errorf("parameter is nil pointer")
		}
		v = v.Elem()
		t = t.Elem()
	}

	// 如果是基础类型，直接使用
	if isBasicType(v.Kind()) {
		for _, match := range matches {
			args = append(args, parameter)
			processedSQL = strings.Replace(processedSQL, match[0], "?", 1)
		}
		return args, processedSQL, nil
	}

	// 如果是结构体，按字段名绑定
	if v.Kind() == reflect.Struct {
		fieldMap := make(map[string]interface{})

		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			// 获取字段名，优先使用 db 标签
			fieldName := field.Name
			if dbTag := field.Tag.Get("db"); dbTag != "" {
				fieldName = dbTag
			}

			if fieldValue.CanInterface() {
				fieldMap[fieldName] = fieldValue.Interface()
			}
		}

		for _, match := range matches {
			paramName := strings.TrimSpace(match[1])
			if value, exists := fieldMap[paramName]; exists {
				args = append(args, value)
			} else {
				args = append(args, nil)
			}
			processedSQL = strings.Replace(processedSQL, match[0], "?", 1)
		}

		return args, processedSQL, nil
	}

	return nil, "", fmt.Errorf("unsupported parameter type: %T", parameter)
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

// ConvertValue 转换值类型
func ConvertValue(value interface{}, targetType reflect.Type) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	sourceValue := reflect.ValueOf(value)
	if sourceValue.Type().ConvertibleTo(targetType) {
		return sourceValue.Convert(targetType).Interface(), nil
	}

	// 特殊处理字符串到数字的转换
	if sourceValue.Kind() == reflect.String && isNumericType(targetType.Kind()) {
		str := sourceValue.String()
		switch targetType.Kind() {
		case reflect.Int, reflect.Int64:
			if i, err := strconv.ParseInt(str, 10, 64); err == nil {
				return i, nil
			}
		case reflect.Float64:
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return f, nil
			}
		}
	}

	return value, nil
}

// isNumericType 判断是否为数字类型
func isNumericType(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}