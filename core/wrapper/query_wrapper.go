package wrapper

import (
	"fmt"
	"strings"
)

// QueryWrapper 查询条件构造器
type QueryWrapper struct {
	conditions []Condition
	orderBy    []OrderBy
	groupBy    []string
	having     []Condition
	limit      *int
	offset     *int
}

// Condition 查询条件
type Condition struct {
	Column   string
	Operator string
	Value    interface{}
	Logic    string // AND, OR
}

// OrderBy 排序条件
type OrderBy struct {
	Column string
	Desc   bool
}

// NewQueryWrapper 创建新的查询构造器
func NewQueryWrapper() *QueryWrapper {
	return &QueryWrapper{
		conditions: make([]Condition, 0),
		orderBy:    make([]OrderBy, 0),
		groupBy:    make([]string, 0),
		having:     make([]Condition, 0),
	}
}

// Eq 等于条件
func (w *QueryWrapper) Eq(column string, value interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "=",
		Value:    value,
		Logic:    "AND",
	})
	return w
}

// Ne 不等于条件
func (w *QueryWrapper) Ne(column string, value interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "!=",
		Value:    value,
		Logic:    "AND",
	})
	return w
}

// Gt 大于条件
func (w *QueryWrapper) Gt(column string, value interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: ">",
		Value:    value,
		Logic:    "AND",
	})
	return w
}

// Ge 大于等于条件
func (w *QueryWrapper) Ge(column string, value interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: ">=",
		Value:    value,
		Logic:    "AND",
	})
	return w
}

// Lt 小于条件
func (w *QueryWrapper) Lt(column string, value interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "<",
		Value:    value,
		Logic:    "AND",
	})
	return w
}

// Le 小于等于条件
func (w *QueryWrapper) Le(column string, value interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "<=",
		Value:    value,
		Logic:    "AND",
	})
	return w
}

// Like 模糊查询条件
func (w *QueryWrapper) Like(column string, value interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "LIKE",
		Value:    fmt.Sprintf("%%%v%%", value),
		Logic:    "AND",
	})
	return w
}

// NotLike 不包含条件
func (w *QueryWrapper) NotLike(column string, value interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "NOT LIKE",
		Value:    fmt.Sprintf("%%%v%%", value),
		Logic:    "AND",
	})
	return w
}

// In 包含条件
func (w *QueryWrapper) In(column string, values ...interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "IN",
		Value:    values,
		Logic:    "AND",
	})
	return w
}

// NotIn 不包含条件
func (w *QueryWrapper) NotIn(column string, values ...interface{}) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "NOT IN",
		Value:    values,
		Logic:    "AND",
	})
	return w
}

// IsNull 为空条件
func (w *QueryWrapper) IsNull(column string) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "IS NULL",
		Value:    nil,
		Logic:    "AND",
	})
	return w
}

// IsNotNull 不为空条件
func (w *QueryWrapper) IsNotNull(column string) *QueryWrapper {
	w.conditions = append(w.conditions, Condition{
		Column:   column,
		Operator: "IS NOT NULL",
		Value:    nil,
		Logic:    "AND",
	})
	return w
}

// Or 或条件
func (w *QueryWrapper) Or() *QueryWrapper {
	if len(w.conditions) > 0 {
		w.conditions[len(w.conditions)-1].Logic = "OR"
	}
	return w
}

// OrderByAsc 升序排序
func (w *QueryWrapper) OrderByAsc(columns ...string) *QueryWrapper {
	for _, column := range columns {
		w.orderBy = append(w.orderBy, OrderBy{
			Column: column,
			Desc:   false,
		})
	}
	return w
}

// OrderByDesc 降序排序
func (w *QueryWrapper) OrderByDesc(columns ...string) *QueryWrapper {
	for _, column := range columns {
		w.orderBy = append(w.orderBy, OrderBy{
			Column: column,
			Desc:   true,
		})
	}
	return w
}

