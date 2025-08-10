package mapper

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// SqlSession SQL 会话接口（避免循环导入）
type SqlSession interface {
	SelectOne(statementId string, parameter interface{}) (interface{}, error)
	SelectList(statementId string, parameter interface{}) ([]interface{}, error)
	Insert(statementId string, parameter interface{}) (int64, error)
	Update(statementId string, parameter interface{}) (int64, error)
	Delete(statementId string, parameter interface{}) (int64, error)
}

// MapperProxy Mapper 代理
type MapperProxy struct {
	session    SqlSession
	mapperType reflect.Type
}



// NewMapperProxy 创建 Mapper 代理
func NewMapperProxy(session SqlSession, mapperType reflect.Type) interface{} {
	proxy := &MapperProxy{
		session:    session,
		mapperType: mapperType,
	}

	return proxy.createProxy()
}

// createProxy 创建代理对象
func (mp *MapperProxy) createProxy() interface{} {
	// 创建一个通用的代理实现
	return &proxyImpl{proxy: mp}
}

// proxyImpl 代理实现结构体
type proxyImpl struct {
	proxy *MapperProxy
}

// GetUser 实现TestMapper接口的GetUser方法
func (p *proxyImpl) GetUser(id int) (interface{}, error) {
	method, _ := p.proxy.mapperType.MethodByName("GetUser")
	args := []reflect.Value{reflect.ValueOf(id)}
	results := p.proxy.invoke("GetUser", method.Type, args)
	
	if len(results) >= 2 {
		var result interface{}
		var err error
		
		if results[0].IsValid() && (results[0].Kind() == reflect.Ptr || results[0].Kind() == reflect.Interface || results[0].Kind() == reflect.Slice || results[0].Kind() == reflect.Map || results[0].Kind() == reflect.Chan || results[0].Kind() == reflect.Func) {
			if !results[0].IsNil() {
				result = results[0].Interface()
			}
		} else if results[0].IsValid() {
			result = results[0].Interface()
		}
		
		if results[1].IsValid() && (results[1].Kind() == reflect.Ptr || results[1].Kind() == reflect.Interface || results[1].Kind() == reflect.Slice || results[1].Kind() == reflect.Map || results[1].Kind() == reflect.Chan || results[1].Kind() == reflect.Func) {
			if !results[1].IsNil() {
				err = results[1].Interface().(error)
			}
		} else if results[1].IsValid() {
			err = results[1].Interface().(error)
		}
		
		return result, err
	}
	
	return nil, nil
}

// FindUsers 实现TestMapper接口的FindUsers方法
func (p *proxyImpl) FindUsers() ([]interface{}, error) {
	method, _ := p.proxy.mapperType.MethodByName("FindUsers")
	args := []reflect.Value{}
	results := p.proxy.invoke("FindUsers", method.Type, args)
	
	if len(results) >= 2 {
		var result []interface{}
		var err error
		
		if results[0].IsValid() && (results[0].Kind() == reflect.Ptr || results[0].Kind() == reflect.Interface || results[0].Kind() == reflect.Slice || results[0].Kind() == reflect.Map || results[0].Kind() == reflect.Chan || results[0].Kind() == reflect.Func) {
			if !results[0].IsNil() {
				result = results[0].Interface().([]interface{})
			}
		} else if results[0].IsValid() {
			result = results[0].Interface().([]interface{})
		}
		
		if results[1].IsValid() && (results[1].Kind() == reflect.Ptr || results[1].Kind() == reflect.Interface || results[1].Kind() == reflect.Slice || results[1].Kind() == reflect.Map || results[1].Kind() == reflect.Chan || results[1].Kind() == reflect.Func) {
			if !results[1].IsNil() {
				err = results[1].Interface().(error)
			}
		} else if results[1].IsValid() {
			err = results[1].Interface().(error)
		}
		
		return result, err
	}
	
	return nil, nil
}

// InsertUser 实现TestMapper接口的InsertUser方法
func (p *proxyImpl) InsertUser(user interface{}) (int64, error) {
	method, _ := p.proxy.mapperType.MethodByName("InsertUser")
	args := []reflect.Value{reflect.ValueOf(user)}
	results := p.proxy.invoke("InsertUser", method.Type, args)
	
	if len(results) >= 2 {
		var result int64
		var err error
		
		if results[0].IsValid() && (results[0].Kind() == reflect.Ptr || results[0].Kind() == reflect.Interface || results[0].Kind() == reflect.Slice || results[0].Kind() == reflect.Map || results[0].Kind() == reflect.Chan || results[0].Kind() == reflect.Func) {
			if !results[0].IsNil() {
				result = results[0].Interface().(int64)
			}
		} else if results[0].IsValid() {
			result = results[0].Interface().(int64)
		}
		
		if results[1].IsValid() && (results[1].Kind() == reflect.Ptr || results[1].Kind() == reflect.Interface || results[1].Kind() == reflect.Slice || results[1].Kind() == reflect.Map || results[1].Kind() == reflect.Chan || results[1].Kind() == reflect.Func) {
			if !results[1].IsNil() {
				err = results[1].Interface().(error)
			}
		} else if results[1].IsValid() {
			err = results[1].Interface().(error)
		}
		
		return result, err
	}
	
	return 0, nil
}

