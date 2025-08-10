package config

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"gobatis/logger"
	"io/ioutil"
	"reflect"
	"strings"
)

// Configuration 框架配置
type Configuration struct {
	DataSource   *DataSource
	MapperConfig *MapperConfig
	Plugins      []Plugin
	Logger       logger.Interface
}

// DataSource 数据源配置
type DataSource struct {
	DriverName     string
	DataSourceName string
	DB             *sql.DB
}

// MapperConfig Mapper 配置
type MapperConfig struct {
	Mappers map[string]*MapperStatement
}

// MapperStatement SQL 语句配置
type MapperStatement struct {
	ID            string
	SQL           string
	ResultType    reflect.Type
	StatementType StatementType
}

// StatementType SQL 语句类型
type StatementType int

const (
	SELECT StatementType = iota
	INSERT
	UPDATE
	DELETE
)

// Plugin 插件接口
type Plugin interface {
	Intercept(invocation *Invocation) (interface{}, error)
	SetProperties(properties map[string]string)
}

// Invocation 拦截调用信息
type Invocation struct {
	Target interface{}
	Method reflect.Method
	Args   []interface{}
}

// NewConfiguration 创建新的配置
func NewConfiguration() *Configuration {
	return &Configuration{
		MapperConfig: &MapperConfig{
			Mappers: make(map[string]*MapperStatement),
		},
		Plugins: make([]Plugin, 0),
		Logger:  logger.Default,
	}
}

// SetDataSource 设置数据源
func (c *Configuration) SetDataSource(driverName, dataSourceName string) error {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	c.DataSource = &DataSource{
		DriverName:     driverName,
		DataSourceName: dataSourceName,
		DB:             db,
	}

	return nil
}

// AddMapperXML 添加 Mapper XML 配置
func (c *Configuration) AddMapperXML(xmlPath string) error {
	data, err := ioutil.ReadFile(xmlPath)
	if err != nil {
		return fmt.Errorf("failed to read mapper xml: %w", err)
	}

	var mapper XMLMapper
	if err := xml.Unmarshal(data, &mapper); err != nil {
		return fmt.Errorf("failed to parse mapper xml: %w", err)
	}

	// 解析 select 语句
	for _, sel := range mapper.Selects {
		statementId := mapper.Namespace + "." + sel.ID
		c.MapperConfig.Mappers[statementId] = &MapperStatement{
			ID:            statementId,
			SQL:           strings.TrimSpace(sel.SQL),
			StatementType: SELECT,
		}
	}

	// 解析 insert 语句
	for _, ins := range mapper.Inserts {
		statementId := mapper.Namespace + "." + ins.ID
		c.MapperConfig.Mappers[statementId] = &MapperStatement{
			ID:            statementId,
			SQL:           strings.TrimSpace(ins.SQL),
			StatementType: INSERT,
		}
	}

	// 解析 update 语句
	for _, upd := range mapper.Updates {
		statementId := mapper.Namespace + "." + upd.ID
		c.MapperConfig.Mappers[statementId] = &MapperStatement{
			ID:            statementId,
			SQL:           strings.TrimSpace(upd.SQL),
			StatementType: UPDATE,
		}
	}

	// 解析 delete 语句
	for _, del := range mapper.Deletes {
		statementId := mapper.Namespace + "." + del.ID
		c.MapperConfig.Mappers[statementId] = &MapperStatement{
			ID:            statementId,
			SQL:           strings.TrimSpace(del.SQL),
			StatementType: DELETE,
		}
	}

	return nil
}

// AddPlugin 添加插件
func (c *Configuration) AddPlugin(plugin Plugin) {
	c.Plugins = append(c.Plugins, plugin)
}

// GetMapperStatement 获取 Mapper 语句
func (c *Configuration) GetMapperStatement(statementId string) (*MapperStatement, bool) {
	stmt, exists := c.MapperConfig.Mappers[statementId]
	return stmt, exists
}

// XMLMapper XML Mapper 结构
type XMLMapper struct {
	XMLName   xml.Name    `xml:"mapper"`
	Namespace string      `xml:"namespace,attr"`
	Selects   []XMLSelect `xml:"select"`
	Inserts   []XMLInsert `xml:"insert"`
	Updates   []XMLUpdate `xml:"update"`
	Deletes   []XMLDelete `xml:"delete"`
}

// XMLSelect XML Select 语句
type XMLSelect struct {
	ID         string `xml:"id,attr"`
	ResultType string `xml:"resultType,attr"`
	SQL        string `xml:",chardata"`
}

// XMLInsert XML Insert 语句
type XMLInsert struct {
	ID  string `xml:"id,attr"`
	SQL string `xml:",chardata"`
}

// XMLUpdate XML Update 语句
type XMLUpdate struct {
	ID  string `xml:"id,attr"`
	SQL string `xml:",chardata"`
}

// XMLDelete XML Delete 语句
type XMLDelete struct {
	ID  string `xml:"id,attr"`
	SQL string `xml:",chardata"`
}
