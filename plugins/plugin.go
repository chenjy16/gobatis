package plugins

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Plugin 插件接口
type Plugin interface {
	// Intercept 拦截方法调用
	Intercept(invocation *Invocation) (interface{}, error)
	// SetProperties 设置插件属性
	SetProperties(properties map[string]string)
	// GetOrder 获取插件执行顺序（数字越小优先级越高）
	GetOrder() int
}

// Invocation 拦截调用信息
type Invocation struct {
	Target      interface{}                    // 目标对象
	Method      reflect.Method                 // 调用的方法
	Args        []interface{}                  // 方法参数
	StatementId string                         // SQL 语句 ID
	Properties  map[string]interface{}         // 额外属性
	Proceed     func() (interface{}, error)    // 继续执行的函数
	Context     *InvocationContext             // 调用上下文
}

// InvocationContext 调用上下文，用于错误处理和回滚
type InvocationContext struct {
	StartTime     time.Time                    // 开始时间
	PluginStates  map[string]interface{}       // 插件状态
	RollbackFuncs []func() error               // 回滚函数列表
	mutex         sync.RWMutex                 // 保护并发访问
}

// NewInvocationContext 创建新的调用上下文
func NewInvocationContext() *InvocationContext {
	return &InvocationContext{
		StartTime:     time.Now(),
		PluginStates:  make(map[string]interface{}),
		RollbackFuncs: make([]func() error, 0),
	}
}

// SetPluginState 设置插件状态
func (ctx *InvocationContext) SetPluginState(pluginName string, state interface{}) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.PluginStates[pluginName] = state
}

// GetPluginState 获取插件状态
func (ctx *InvocationContext) GetPluginState(pluginName string) (interface{}, bool) {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	state, exists := ctx.PluginStates[pluginName]
	return state, exists
}

// AddRollbackFunc 添加回滚函数
func (ctx *InvocationContext) AddRollbackFunc(rollbackFunc func() error) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.RollbackFuncs = append(ctx.RollbackFuncs, rollbackFunc)
}

