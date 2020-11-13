package session

import (
	"LearnGo/gee-orm/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).refTable
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values)) // 结果数据切片
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	// 查询结果
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	// 将结果转换成对应的对象
	for rows.Next() {
		dest := reflect.New(destType).Elem() // 根据对象类型new一个对象
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface()) // 得到对象中对应字段的内存地址
		}
		if err := rows.Scan(values...); err != nil { // 将字段值填充到对象的属性中，此处处理比较简单，所以应该是要严格保存顺序一致
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest)) // 存储查询好的对象结果，改变指针指向的对象内容，但是并不改变指针
	}

	return rows.Close()
}
