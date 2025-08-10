# gobatis

A ORM framework for Go, providing SQL-business logic decoupling, automatic parameter binding, struct result mapping, plugin extensions, and more.

## Features

- **SQL-Business Logic Decoupling**: Define SQL statements through XML configuration files
- **Automatic Parameter Binding**: Support for automatic binding of named parameters (`#{paramName}`)
- **Struct Result Mapping**: Automatically map query results to Go structs
- **Plugin Extension System**: Support for plugins like pagination with custom extensions
- **Dynamic Proxy**: Automatically generate Mapper interface proxies to simplify data access
- **Parameter Binding**: Support for named parameters and struct parameter binding
- **Advanced Logging System**: GORM-compatible logging with SQL tracing, slow query detection, and third-party logger integration

## Quick Start

### 1. Define Entity

```go
type User struct {
    ID       int64     `db:"id"`
    Username string    `db:"username"`
    Email    string    `db:"email"`
    CreateAt time.Time `db:"create_at"`
}
```

### 2. Define Mapper Interface

```go
type UserMapper interface {
    GetUserById(id int64) (*User, error)
    GetUsersByName(name string) ([]*User, error)
    GetAllUsers() ([]*User, error)
    InsertUser(user *User) (int64, error)
    UpdateUser(user *User) (int64, error)
    DeleteUser(id int64) (int64, error)
    CountUsers() (int64, error)
}
```

### 3. Configure XML Mapper

```xml
<?xml version="1.0" encoding="UTF-8"?>
<mapper namespace="UserMapper">
    <select id="GetUserById" resultType="User">
        SELECT id, username, email, create_at FROM users WHERE id = #{id}
    </select>
    
    <insert id="InsertUser">
        INSERT INTO users (username, email, create_at) 
        VALUES (#{username}, #{email}, #{createAt})
    </insert>
    
    <update id="UpdateUser">
        UPDATE users SET username = #{username}, email = #{email} WHERE id = #{id}
    </update>
    
    <delete id="DeleteUser">
        DELETE FROM users WHERE id = #{id}
    </delete>
</mapper>
```

### 4. Using the Framework

```go
package main

import (
    "fmt"
    "gobatis"
    "gobatis/examples"
    "gobatis/plugins"
)

func main() {
    // Create configuration
    config := gobatis.NewConfiguration()
    
    // Configure logging (optional - uses default logger if not set)
    config.Logger = logger.New(
        log.New(os.Stdout, "", log.LstdFlags),
        logger.Config{
            SlowThreshold: 200 * time.Millisecond,
            LogLevel:      logger.Info,
            Colorful:      true,
        },
    )
    
    // Set data source
    err := config.SetDataSource("mysql", "user:password@tcp(localhost:3306)/dbname?parseTime=true")
    if err != nil {
        panic(err)
    }
    
    // Add Mapper XML
    err = config.AddMapperXML("examples/user_mapper.xml")
    if err != nil {
        panic(err)
    }
    
    // Configure plugins
    pluginManager := plugins.NewPluginBuilder().
        WithCustomPlugin(plugins.NewPaginationPlugin()).
        Build()
    
    // Create Session
    factory := gobatis.NewSqlSessionFactory(config)
    session := factory.OpenSession()
    defer session.Close()
    
    // Get Mapper proxy
    userMapper := session.GetMapper((*examples.UserMapper)(nil)).(examples.UserMapper)
    
    // Use Mapper for CRUD operations
    user, err := userMapper.GetUserById(1)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("User: %+v\n", user)
    }
    
    // Example query builder usage
    ex := example.NewExample()
    criteria := ex.CreateCriteria()
    criteria.AndEqualTo("status", "active")
    criteria.AndGreaterThan("age", 18)
    ex.SetOrderByClause("created_at DESC")
    ex.SetLimit(0, 10)
    
    users, err := userMapper.SelectByExample(ex)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Found %d users\n", len(users))
    }
    
    // Pagination query example
    pageReq := &plugins.PageRequest{Page: 1, Size: 10}
    // pageResult := userService.SearchUsersPaginated("john", pageReq)
}
```

## Logging System

GoBatis provides a powerful and flexible logging system inspired by GORM's design, offering SQL tracing, slow query detection, multi-level logging, and third-party logger integration.

### Features