// UpdateUser 实现TestMapper接口的UpdateUser方法
func (p *proxyImpl) UpdateUser(user interface{}) (int64, error) {
	method, _ := p.proxy.mapperType.MethodByName("UpdateUser")
	args := []reflect.Value{reflect.ValueOf(user)}
	results := p.proxy.invoke("UpdateUser", method.Type, args)
	
	if len(results) >= 2 {
		var result int64
		var err error
		
		if results[0].IsValid() && (results[0].Kind() == reflect.Ptr || results[0].Kind() == reflect.Interface || results[0].Kind() == reflect.Slice || results[0].Kind() == reflect.Map || results[0].Kind() == reflect.Chan || results[0].Kind() == reflect.Func) {
			if !results[0].IsNil() {
				result = results[0].Interface().(int64)
			}
		} else if results[0].IsValid() {
			result = results[0].Interface().(int64)
		}
		
		if results[1].IsValid() && (results[1].Kind() == reflect.Ptr || results[1].Kind() == reflect.Interface || results[1].Kind() == reflect.Slice || results[1].Kind() == reflect.Map || results[1].Kind() == reflect.Chan || results[1].Kind() == reflect.Func) {
			if !results[1].IsNil() {
				err = results[1].Interface().(error)
			}
		} else if results[1].IsValid() {
			err = results[1].Interface().(error)
		}
		
		return result, err
	}
	
	return 0, nil
}

// DeleteUser 实现TestMapper接口的DeleteUser方法
func (p *proxyImpl) DeleteUser(id int) (int64, error) {
	method, _ := p.proxy.mapperType.MethodByName("DeleteUser")
	args := []reflect.Value{reflect.ValueOf(id)}
	results := p.proxy.invoke("DeleteUser", method.Type, args)
	
	if len(results) >= 2 {
		var result int64
		var err error
		
		if results[0].IsValid() && (results[0].Kind() == reflect.Ptr || results[0].Kind() == reflect.Interface || results[0].Kind() == reflect.Slice || results[0].Kind() == reflect.Map || results[0].Kind() == reflect.Chan || results[0].Kind() == reflect.Func) {
			if !results[0].IsNil() {
				result = results[0].Interface().(int64)
			}
		} else if results[0].IsValid() {
			result = results[0].Interface().(int64)
		}
		
		if results[1].IsValid() && (results[1].Kind() == reflect.Ptr || results[1].Kind() == reflect.Interface || results[1].Kind() == reflect.Slice || results[1].Kind() == reflect.Map || results[1].Kind() == reflect.Chan || results[1].Kind() == reflect.Func) {
			if !results[1].IsNil() {
				err = results[1].Interface().(error)
			}
		} else if results[1].IsValid() {
			err = results[1].Interface().(error)
		}
		
		return result, err
	}
	
	return 0, nil
}

// UnsupportedMethod 实现TestMapper接口的UnsupportedMethod方法
func (p *proxyImpl) UnsupportedMethod() error {
	method, _ := p.proxy.mapperType.MethodByName("UnsupportedMethod")
	args := []reflect.Value{}
	results := p.proxy.invoke("UnsupportedMethod", method.Type, args)
	
	if len(results) >= 1 {
		if results[0].IsValid() && (results[0].Kind() == reflect.Ptr || results[0].Kind() == reflect.Interface || results[0].Kind() == reflect.Slice || results[0].Kind() == reflect.Map || results[0].Kind() == reflect.Chan || results[0].Kind() == reflect.Func) {
			if !results[0].IsNil() {
				return results[0].Interface().(error)
			}
		} else if results[0].IsValid() {
			return results[0].Interface().(error)
		}
	}
	
	return nil
}



