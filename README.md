# gobatis

A ORM framework for Go, providing SQL-business logic decoupling, automatic parameter binding, struct result mapping, plugin extensions, and more.

## Features

- **SQL-Business Logic Decoupling**: Define SQL statements through XML configuration files
- **Automatic Parameter Binding**: Support for automatic binding of named parameters (`#{paramName}`)
- **Struct Result Mapping**: Automatically map query results to Go structs
- **Plugin Extension System**: Support for plugins like pagination with custom extensions
- **Dynamic Proxy**: Automatically generate Mapper interface proxies to simplify data access
- **Parameter Binding**: Support for named parameters and struct parameter binding

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