package land

type ColOpts struct {
	Default any
	Limit   int
	PK      bool
	NotNull bool
	Unique  bool
	Exclude bool
}

type column struct {
	entity   *entity
	name     string
	dataType string
	alias    string
	options  ColOpts
}

func createColumn(name, dataType string, options ColOpts) *column {
	c := &column{
		name:     name,
		dataType: dataType,
		options:  options,
	}
	return c
}
