package land

type CoalesceQuery interface {
	Column(column ...string) CoalesceQuery
	Value(value ...any) CoalesceQuery
	Entity(entity Entity) CoalesceQuery
	Use(use bool) CoalesceQuery
	getPtr() *coalesceQueryBuilder
}

type coalesceQueryBuilder struct {
	entity  *entity
	use     bool
	columns []string
	values  []any
}

func createCoalesce() *coalesceQueryBuilder {
	return &coalesceQueryBuilder{
		use:     true,
		columns: make([]string, 0),
		values:  make([]any, 0),
	}
}

func (q *coalesceQueryBuilder) getPtr() *coalesceQueryBuilder {
	return q
}

func (q *coalesceQueryBuilder) Column(column ...string) CoalesceQuery {
	q.columns = column
	return q
}

func (q *coalesceQueryBuilder) Value(value ...any) CoalesceQuery {
	q.values = value
	return q
}

func (q *coalesceQueryBuilder) Entity(entity Entity) CoalesceQuery {
	q.entity = entity.getPtr()
	return q
}

func (q *coalesceQueryBuilder) Use(use bool) CoalesceQuery {
	q.use = use
	return q
}