// GroupBy 分组
func (w *QueryWrapper) GroupBy(columns ...string) *QueryWrapper {
	w.groupBy = append(w.groupBy, columns...)
	return w
}

// Having Having 条件
func (w *QueryWrapper) Having(column, operator string, value interface{}) *QueryWrapper {
	w.having = append(w.having, Condition{
		Column:   column,
		Operator: operator,
		Value:    value,
		Logic:    "AND",
	})
	return w
}

// Limit 限制条数
func (w *QueryWrapper) Limit(limit int) *QueryWrapper {
	w.limit = &limit
	return w
}

// Offset 偏移量
func (w *QueryWrapper) Offset(offset int) *QueryWrapper {
	w.offset = &offset
	return w
}

// BuildSQL 构建 SQL 语句
func (w *QueryWrapper) BuildSQL(baseSQL string) (string, []interface{}) {
	var args []interface{}
	sql := baseSQL

	// 构建 WHERE 条件
	if len(w.conditions) > 0 {
		whereClause, whereArgs := w.buildWhereClause()
		sql += " WHERE " + whereClause
		args = append(args, whereArgs...)
	}

	// 构建 GROUP BY
	if len(w.groupBy) > 0 {
		sql += " GROUP BY " + strings.Join(w.groupBy, ", ")
	}

	// 构建 HAVING
	if len(w.having) > 0 {
		havingClause, havingArgs := w.buildHavingClause()
		sql += " HAVING " + havingClause
		args = append(args, havingArgs...)
	}

	// 构建 ORDER BY
	if len(w.orderBy) > 0 {
		orderClause := w.buildOrderByClause()
		sql += " ORDER BY " + orderClause
	}

	// 构建 LIMIT 和 OFFSET
	if w.limit != nil {
		sql += fmt.Sprintf(" LIMIT %d", *w.limit)
	}
	if w.offset != nil {
		sql += fmt.Sprintf(" OFFSET %d", *w.offset)
	}

	return sql, args
}

// buildWhereClause 构建 WHERE 子句
func (w *QueryWrapper) buildWhereClause() (string, []interface{}) {
	var clauses []string
	var args []interface{}

	for i, condition := range w.conditions {
		clause, conditionArgs := w.buildCondition(condition)

		if i > 0 {
			clause = condition.Logic + " " + clause
		}

		clauses = append(clauses, clause)
		args = append(args, conditionArgs...)
	}

	return strings.Join(clauses, " "), args
}

// buildHavingClause 构建 HAVING 子句
func (w *QueryWrapper) buildHavingClause() (string, []interface{}) {
	var clauses []string
	var args []interface{}

	for i, condition := range w.having {
		clause, conditionArgs := w.buildCondition(condition)

		if i > 0 {
			clause = condition.Logic + " " + clause
		}

		clauses = append(clauses, clause)
		args = append(args, conditionArgs...)
	}

	return strings.Join(clauses, " "), args
}

// buildCondition 构建单个条件
func (w *QueryWrapper) buildCondition(condition Condition) (string, []interface{}) {
	switch condition.Operator {
	case "IN", "NOT IN":
		values := condition.Value.([]interface{})
		placeholders := make([]string, len(values))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		clause := fmt.Sprintf("%s %s (%s)", condition.Column, condition.Operator, strings.Join(placeholders, ", "))
		return clause, values
	case "IS NULL", "IS NOT NULL":
		return fmt.Sprintf("%s %s", condition.Column, condition.Operator), []interface{}{}
	default:
		return fmt.Sprintf("%s %s ?", condition.Column, condition.Operator), []interface{}{condition.Value}
	}
}

// buildOrderByClause 构建 ORDER BY 子句
func (w *QueryWrapper) buildOrderByClause() string {
	var clauses []string
	for _, order := range w.orderBy {
		direction := "ASC"
		if order.Desc {
			direction = "DESC"
		}
		clauses = append(clauses, fmt.Sprintf("%s %s", order.Column, direction))
	}
	return strings.Join(clauses, ", ")
}
