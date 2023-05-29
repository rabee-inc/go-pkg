package cloudfirestore

func NewTypedBatchGetter[T any](
	bg BatchGetter,
	getDoc FuncGetDoc,
	getID FuncGetID[T],
	convert FuncConvert[T, T],
) TypedBatchGetter[T] {
	if convert == nil {
		convert = func(t *T) *T { return t }
	}
	cbg := &convertibleBatchGetter[T, T]{
		bg:      bg,
		dstMap:  map[string]*T{},
		items:   []*convertibleBatchGetterItem[*T]{},
		getDoc:  getDoc,
		getID:   getID,
		convert: convert,
	}

	bg.OnCommit(func() {
		cbg.convertAll()
	})

	bg.OnEnd(func() {
		cbg.convertAll()
	})

	return cbg
}

type FuncGetModel[E, M any] func(e *E) (id string, m *M)

func GetModelMapByBatchGetter[E, M any](
	bg TypedBatchGetter[E],
	getModel FuncGetModel[E, M],
) map[string]*M {
	es := bg.GetMap()
	ms := map[string]*M{}
	for _, e := range es {
		id, m := getModel(e)
		if m != nil {
			ms[id] = m
		}
	}
	return ms
}
