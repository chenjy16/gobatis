package mapping

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestUser 测试用户结构体
type TestUser struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"created_at"`
}

// TestUserWithPointer 带指针字段的测试用户结构体
type TestUserWithPointer struct {
	ID    int     `db:"id"`
	Name  *string `db:"name"`
	Email *string `db:"email"`
}

func TestNewResultMapper(t *testing.T) {
	mapper := NewResultMapper()

	if mapper == nil {
		t.Error("Expected non-nil mapper")
	}

	_, ok := mapper.(*DefaultResultMapper)
	if !ok {
		t.Error("Expected DefaultResultMapper type")
	}
}

func TestDefaultResultMapper_MapResult_BasicType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Mock single string result
	rows := sqlmock.NewRows([]string{"name"}).AddRow("test_name")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT name FROM users")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
	defer queryRows.Close()

	mapper := NewResultMapper()
	result, err := mapper.MapResult(queryRows, reflect.TypeOf(""))

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result != "test_name" {
		t.Errorf("Expected 'test_name', got: %v", result)
	}
}

func TestDefaultResultMapper_MapResult_EmptyResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Mock empty result
	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT name FROM users")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
	defer queryRows.Close()

	mapper := NewResultMapper()
	result, err := mapper.MapResult(queryRows, reflect.TypeOf(""))

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got: %v", result)
	}
}

func TestDefaultResultMapper_MapResults_BasicType(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Mock multiple string results
	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("name1").
		AddRow("name2").
		AddRow("name3")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT name FROM users")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
	defer queryRows.Close()

	mapper := NewResultMapper()
	results, err := mapper.MapResults(queryRows, reflect.TypeOf(""))

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got: %d", len(results))
	}

	if results[0] != "name1" || results[1] != "name2" || results[2] != "name3" {
		t.Errorf("Expected ['name1', 'name2', 'name3'], got: %v", results)
	}
}

func TestDefaultResultMapper_MapResults_Struct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Mock struct results
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "created_at"}).
		AddRow(1, "John", "john@example.com", 25, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)).
		AddRow(2, "Jane", "jane@example.com", 30, time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, email, age, created_at FROM users")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
	defer queryRows.Close()

	mapper := NewResultMapper()
	results, err := mapper.MapResults(queryRows, reflect.TypeOf(TestUser{}))

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got: %d", len(results))
	}

	user1, ok := results[0].(TestUser)
	if !ok {
		t.Error("Expected TestUser type")
	}

	if user1.ID != 1 || user1.Name != "John" || user1.Email != "john@example.com" || user1.Age != 25 {
		t.Errorf("Unexpected user1 values: %+v", user1)
	}
}

func TestDefaultResultMapper_MapResults_StructPointer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Mock struct pointer results
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "created_at"}).
		AddRow(1, "John", "john@example.com", 25, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, email, age, created_at FROM users")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
	defer queryRows.Close()

	mapper := NewResultMapper()
	results, err := mapper.MapResults(queryRows, reflect.TypeOf(&TestUser{}))

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got: %d", len(results))
	}

	user1, ok := results[0].(*TestUser)
	if !ok {
		t.Error("Expected *TestUser type")
	}

	if user1.ID != 1 || user1.Name != "John" {
		t.Errorf("Unexpected user1 values: %+v", user1)
	}
}

func TestConvertToType(t *testing.T) {
	// Test nil value
	result, err := convertToType(nil, reflect.TypeOf(""))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string, got: %v", result)
	}

	// Test assignable type
	result, err = convertToType("test", reflect.TypeOf(""))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != "test" {
		t.Errorf("Expected 'test', got: %v", result)
	}

	// Test convertible type
	result, err = convertToType(123, reflect.TypeOf(int64(0)))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != int64(123) {
		t.Errorf("Expected int64(123), got: %v", result)
	}
}

func TestConvertToFieldType(t *testing.T) {
	// Test nil value
	result, err := convertToFieldType(nil, reflect.TypeOf(""))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil, got: %v", result)
	}

	// Test time conversion from string
	timeStr := "2023-01-01 12:00:00"
	result, err = convertToFieldType(timeStr, reflect.TypeOf(time.Time{}))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if _, ok := result.(time.Time); !ok {
		t.Error("Expected time.Time type")
	}

	// Test date conversion from string
	dateStr := "2023-01-01"
	result, err = convertToFieldType(dateStr, reflect.TypeOf(time.Time{}))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if _, ok := result.(time.Time); !ok {
		t.Error("Expected time.Time type")
	}

	// Test pointer type conversion
	str := "test"
	result, err = convertToFieldType(str, reflect.TypeOf(&str))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if resultPtr, ok := result.(*string); !ok || *resultPtr != "test" {
		t.Error("Expected *string with value 'test'")
	}

	// Test assignable type
	result, err = convertToFieldType("test", reflect.TypeOf(""))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != "test" {
		t.Errorf("Expected 'test', got: %v", result)
	}

	// Test convertible type
	result, err = convertToFieldType(123, reflect.TypeOf(int64(0)))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result != int64(123) {
		t.Errorf("Expected int64(123), got: %v", result)
	}
}

func TestIsBasicType(t *testing.T) {
	basicTypes := []reflect.Kind{
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String,
	}

	for _, kind := range basicTypes {
		if !isBasicType(kind) {
			t.Errorf("Expected %s to be basic type", kind)
		}
	}

	nonBasicTypes := []reflect.Kind{
		reflect.Struct,
		reflect.Slice,
		reflect.Map,
		reflect.Ptr,
		reflect.Interface,
	}

	for _, kind := range nonBasicTypes {
		if isBasicType(kind) {
			t.Errorf("Expected %s to not be basic type", kind)
		}
	}
}

func TestCamelToSnake(t *testing.T) {
	testCases := map[string]string{
		"ID":        "i_d",
		"Name":      "name",
		"FirstName": "first_name",
		"CreatedAt": "created_at",
		"XMLData":   "x_m_l_data",
		"HTTPCode":  "h_t_t_p_code",
	}

	for input, expected := range testCases {
		result := camelToSnake(input)
		if result != expected {
			t.Errorf("camelToSnake(%s) = %s, expected %s", input, result, expected)
		}
	}
}

func TestDefaultResultMapper_ScanStruct_WithDBTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Mock results with db tag mapping
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "John", "john@example.com")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
	defer queryRows.Close()

	mapper := NewResultMapper()
	results, err := mapper.MapResults(queryRows, reflect.TypeOf(TestUser{}))

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got: %d", len(results))
	}

	user, ok := results[0].(TestUser)
	if !ok {
		t.Error("Expected TestUser type")
	}

	if user.ID != 1 || user.Name != "John" || user.Email != "john@example.com" {
		t.Errorf("Unexpected user values: %+v", user)
	}
}

func TestDefaultResultMapper_ScanStruct_MissingColumns(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Mock results with missing columns
	rows := sqlmock.NewRows([]string{"id", "name", "unknown_column"}).
		AddRow(1, "John", "ignored_value")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	queryRows, err := db.Query("SELECT id, name, unknown_column FROM users")
	if err != nil {
		t.Fatalf("Failed to execute query: %v", err)
	}
	defer queryRows.Close()

	mapper := NewResultMapper()
	results, err := mapper.MapResults(queryRows, reflect.TypeOf(TestUser{}))

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got: %d", len(results))
	}

	user, ok := results[0].(TestUser)
	if !ok {
		t.Error("Expected TestUser type")
	}

	if user.ID != 1 || user.Name != "John" {
		t.Errorf("Unexpected user values: %+v", user)
	}
}