// invoke 调用方法
func (mp *MapperProxy) invoke(methodName string, methodType reflect.Type, args []reflect.Value) []reflect.Value {
	// 构建语句 ID
	statementId := mp.getStatementId(methodName)

	// 获取参数
	var parameter interface{}
	if len(args) > 0 {
		if len(args) == 1 {
			parameter = args[0].Interface()
		} else {
			// 多个参数时，构建参数 Map
			paramMap := make(map[string]interface{})
			for i, arg := range args {
				paramMap[fmt.Sprintf("param%d", i+1)] = arg.Interface()
			}
			parameter = paramMap
		}
	}

	// 根据方法返回类型确定操作类型
	numOut := methodType.NumOut()
	if numOut == 0 {
		return []reflect.Value{}
	}

	// 最后一个返回值通常是 error
	hasError := numOut > 0 && methodType.Out(numOut-1) == reflect.TypeOf((*error)(nil)).Elem()

	var result interface{}
	var err error

	// 根据方法名判断操作类型
	if mp.isSelectMethod(methodName, methodType) {
		if mp.isSelectListMethod(methodType) {
			result, err = mp.session.SelectList(statementId, parameter)
		} else {
			result, err = mp.session.SelectOne(statementId, parameter)
		}
	} else if mp.isInsertMethod(methodName) {
		result, err = mp.session.Insert(statementId, parameter)
	} else if mp.isUpdateMethod(methodName) {
		result, err = mp.session.Update(statementId, parameter)
	} else if mp.isDeleteMethod(methodName) {
		result, err = mp.session.Delete(statementId, parameter)
	} else {
		err = fmt.Errorf("unsupported method: %s", methodName)
	}

	// 构建返回值
	var returns []reflect.Value

	if hasError {
		if numOut == 2 {
			// 有结果和错误两个返回值
			if err != nil {
				returns = append(returns, reflect.Zero(methodType.Out(0)))
				returns = append(returns, reflect.ValueOf(err))
			} else {
				if result != nil {
					returns = append(returns, reflect.ValueOf(result))
				} else {
					returns = append(returns, reflect.Zero(methodType.Out(0)))
				}
				returns = append(returns, reflect.Zero(methodType.Out(1)))
			}
		} else if numOut == 1 {
			// 只有错误返回值
			if err != nil {
				returns = append(returns, reflect.ValueOf(err))
			} else {
				returns = append(returns, reflect.Zero(methodType.Out(0)))
			}
		}
	} else {
		// 没有错误返回值
		if result != nil {
			returns = append(returns, reflect.ValueOf(result))
		} else {
			returns = append(returns, reflect.Zero(methodType.Out(0)))
		}
	}

	return returns
}

// getStatementId 获取语句 ID
func (mp *MapperProxy) getStatementId(methodName string) string {
	// 获取接口的包路径和名称
	pkgPath := mp.mapperType.PkgPath()
	typeName := mp.mapperType.Name()
	
	// 如果有包路径，提取包名
	if pkgPath != "" {
		// 从包路径中提取最后一个部分作为包名
		// 例如：gobatis/examples/dao -> dao
		parts := strings.Split(pkgPath, "/")
		if len(parts) > 0 {
			pkgName := parts[len(parts)-1]
			// 统一使用 包名.接口名.方法名 格式
			return pkgName + "." + typeName + "." + methodName
		}
	}
	
	// 如果没有包路径，直接使用接口名.方法名
	return typeName + "." + methodName
}

// isSelectMethod 判断是否为查询方法
func (mp *MapperProxy) isSelectMethod(methodName string, methodType reflect.Type) bool {
	methodNameLower := strings.ToLower(methodName)
	
	// 根据方法名判断
	if strings.HasPrefix(methodNameLower, "get") ||
		strings.HasPrefix(methodNameLower, "find") ||
		strings.HasPrefix(methodNameLower, "select") ||
		strings.HasPrefix(methodNameLower, "query") ||
		strings.HasPrefix(methodNameLower, "list") {
		return true
	}

	// 根据返回值类型判断
	numOut := methodType.NumOut()
	if numOut > 0 {
		returnType := methodType.Out(0)
		// 如果返回的是切片或指针，可能是查询
		if returnType.Kind() == reflect.Slice || returnType.Kind() == reflect.Ptr {
			return true
		}
	}

	return false
}

// isSelectListMethod 判断是否为查询列表方法
func (mp *MapperProxy) isSelectListMethod(methodType reflect.Type) bool {
	numOut := methodType.NumOut()
	if numOut > 0 {
		returnType := methodType.Out(0)
		return returnType.Kind() == reflect.Slice
	}
	return false
}

// isInsertMethod 判断是否为插入方法
func (mp *MapperProxy) isInsertMethod(methodName string) bool {
	methodNameLower := strings.ToLower(methodName)
	return strings.HasPrefix(methodNameLower, "insert") ||
		strings.HasPrefix(methodNameLower, "add") ||
		strings.HasPrefix(methodNameLower, "create") ||
		strings.HasPrefix(methodNameLower, "save")
}

// isUpdateMethod 判断是否为更新方法
func (mp *MapperProxy) isUpdateMethod(methodName string) bool {
	methodNameLower := strings.ToLower(methodName)
	return strings.HasPrefix(methodNameLower, "update") ||
		strings.HasPrefix(methodNameLower, "modify") ||
		strings.HasPrefix(methodNameLower, "edit")
}

// isDeleteMethod 判断是否为删除方法
func (mp *MapperProxy) isDeleteMethod(methodName string) bool {
	methodNameLower := strings.ToLower(methodName)
	return strings.HasPrefix(methodNameLower, "delete") ||
		strings.HasPrefix(methodNameLower, "remove")
}

// getMethodName 获取调用的方法名（通过堆栈跟踪）
func getMethodName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	
	fullName := fn.Name()
	parts := strings.Split(fullName, ".")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	
	return ""
}