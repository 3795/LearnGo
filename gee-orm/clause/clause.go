package clause

import "strings"

type Type int

// 拼接SQL语句
type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

// 赋值SQL的各项变量
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// 组合SQL属性，拼接成一条完整SQL
// 此处orders（sql命令）是一个数组，例如：select * from where a = x，这种情况
func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok { // 将多个命令组合成一条SQL语句
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
