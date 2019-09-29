package builder

type OrderPair struct {
	Key   string
	Value string
}

type WherePair struct {
	Query string
	Args  []interface{}
}

// 用与逻辑合并两个查询条件
func (w *WherePair) And(where *WherePair) *WherePair {
	if w.Query == "" {
		return where
	} else {
		return &WherePair{Query: w.Query + "AND" + where.Query, Args: append(w.Args, where.Args...)}
	}
}

// 用或逻辑合并两个查询条件
func (w *WherePair) Or(where *WherePair) *WherePair {
	if w.Query == "" {
		return where
	} else {
		return &WherePair{Query: w.Query + "OR" + where.Query, Args: append(w.Args, where.Args...)}
	}
}