- **Interface-Driven Design** - Based on `logger.Interface`, supports custom implementations
- **Multi-Level Logging** - Silent, Error, Warn, Info levels
- **SQL Tracing** - Automatically logs SQL statements, parameters, execution time, and affected rows
- **Slow Query Detection** - Configurable threshold with automatic slow query warnings
- **Colorful Output** - Supports colored console output for better development experience
- **Context Support** - Supports `context.Context` for distributed tracing
- **High Performance** - Optimized implementation with minimal performance impact
- **Third-Party Integration** - Easy integration with Logrus, Zap, Zerolog, and other popular loggers

### Quick Start

#### 1. Using Default Logger

```go
import "gobatis/logger"

// Use default configuration
log := logger.Default

// Set log level
log = log.LogMode(logger.Info)

// Log messages
ctx := context.Background()
log.Info(ctx, "Application started")
log.Warn(ctx, "This is a warning")
log.Error(ctx, "This is an error")
```

#### 2. Custom Logger Configuration

```go
import (
    "gobatis/logger"
    "log"
    "os"
    "time"
)

// Create custom logger
customLogger := logger.New(
    log.New(os.Stdout, "[GOBATIS] ", log.LstdFlags),
    logger.Config{
        SlowThreshold:             200 * time.Millisecond, // Slow query threshold
        LogLevel:                  logger.Warn,            // Log level
        IgnoreRecordNotFoundError: true,                   // Ignore record not found errors
        Colorful:                  true,                   // Enable colored output
        ParameterizedQueries:      false,                  // Show full SQL
    },
)
```

#### 3. Using in GoBatis

```go
import (
    "gobatis"
    "gobatis/core/config"
    "gobatis/logger"
)

// Create configuration
configuration := config.NewConfiguration()

// Set custom logger
configuration.Logger = logger.New(
    log.New(os.Stdout, "", log.LstdFlags),
    logger.Config{
        SlowThreshold: 100 * time.Millisecond,
        LogLevel:      logger.Info,
        Colorful:      true,
    },
)

// Create SQL session factory
factory := gobatis.NewSqlSessionFactory(configuration)
session := factory.OpenSession()

// Queries will automatically log
user, err := session.SelectOne("UserMapper.GetUserById", 1)
```

### Log Levels

#### Silent
Silent mode, no logs output. Suitable for production environments or performance-sensitive scenarios.

```go
logger := logger.Default.LogMode(logger.Silent)
```

#### Error
Only output error logs, including SQL execution errors.

```go
logger := logger.Default.LogMode(logger.Error)
```

#### Warn
Output warning and error logs, including slow query warnings.

```go
logger := logger.Default.LogMode(logger.Warn)
```

#### Info
Output all logs, including SQL statement tracing. Suitable for development environments.

```go
logger := logger.Default.LogMode(logger.Info)
```

### SQL Tracing

The logging system automatically traces all SQL executions, recording:

- **Execution Time** - Accurate to milliseconds
- **SQL Statements** - Complete SQL statements and parameters
- **Affected Rows** - Number of query results or updated rows
- **Error Information** - Detailed error messages if execution fails

#### Example Output



### Custom Logger Implementation

You can implement the `logger.Interface` to create custom loggers:

```go
type CustomLogger struct {
    // Your fields
}

func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
    // Implement log level setting
}

func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
    // Implement info logging
}

func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
    // Implement warning logging
}

func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
    // Implement error logging
}

func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
    // Implement SQL tracing
}
```

### Third-Party Logger Integration

GoBatis logging system fully supports integration with third-party logging libraries like Logrus, Zap, Zerolog, etc. By implementing the `logger.Interface`, you can easily integrate any logging library.

#### Supported Third-Party Loggers

We provide adapters for the following popular logging libraries:

1. **Logrus** - Structured logging library
2. **Zap** - High-performance logging library  
3. **Zerolog** - Zero-allocation logging library

#### Adapter Usage

```go
package main

import (
    "gobatis/core/config"
    "gobatis/logger"
)

func main() {
    // 1. Using Logrus adapter
    logrusAdapter := &logger.LogrusAdapter{}
    config1 := config.NewConfiguration()
    config1.Logger = logrusAdapter.LogMode(logger.Info)
    
    // 2. Using Zap adapter
    zapAdapter := &logger.ZapAdapter{}
    config2 := config.NewConfiguration()
    config2.Logger = zapAdapter.LogMode(logger.Warn)
    
    // 3. Using Zerolog adapter
    zerologAdapter := &logger.ZerologAdapter{}
    config3 := config.NewConfiguration()
    config3.Logger = zerologAdapter.LogMode(logger.Error)
}
```

