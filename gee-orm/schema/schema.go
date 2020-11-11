package schema

import (
	"LearnGo/gee-orm/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string // 字段名称
	Type string // 字段类型
	Tag  string // 额外的约束条件，非空，主键等
}

// 被映射对象Model，表名Name和字段Fields
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// 将任意对象解析为Schema实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(), // 将结构体的类名作为表名
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{Name: p.Name, Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type)))}

			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}

	return schema
}
