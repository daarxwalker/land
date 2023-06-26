package land

type Param struct {
	Id       int     `json:"id"`
	Ids      int     `json:"ids"`
	Fulltext string  `json:"fulltext"`
	Offset   int     `json:"offset"`
	Limit    int     `json:"limit"`
	Order    []Order `json:"order"`
	All      bool    `json:"all"`
}

type Order struct {
	Key       string `json:"key"`
	Direction string `json:"direction"`

	dynamic bool `json:"-"`
}