#### Custom Adapter Example (Zerolog)

```go
package main

import (
    "context"
    "time"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "gobatis/logger"
)

// ZerologAdapter implements logger.Interface
type ZerologAdapter struct {
    level logger.LogLevel
}

func (z *ZerologAdapter) LogMode(level logger.LogLevel) logger.Interface {
    return &ZerologAdapter{level: level}
}

func (z *ZerologAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
    if z.level >= logger.Info {
        log.Info().Msgf(msg, data...)
    }
}

func (z *ZerologAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
    if z.level >= logger.Warn {
        log.Warn().Msgf(msg, data...)
    }
}

func (z *ZerologAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
    if z.level >= logger.Error {
        log.Error().Msgf(msg, data...)
    }
}

func (z *ZerologAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    if z.level <= logger.Silent {
        return
    }
    
    elapsed := time.Since(begin)
    sql, rows := fc()
    
    if err != nil && err != logger.ErrRecordNotFound {
        log.Error().
            Err(err).
            Dur("elapsed", elapsed).
            Int64("rows", rows).
            Str("sql", sql).
            Msg("SQL execution failed")
    } else if elapsed > 200*time.Millisecond {
        log.Warn().
            Dur("elapsed", elapsed).
            Int64("rows", rows).
            Str("sql", sql).
            Msg("Slow SQL detected")
    } else if z.level >= logger.Info {
        log.Info().
            Dur("elapsed", elapsed).
            Int64("rows", rows).
            Str("sql", sql).
            Msg("SQL executed")
    }
}
```

### Configuration Options

#### SlowThreshold
Slow query threshold. Queries exceeding this time will be marked as slow queries.

```go
Config{
    SlowThreshold: 200 * time.Millisecond,
}
```

#### Colorful
Whether to enable colored output, improving readability in development environments.

```go
Config{
    Colorful: true, // Enable colored output
}
```

#### IgnoreRecordNotFoundError
Whether to ignore "record not found" error logging.

```go
Config{
    IgnoreRecordNotFoundError: true,
}
```

#### ParameterizedQueries
Whether to show parameterized queries (hide specific parameter values).

```go
Config{
    ParameterizedQueries: true, // Hide parameters, show only placeholders
}
```

### Performance Considerations

- The logging system is optimized for minimal performance impact
- Use `Silent` or `Error` levels in production environments
- SQL tracing uses lazy execution, formatting SQL only when needed
- Supports conditional compilation to completely disable logging at compile time

### Best Practices

1. **Development Environment** - Use `Info` level with colored output enabled
2. **Testing Environment** - Use `Warn` level to log slow queries and errors
3. **Production Environment** - Use `Error` level or `Silent` to minimize log output
4. **Performance Tuning** - Use slow query detection to identify performance bottlenecks
5. **Error Troubleshooting** - Temporarily increase log level when issues occur

### Example Usage

