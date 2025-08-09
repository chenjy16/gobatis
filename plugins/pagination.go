package plugins

import (
	"fmt"
	"reflect"
	"strings"
)

// PageRequest 分页请求
type PageRequest struct {
	Page     int `json:"page"`     // 页码（从1开始）
	Size     int `json:"size"`     // 每页大小
	Offset   int `json:"offset"`   // 偏移量
	SortBy   string `json:"sortBy"` // 排序字段
	SortDir  string `json:"sortDir"` // 排序方向（ASC/DESC）
}

// PageResult 分页结果
type PageResult struct {
	Data       interface{} `json:"data"`       // 数据列表
	Total      int64       `json:"total"`      // 总记录数
	Page       int         `json:"page"`       // 当前页码
	Size       int         `json:"size"`       // 每页大小
	TotalPages int         `json:"totalPages"` // 总页数
	HasNext    bool        `json:"hasNext"`    // 是否有下一页
	HasPrev    bool        `json:"hasPrev"`    // 是否有上一页
}

// PaginationPlugin 分页插件
type PaginationPlugin struct {
	properties map[string]string
	order      int
}

// NewPaginationPlugin 创建分页插件
func NewPaginationPlugin() *PaginationPlugin {
	return &PaginationPlugin{
		properties: make(map[string]string),
		order:      100, // 较低优先级，在其他插件之后执行
	}
}

// Intercept 拦截方法调用
func (p *PaginationPlugin) Intercept(invocation *Invocation) (interface{}, error) {
	// 检查是否需要分页
	pageRequest := p.extractPageRequest(invocation.Args)
	if pageRequest == nil {
		// 不需要分页，继续执行
		return invocation.Proceed()
	}

	// 修改 SQL 添加分页
	originalSQL := p.getOriginalSQL(invocation)
	if originalSQL == "" {
		return invocation.Proceed()
	}

	// 先查询总数
	countSQL := p.buildCountSQL(originalSQL)
	total, err := p.executeCountQuery(invocation, countSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to execute count query: %w", err)
	}

	// 构建分页 SQL
	pagedSQL := p.buildPagedSQL(originalSQL, pageRequest)
	
	// 更新调用参数中的 SQL
	p.updateSQL(invocation, pagedSQL)

	// 执行分页查询
	result, err := invocation.Proceed()
	if err != nil {
		return nil, err
	}

	// 构建分页结果
	pageResult := &PageResult{
		Data:       result,
		Total:      total,
		Page:       pageRequest.Page,
		Size:       pageRequest.Size,
		TotalPages: int((total + int64(pageRequest.Size) - 1) / int64(pageRequest.Size)),
		HasNext:    pageRequest.Page < int((total+int64(pageRequest.Size)-1)/int64(pageRequest.Size)),
		HasPrev:    pageRequest.Page > 1,
	}

	return pageResult, nil
}

// SetProperties 设置插件属性
func (p *PaginationPlugin) SetProperties(properties map[string]string) {
	p.properties = properties
}

// GetOrder 获取插件执行顺序
func (p *PaginationPlugin) GetOrder() int {
	return p.order
}

// extractPageRequest 从参数中提取分页请求
func (p *PaginationPlugin) extractPageRequest(args []interface{}) *PageRequest {
	for _, arg := range args {
		if pageReq, ok := arg.(*PageRequest); ok {
			// 计算偏移量
			if pageReq.Offset == 0 && pageReq.Page > 0 {
				pageReq.Offset = (pageReq.Page - 1) * pageReq.Size
			}
			return pageReq
		}
		
		// 检查是否为包含分页信息的结构体
		v := reflect.ValueOf(arg)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			pageField := v.FieldByName("Page")
			sizeField := v.FieldByName("Size")
			if pageField.IsValid() && sizeField.IsValid() {
				page := int(pageField.Int())
				size := int(sizeField.Int())
				if page > 0 && size > 0 {
					return &PageRequest{
						Page:   page,
						Size:   size,
						Offset: (page - 1) * size,
					}
				}
			}
		}
	}
	return nil
}

// getOriginalSQL 获取原始 SQL
func (p *PaginationPlugin) getOriginalSQL(invocation *Invocation) string {
	if sql, exists := invocation.Properties["sql"]; exists {
		if sqlStr, ok := sql.(string); ok {
			return sqlStr
		}
	}
	return ""
}

// buildCountSQL 构建计数 SQL
func (p *PaginationPlugin) buildCountSQL(originalSQL string) string {
	// 简单的计数 SQL 构建
	// 实际项目中可能需要更复杂的 SQL 解析
	lowerSQL := strings.ToLower(strings.TrimSpace(originalSQL))
	
	// 查找 FROM 子句
	fromIndex := strings.Index(lowerSQL, "from")
	if fromIndex == -1 {
		return ""
	}
	
	// 查找 ORDER BY 子句并移除
	orderByIndex := strings.LastIndex(lowerSQL, "order by")
	fromClause := originalSQL[fromIndex:]
	if orderByIndex > fromIndex {
		fromClause = originalSQL[fromIndex:orderByIndex]
	}
	
	return fmt.Sprintf("SELECT COUNT(*) %s", fromClause)
}

// buildPagedSQL 构建分页 SQL
func (p *PaginationPlugin) buildPagedSQL(originalSQL string, pageRequest *PageRequest) string {
	sql := originalSQL
	
	// 添加排序
	if pageRequest.SortBy != "" {
		sortDir := "ASC"
		if strings.ToUpper(pageRequest.SortDir) == "DESC" {
			sortDir = "DESC"
		}
		
		// 检查是否已有 ORDER BY
		if !strings.Contains(strings.ToLower(sql), "order by") {
			sql += fmt.Sprintf(" ORDER BY %s %s", pageRequest.SortBy, sortDir)
		}
	}
	
	// 添加 LIMIT 和 OFFSET
	sql += fmt.Sprintf(" LIMIT %d OFFSET %d", pageRequest.Size, pageRequest.Offset)
	
	return sql
}

// executeCountQuery 执行计数查询
func (p *PaginationPlugin) executeCountQuery(invocation *Invocation, countSQL string) (int64, error) {
	// 这里需要实际的数据库执行逻辑
	// 在实际实现中，需要访问 SqlSession 来执行查询
	// 暂时返回模拟数据
	return 100, nil
}

// updateSQL 更新调用中的 SQL
func (p *PaginationPlugin) updateSQL(invocation *Invocation, newSQL string) {
	if invocation.Properties == nil {
		invocation.Properties = make(map[string]interface{})
	}
	invocation.Properties["sql"] = newSQL
}