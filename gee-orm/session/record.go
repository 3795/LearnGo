package session

import (
	"LearnGo/gee-orm/clause"
	"errors"
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

func (s *Session) Update(kv ...interface{}) (int64, error) {
	// 这个地方做两个判断，参数要么可以传map键值对，要么可以直接传update中的键值语句，然后在函数中进行处理
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable(), m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