```go
// ExampleCustomLogger custom logger implementation example
type ExampleCustomLogger struct {
    *log.Logger
    LogLevel logger.LogLevel
}

// NewExampleCustomLogger creates custom logger instance
func NewExampleCustomLogger() logger.Interface {
    return &ExampleCustomLogger{
        Logger:   log.New(os.Stdout, "[GOBATIS] ", log.LstdFlags),
        LogLevel: logger.Info,
    }
}

// LogMode sets log level
func (l *ExampleCustomLogger) LogMode(level logger.LogLevel) logger.Interface {
    newLogger := *l
    newLogger.LogLevel = level
    return &newLogger
}

// Info outputs info logs
func (l *ExampleCustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
    if l.LogLevel >= logger.Info {
        l.Printf("[INFO] "+msg, data...)
    }
}

// Warn outputs warning logs
func (l *ExampleCustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
    if l.LogLevel >= logger.Warn {
        l.Printf("[WARN] "+msg, data...)
    }
}

// Error outputs error logs
func (l *ExampleCustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
    if l.LogLevel >= logger.Error {
        l.Printf("[ERROR] "+msg, data...)
    }
}

// Trace traces SQL execution
func (l *ExampleCustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
    if l.LogLevel <= logger.Silent {
        return
    }

    elapsed := time.Since(begin)
    sql, rows := fc()

    if err != nil && l.LogLevel >= logger.Error {
        l.Printf("[ERROR] [%.3fms] [rows:%v] %s | Error: %v", 
            float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
    } else if l.LogLevel >= logger.Info {
        l.Printf("[SQL] [%.3fms] [rows:%v] %s", 
            float64(elapsed.Nanoseconds())/1e6, rows, sql)
    }
}

// Example usage demonstration
func ExampleUsage() {
    // 1. Use default logger
    defaultLogger := logger.Default.LogMode(logger.Info)
    
    // 2. Create custom logger
    customLogger := NewExampleCustomLogger()
    
    // 3. Create logger with custom configuration
    configuredLogger := logger.New(log.New(os.Stdout, "[CUSTOM] ", log.LstdFlags), logger.Config{
        SlowThreshold:             100 * time.Millisecond,
        LogLevel:                  logger.Warn,
        IgnoreRecordNotFoundError: true,
        Colorful:                  false,
        ParameterizedQueries:      true,
    })

    // Usage examples
    ctx := context.Background()
    
    defaultLogger.Info(ctx, "Using default logger")
    customLogger.Warn(ctx, "Using custom logger")
    configuredLogger.Error(ctx, "Using configured logger")
    
    // SQL tracing example
    begin := time.Now()
    time.Sleep(50 * time.Millisecond) // Simulate SQL execution time
    
    defaultLogger.Trace(ctx, begin, func() (string, int64) {
        return "SELECT * FROM users WHERE id = ?", 1
    }, nil)
}
```

### Running Demos

```bash
# View basic logger demo
go run cmd/logger_demo/main.go

# View third-party logger integration demo
go run cmd/third_party_logger_demo/main.go

# Run logger tests
go test ./logger/... -v
```

## Plugin System Overview

### Example Query Builder

The Example query builder provides a type-safe way to build dynamic SQL queries without writing SQL strings directly. It supports various conditions, sorting, and pagination.

**1. Basic Usage**

```go
import "gobatis/core/example"

// Create a new Example for User table
ex := example.NewExample()

// Create criteria for conditions
criteria := ex.CreateCriteria()

// Add various conditions
criteria.AndEqualTo("username", "john")
criteria.AndGreaterThan("age", 18)
criteria.AndLike("email", "%@gmail.com")

// Set sorting
ex.SetOrderByClause("created_at DESC, username ASC")

// Set pagination
ex.SetLimit(0, 10) // LIMIT 0, 10

// Use in Mapper method
users, err := userMapper.SelectByExample(ex)
```

**2. Complex Conditions with OR**

```go
ex := example.NewExample()

// First condition group: username = 'john' AND age > 18
criteria1 := ex.CreateCriteria()
criteria1.AndEqualTo("username", "john")
criteria1.AndGreaterThan("age", 18)

// Second condition group: email LIKE '%admin%' AND status = 'active'
criteria2 := ex.Or()
criteria2.AndLike("email", "%admin%")
criteria2.AndEqualTo("status", "active")

// Final SQL: WHERE (username = ? AND age > ?) OR (email LIKE ? AND status = ?)
users, err := userMapper.SelectByExample(ex)
```

**3. Available Condition Methods**

```go
criteria := ex.CreateCriteria()

// Equality conditions
criteria.AndEqualTo("field", value)
criteria.AndNotEqualTo("field", value)

// Comparison conditions
criteria.AndGreaterThan("field", value)
criteria.AndGreaterThanOrEqualTo("field", value)
criteria.AndLessThan("field", value)
criteria.AndLessThanOrEqualTo("field", value)

// NULL conditions
criteria.AndIsNull("field")
criteria.AndIsNotNull("field")

// Pattern matching
criteria.AndLike("field", "%pattern%")
criteria.AndNotLike("field", "%pattern%")

// Range conditions
criteria.AndBetween("field", value1, value2)
criteria.AndNotBetween("field", value1, value2)

// List conditions
criteria.AndIn("field", []interface{}{value1, value2, value3})
criteria.AndNotIn("field", []interface{}{value1, value2, value3})
```

