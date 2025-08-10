package example

import (
	"fmt"
	"regexp"
	"strings"
)

// Example MyBatis 风格的查询条件构建器
type Example struct {
	oredCriteria  []Criteria
	orderByClause string
	distinct      bool
	limitStart    *int
	limitEnd      *int
}

// Criteria 查询条件组
type Criteria struct {
	criteria []Criterion
	valid    bool
}

// Criterion 单个查询条件
type Criterion struct {
	condition    string
	value        interface{}
	secondValue  interface{}
	noValue      bool
	singleValue  bool
	betweenValue bool
	listValue    bool
	typeHandler  string
}

// NewExample 创建新的 Example 实例
func NewExample() *Example {
	return &Example{
		oredCriteria: make([]Criteria, 0),
		distinct:     false,
	}
}

// isValidOrderByClause 验证ORDER BY子句是否安全
func isValidOrderByClause(orderBy string) bool {
	if orderBy == "" {
		return true
	}

	// 允许的ORDER BY模式：列名、表名.列名、ASC/DESC、逗号和空格
	// 这个正则表达式匹配安全的ORDER BY子句
	pattern := `^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?(\s+(ASC|DESC))?(\s*,\s*[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?(\s+(ASC|DESC))?)*$`
	matched, _ := regexp.MatchString(pattern, strings.TrimSpace(orderBy))
	return matched
}

// SetOrderByClause 设置排序子句
func (e *Example) SetOrderByClause(orderByClause string) {
	if isValidOrderByClause(orderByClause) {
		e.orderByClause = orderByClause
	}
	// 如果ORDER BY子句不安全，则忽略它
}

// SetDistinct 设置是否去重
func (e *Example) SetDistinct(distinct bool) {
	e.distinct = distinct
}

// SetLimit 设置分页限制
func (e *Example) SetLimit(start, end int) {
	e.limitStart = &start
	e.limitEnd = &end
}

// GetOredCriteria 获取所有条件组
func (e *Example) GetOredCriteria() []Criteria {
	return e.oredCriteria
}

// Or 添加 OR 条件组
func (e *Example) Or(criteria Criteria) {
	e.oredCriteria = append(e.oredCriteria, criteria)
}

// CreateCriteria 创建条件组
func (e *Example) CreateCriteria() *Criteria {
	criteria := &Criteria{
		criteria: make([]Criterion, 0),
		valid:    false,
	}
	// 自动将第一个 criteria 添加到 oredCriteria 中
	if len(e.oredCriteria) == 0 {
		e.oredCriteria = append(e.oredCriteria, *criteria)
		// 返回对 oredCriteria 中实际元素的引用
		return &e.oredCriteria[0]
	}
	return criteria
}

// CreateCriteriaInternal 创建内部条件组
func (e *Example) CreateCriteriaInternal() *Criteria {
	criteria := e.CreateCriteria()
	if len(e.oredCriteria) == 0 {
		e.oredCriteria = append(e.oredCriteria, *criteria)
	}
	return criteria
}

// Clear 清空所有条件
func (e *Example) Clear() {
	e.oredCriteria = make([]Criteria, 0)
	e.orderByClause = ""
	e.distinct = false
	e.limitStart = nil
	e.limitEnd = nil
}

// IsValid 检查 Example 是否有效
func (e *Example) IsValid() bool {
	return len(e.oredCriteria) > 0
}

// BuildSQL 构建 SQL 语句
func (e *Example) BuildSQL(baseSQL string) (string, []interface{}) {
	var args []interface{}
	sql := baseSQL

	// 添加 DISTINCT
	if e.distinct {
		sql = strings.Replace(sql, "SELECT", "SELECT DISTINCT", 1)
	}

	// 构建 WHERE 条件
	if e.IsValid() {
		whereClause, whereArgs := e.buildWhereClause()
		sql += " WHERE " + whereClause
		args = append(args, whereArgs...)
	}

	// 添加 ORDER BY
	if e.orderByClause != "" {
		sql += " ORDER BY " + e.orderByClause
	}

	// 添加 LIMIT
	if e.limitStart != nil && e.limitEnd != nil {
		sql += fmt.Sprintf(" LIMIT %d, %d", *e.limitStart, *e.limitEnd)
	}

	return sql, args
}

// buildWhereClause 构建 WHERE 子句
func (e *Example) buildWhereClause() (string, []interface{}) {
	var clauses []string
	var args []interface{}

	for i, criteria := range e.oredCriteria {
		if criteria.IsValid() {
			clause, criteriaArgs := criteria.buildClause()
			if i > 0 {
				clause = "OR (" + clause + ")"
			} else {
				clause = "(" + clause + ")"
			}
			clauses = append(clauses, clause)
			args = append(args, criteriaArgs...)
		}
	}

	return strings.Join(clauses, " "), args
}

// Criteria 方法实现

// IsValid 检查条件组是否有效
func (c *Criteria) IsValid() bool {
	return len(c.criteria) > 0
}

// GetCriteria 获取所有条件
func (c *Criteria) GetCriteria() []Criterion {
	return c.criteria
}