// ExecuteRollback 执行回滚
func (ctx *InvocationContext) ExecuteRollback() []error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	
	var errors []error
	// 逆序执行回滚函数
	for i := len(ctx.RollbackFuncs) - 1; i >= 0; i-- {
		if err := ctx.RollbackFuncs[i](); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// PluginChain 插件链 - 线程安全版本
type PluginChain struct {
	plugins       []Plugin
	currentIndex  int
	target        interface{}
	invocation    *Invocation
	completed     bool
	mutex         sync.RWMutex
}

// NewPluginChain 创建插件链 - 每次调用都创建新实例
func NewPluginChain(plugins []Plugin, target interface{}) *PluginChain {
	// 创建插件副本，避免并发修改
	pluginsCopy := make([]Plugin, len(plugins))
	copy(pluginsCopy, plugins)
	
	return &PluginChain{
		plugins:      pluginsCopy,
		currentIndex: 0,
		target:       target,
		completed:    false,
	}
}

// Proceed 执行插件链 - 线程安全版本
func (chain *PluginChain) Proceed(invocation *Invocation) (interface{}, error) {
	chain.mutex.Lock()
	defer chain.mutex.Unlock()
	
	// 检查是否已完成
	if chain.completed {
		return nil, fmt.Errorf("plugin chain already completed")
	}
	
	// 设置调用上下文
	if invocation.Context == nil {
		invocation.Context = NewInvocationContext()
	}
	
	chain.invocation = invocation
	
	// 保存原始的 Proceed 函数
	originalProceed := invocation.Proceed
	
	// 设置新的 Proceed 函数
	invocation.Proceed = func() (interface{}, error) {
		return chain.proceedNext(originalProceed)
	}
	
	// 开始执行插件链
	result, err := chain.proceedNext(originalProceed)
	
	// 如果发生错误，执行回滚
	if err != nil {
		rollbackErrors := invocation.Context.ExecuteRollback()
		if len(rollbackErrors) > 0 {
			// 将回滚错误附加到原始错误
			return nil, fmt.Errorf("original error: %v, rollback errors: %v", err, rollbackErrors)
		}
	}
	
	chain.completed = true
	return result, err
}

// proceedNext 执行下一个插件
func (chain *PluginChain) proceedNext(originalProceed func() (interface{}, error)) (interface{}, error) {
	if chain.currentIndex >= len(chain.plugins) {
		// 所有插件都执行完毕，调用原始方法
		return chain.invokeTarget(originalProceed)
	}

	plugin := chain.plugins[chain.currentIndex]
	chain.currentIndex++
	
	// 执行插件拦截
	return plugin.Intercept(chain.invocation)
}

// 反射调用缓存
type methodCache struct {
	cache map[string]*methodInfo
	mutex sync.RWMutex
}

type methodInfo struct {
	method     reflect.Method
	paramTypes []reflect.Type
	returnType reflect.Type
	hasError   bool
}

var globalMethodCache = &methodCache{
	cache: make(map[string]*methodInfo),
}

// getMethodInfo 获取方法信息（带缓存）
func (mc *methodCache) getMethodInfo(target interface{}, methodName string) (*methodInfo, error) {
	targetType := reflect.TypeOf(target)
	cacheKey := fmt.Sprintf("%s.%s", targetType.String(), methodName)
	
	// 先尝试从缓存获取
	mc.mutex.RLock()
	if info, exists := mc.cache[cacheKey]; exists {
		mc.mutex.RUnlock()
		return info, nil
	}
	mc.mutex.RUnlock()
	
	// 缓存未命中，使用反射获取方法信息
	method, found := targetType.MethodByName(methodName)
	if !found {
		return nil, fmt.Errorf("method %s not found on target %T", methodName, target)
	}
	
	// 构建方法信息
	methodType := method.Type
	paramTypes := make([]reflect.Type, methodType.NumIn())
	for i := 0; i < methodType.NumIn(); i++ {
		paramTypes[i] = methodType.In(i)
	}
	
	var returnType reflect.Type
	hasError := false
	if methodType.NumOut() > 0 {
		returnType = methodType.Out(0)
		if methodType.NumOut() > 1 {
			// 检查最后一个返回值是否为 error
			lastType := methodType.Out(methodType.NumOut() - 1)
			hasError = lastType.Implements(reflect.TypeOf((*error)(nil)).Elem())
		}
	}
	
	info := &methodInfo{
		method:     method,
		paramTypes: paramTypes,
		returnType: returnType,
		hasError:   hasError,
	}
	
	// 存入缓存
	mc.mutex.Lock()
	mc.cache[cacheKey] = info
	mc.mutex.Unlock()
	
	return info, nil
}

// invokeTarget 调用目标方法 - 优化版本
func (chain *PluginChain) invokeTarget(originalProceed func() (interface{}, error)) (interface{}, error) {
	// 如果有原始的 Proceed 函数，直接调用
	if originalProceed != nil {
		return originalProceed()
	}
	
	// 使用缓存的反射调用
	methodInfo, err := globalMethodCache.getMethodInfo(chain.target, chain.invocation.Method.Name)
	if err != nil {
		return nil, err
	}
	
	// 类型安全的参数准备
	args, err := chain.prepareArgs(methodInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare arguments: %w", err)
	}
	
	// 调用方法
	targetValue := reflect.ValueOf(chain.target)
	results := methodInfo.method.Func.Call(append([]reflect.Value{targetValue}, args...))
	
	// 类型安全的结果处理
	return chain.processResults(results, methodInfo)
}

// prepareArgs 准备参数 - 增强类型安全
func (chain *PluginChain) prepareArgs(methodInfo *methodInfo) ([]reflect.Value, error) {
	args := make([]reflect.Value, len(chain.invocation.Args))
	
	// 跳过第一个参数（receiver）
	paramOffset := 1
	
	for i, arg := range chain.invocation.Args {
		paramIndex := i + paramOffset
		if paramIndex >= len(methodInfo.paramTypes) {
			return nil, fmt.Errorf("too many arguments: expected %d, got %d", 
				len(methodInfo.paramTypes)-paramOffset, len(chain.invocation.Args))
		}
		
		expectedType := methodInfo.paramTypes[paramIndex]
		
		if arg == nil {
			// 处理 nil 参数
			if expectedType.Kind() == reflect.Ptr || 
			   expectedType.Kind() == reflect.Interface ||
			   expectedType.Kind() == reflect.Slice ||
			   expectedType.Kind() == reflect.Map ||
			   expectedType.Kind() == reflect.Chan ||
			   expectedType.Kind() == reflect.Func {
				args[i] = reflect.Zero(expectedType)
			} else {
				return nil, fmt.Errorf("cannot pass nil to non-pointer parameter %d of type %s", i, expectedType)
			}
		} else {
			argValue := reflect.ValueOf(arg)
			argType := argValue.Type()
			
			// 类型检查和转换
			if !argType.AssignableTo(expectedType) {
				// 尝试类型转换
				if argType.ConvertibleTo(expectedType) {
					args[i] = argValue.Convert(expectedType)
				} else {
					return nil, fmt.Errorf("argument %d: cannot convert %s to %s", i, argType, expectedType)
				}
			} else {
				args[i] = argValue
			}
		}
	}
	
	return args, nil
}

// processResults 处理返回值 - 增强类型安全
func (chain *PluginChain) processResults(results []reflect.Value, methodInfo *methodInfo) (interface{}, error) {
	if len(results) == 0 {
		return nil, nil
	}
	
	// 处理单个返回值
	if len(results) == 1 {
		result := results[0]
		if methodInfo.hasError && result.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			// 单个返回值是 error
			if result.IsNil() {
				return nil, nil
			}
			return nil, result.Interface().(error)
		}
		return result.Interface(), nil
	}
	
	// 处理多个返回值
	if len(results) == 2 && methodInfo.hasError {
		result := results[0]
		errResult := results[1]
		
		var err error
		if !errResult.IsNil() {
			err = errResult.Interface().(error)
		}
		
		return result.Interface(), err
	}
	
	// 其他情况，返回第一个值
	return results[0].Interface(), nil
}