**4. Advanced Features**

```go
ex := example.NewExample()

// Enable DISTINCT
ex.SetDistinct(true)

// Set ORDER BY (with security validation)
ex.SetOrderByClause("created_at DESC, username ASC")

// Set LIMIT
ex.SetLimit(10, 20) // LIMIT 10, 20 (offset, count)

// Clear all conditions
ex.Clear()

// Check if example has valid conditions
if ex.IsValid() {
    users, err := userMapper.SelectByExample(ex)
}
```

**5. Mapper Integration**

Define methods in your Mapper interface that accept Example parameters:

```go
type UserMapper interface {
    SelectByExample(example *example.Example) ([]*User, error)
    CountByExample(example *example.Example) (int64, error)
    UpdateByExample(user *User, example *example.Example) (int64, error)
    DeleteByExample(example *example.Example) (int64, error)
}
```

**6. Security Features**

The Example query builder includes built-in security features:

- **SQL Injection Prevention**: All values are passed through parameterized queries
- **ORDER BY Validation**: ORDER BY clauses are validated to prevent injection
- **Safe Column Names**: Only valid column name patterns are allowed

```go
// Safe - will be accepted
ex.SetOrderByClause("name ASC, created_at DESC")

// Unsafe - will be ignored
ex.SetOrderByClause("name; DROP TABLE users;")
```

### Pagination Plugin

The pagination plugin can automatically intercept queries with pagination parameters and return paginated results.

**1. Define Mapper Method**

Define a method in the Mapper interface whose parameters include `*plugins.PageRequest`.

```go
type UserMapper interface {
    // ... other methods
    FindUsers(name string, pageReq *plugins.PageRequest) ([]*User, error)
}
```

**2. Configure Mapper XML**

The corresponding XML statement does not need to include pagination logic.

```xml
<select id="FindUsers" resultType="User">
    SELECT id, username, email, create_at 
    FROM users 
    WHERE username LIKE #{name}
</select>
```

**3. Call Pagination Query**

In business code, create a `PageRequest` object and call the Mapper method.

```go
// Add pagination plugin
pluginManager := plugins.NewPluginBuilder().
    WithCustomPlugin(plugins.NewPaginationPlugin()).
    Build()

// ... (get session and mapper)

// Create pagination request
pageReq := &plugins.PageRequest{
    Page:    1,        // Page number (starting from 1)
    Size:    10,       // Page size
    SortBy:  "id",     // Sort field
    SortDir: "ASC",    // Sort direction
}

// Execute pagination query
// Plugin will automatically modify SQL to add LIMIT/OFFSET and ORDER BY
// Return type is *plugins.PageResult
pageResult, err := userMapper.FindUsers("test", pageReq)
if err != nil {
    // ... handle error
}

// Handle pagination results
fmt.Printf("Current page: %d, Total pages: %d, Total records: %d\n", 
    pageResult.Page, pageResult.TotalPages, pageResult.Total)

for _, user := range pageResult.Data.([]*User) {
    fmt.Printf("  - User: %+v\n", user)
}
```

The pagination plugin automatically completes the following tasks:
1.  Execute `COUNT(*)` query to get the total number of records.
2.  Modify the original SQL to add `ORDER BY`, `LIMIT`, and `OFFSET` clauses.
3.  Execute the query and return `*plugins.PageResult`, which contains paginated data and metadata.



### Custom Plugin