// addCriterion 添加条件
func (c *Criteria) addCriterion(condition string, value interface{}, property string) *Criteria {
	if value == nil {
		return c
	}
	c.criteria = append(c.criteria, Criterion{
		condition:   condition,
		value:       value,
		singleValue: true,
	})
	c.valid = true
	return c
}

// addCriterionForJDBCType 添加带类型的条件
func (c *Criteria) addCriterionForJDBCType(condition string, value interface{}, property string, typeHandler string) *Criteria {
	if value == nil {
		return c
	}
	c.criteria = append(c.criteria, Criterion{
		condition:   condition,
		value:       value,
		singleValue: true,
		typeHandler: typeHandler,
	})
	c.valid = true
	return c
}

// AndIsNull 添加 IS NULL 条件
func (c *Criteria) AndIsNull(property string) *Criteria {
	c.criteria = append(c.criteria, Criterion{
		condition: property + " IS NULL",
		noValue:   true,
	})
	c.valid = true
	return c
}

// AndIsNotNull 添加 IS NOT NULL 条件
func (c *Criteria) AndIsNotNull(property string) *Criteria {
	c.criteria = append(c.criteria, Criterion{
		condition: property + " IS NOT NULL",
		noValue:   true,
	})
	c.valid = true
	return c
}

// AndEqualTo 添加等于条件
func (c *Criteria) AndEqualTo(property string, value interface{}) *Criteria {
	return c.addCriterion(property+" =", value, property)
}

// AndNotEqualTo 添加不等于条件
func (c *Criteria) AndNotEqualTo(property string, value interface{}) *Criteria {
	return c.addCriterion(property+" <>", value, property)
}

// AndGreaterThan 添加大于条件
func (c *Criteria) AndGreaterThan(property string, value interface{}) *Criteria {
	return c.addCriterion(property+" >", value, property)
}

// AndGreaterThanOrEqualTo 添加大于等于条件
func (c *Criteria) AndGreaterThanOrEqualTo(property string, value interface{}) *Criteria {
	return c.addCriterion(property+" >=", value, property)
}

// AndLessThan 添加小于条件
func (c *Criteria) AndLessThan(property string, value interface{}) *Criteria {
	return c.addCriterion(property+" <", value, property)
}

// AndLessThanOrEqualTo 添加小于等于条件
func (c *Criteria) AndLessThanOrEqualTo(property string, value interface{}) *Criteria {
	return c.addCriterion(property+" <=", value, property)
}

// AndLike 添加 LIKE 条件
func (c *Criteria) AndLike(property string, value interface{}) *Criteria {
	return c.addCriterion(property+" LIKE", value, property)
}

// AndNotLike 添加 NOT LIKE 条件
func (c *Criteria) AndNotLike(property string, value interface{}) *Criteria {
	return c.addCriterion(property+" NOT LIKE", value, property)
}

// AndIn 添加 IN 条件
func (c *Criteria) AndIn(property string, values []interface{}) *Criteria {
	if len(values) == 0 {
		return c
	}
	c.criteria = append(c.criteria, Criterion{
		condition: property + " IN",
		value:     values,
		listValue: true,
	})
	c.valid = true
	return c
}

// AndNotIn 添加 NOT IN 条件
func (c *Criteria) AndNotIn(property string, values []interface{}) *Criteria {
	if len(values) == 0 {
		return c
	}
	c.criteria = append(c.criteria, Criterion{
		condition: property + " NOT IN",
		value:     values,
		listValue: true,
	})
	c.valid = true
	return c
}

// AndBetween 添加 BETWEEN 条件
func (c *Criteria) AndBetween(property string, value1, value2 interface{}) *Criteria {
	c.criteria = append(c.criteria, Criterion{
		condition:    property + " BETWEEN",
		value:        value1,
		secondValue:  value2,
		betweenValue: true,
	})
	c.valid = true
	return c
}

// AndNotBetween 添加 NOT BETWEEN 条件
func (c *Criteria) AndNotBetween(property string, value1, value2 interface{}) *Criteria {
	c.criteria = append(c.criteria, Criterion{
		condition:    property + " NOT BETWEEN",
		value:        value1,
		secondValue:  value2,
		betweenValue: true,
	})
	c.valid = true
	return c
}

// buildClause 构建条件子句
func (c *Criteria) buildClause() (string, []interface{}) {
	var clauses []string
	var args []interface{}

	for i, criterion := range c.criteria {
		if i > 0 {
			clauses = append(clauses, "AND")
		}

		if criterion.noValue {
			clauses = append(clauses, criterion.condition)
		} else if criterion.singleValue {
			clauses = append(clauses, criterion.condition+" ?")
			args = append(args, criterion.value)
		} else if criterion.betweenValue {
			clauses = append(clauses, criterion.condition+" ? AND ?")
			args = append(args, criterion.value, criterion.secondValue)
		} else if criterion.listValue {
			values := criterion.value.([]interface{})
			placeholders := make([]string, len(values))
			for j := range placeholders {
				placeholders[j] = "?"
			}
			clauses = append(clauses, criterion.condition+" ("+strings.Join(placeholders, ", ")+")")
			args = append(args, values...)
		}
	}

	return strings.Join(clauses, " "), args
}