```go
type MyPlugin struct {
    order int
}

func (p *MyPlugin) Intercept(invocation *plugins.Invocation) (interface{}, error) {
    // Pre-processing
    fmt.Println("Before method execution")
    
    // Call next plugin or target method
    result, err := invocation.Proceed()
    
    // Post-processing
    fmt.Println("After method execution")
    
    return result, err
}

func (p *MyPlugin) SetProperties(properties map[string]string) {
    // Set plugin properties
}

func (p *MyPlugin) GetOrder() int {
    return p.order // Return execution order
}
```

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run plugin tests
go test -v ./plugins
```

## Summary



### ✅ Implemented Features

1. **Configuration Management System**
   - XML configuration file parsing
   - Data source configuration
   - Mapper statement management

2. **SQL Session Management**
   - SqlSession interface
   - Connection pool management
   - Transaction control

3. **Dynamic Proxy System**
   - Automatic interface proxy
   - Method call routing
   - Parameter binding

4. **Example Query Builder**
   - Type-safe dynamic SQL query construction
   - Support for complex conditions (AND, OR, LIKE, IN, BETWEEN, etc.)
   - Built-in SQL injection prevention
   - ORDER BY and LIMIT clause support
   - Security validation for all user inputs

5. **Plugin Extension System**
   - Pagination plugin (automatic pagination queries and counting)
   - Plugin manager (plugin registration, sorting, execution chain)
   - Security features (parameter validation, SQL injection prevention)

6. **Parameter Binding and Result Mapping**
   - Named parameter binding
   - Struct field mapping
   - Type conversion

### 🎯 Technical Highlights

- **Plugin Architecture**: Uses interceptor pattern, supports plugin chain execution
- **Concurrency Safety**: Plugin manager supports concurrent access
- **Flexible Configuration**: Supports both XML configuration and code configuration
- **Test Coverage**: Complete unit tests and integration tests
- **Performance Optimization**: Connection pooling, batch operation support





## Core Components

### 1. Configuration Management (Configuration)
- Data source configuration
- Mapper XML parsing
- Plugin management

### 2. Session Management (SqlSession)
- Database connection management
- Transaction control
- Mapper proxy creation

### 3. Example Query Builder (Example)
- Type-safe dynamic SQL query construction
- Complex condition building (AND, OR, comparison, pattern matching)
- Security validation and SQL injection prevention
- ORDER BY and LIMIT clause support

### 4. Parameter Binding (ParameterBinder)
- Named parameter binding
- Struct field mapping
- Type conversion

### 5. Result Mapping (ResultMapper)
- Query result to struct mapping
- Field name conversion (camelCase ↔ snake_case)
- Type conversion

### 6. Dynamic Proxy (MapperProxy)
- Interface method proxy
- Method call routing
- Return value handling

### 7. SQL Executor (Executor)
- SQL execution
- Parameter binding
- Result processing

### 8. Plugin System (Plugins)
- **Pagination Plugin**: Automatic pagination queries with sorting and counting support
- **Plugin Manager**: Plugin registration, sorting, and execution chain management
- **Security Features**: Parameter validation and SQL injection prevention

## Project Structure

```
gobatis/
├── binding/              # Parameter binding module
│   ├── parameter_binder.go
│   └── parameter_binder_test.go
├── core/                 # Core modules
│   ├── config/          # Configuration management
│   │   ├── configuration.go
│   │   └── configuration_test.go
│   ├── example/         # Example query builder
│   │   ├── example.go   # Example implementation
│   │   ├── example_test.go # Example tests
│   │   └── README.md    # Example documentation
│   ├── executor/        # SQL executor
│   │   ├── executor.go
│   │   └── executor_test.go
│   ├── mapper/          # Mapper proxy
│   │   ├── mapper_proxy.go
│   │   └── mapper_proxy_test.go
│   └── session/         # Session management
│       ├── sql_session.go
│       └── sql_session_test.go
├── plugins/             # Plugin system
│   ├── manager.go      # Plugin manager
│   ├── pagination.go   # Pagination plugin
│   ├── plugin.go       # Plugin interface
│   └── plugins_test.go # Plugin tests
├── examples/            # Example code
│   ├── complete_example.go # Complete usage example
│   ├── complete_example_test.go
│   ├── dao.go          # Data access objects
│   ├── main.go         # Main example
│   ├── models.go       # Model definitions
│   ├── schema.sql      # Database schema
│   └── README.md       # Examples documentation
├── mapping/             # Result mapping module
│   ├── result_mapper.go
│   └── result_mapper_test.go
├── gobatis.go          # Main entry file
├── go.mod              # Go module file
├── go.sum              # Go dependencies
├── README.md           # Project documentation
└── SECURITY.md         # Security guidelines
```


## Technical Implementation

- **Dynamic Proxy**: Use `reflect.MakeFunc` to create interface proxies
- **SQL Parsing**: Support named parameter parsing and binding
- **Result Mapping**: Automatically map query results to Go structs
- **Connection Pool**: Connection pool management based on `database/sql`

## Dependencies

- `database/sql`: Go standard database interface
- `github.com/go-sql-driver/mysql`: MySQL driver (optional)

## License

MIT License