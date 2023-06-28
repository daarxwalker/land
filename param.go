package land

type Param struct {
	Id       int          `json:"id"`
	Ids      int          `json:"ids"`
	Fulltext string       `json:"fulltext"`
	Offset   int          `json:"offset"`
	Limit    int          `json:"limit"`
	Order    []OrderParam `json:"order"`
	All      bool         `json:"all"`
}

type OrderParam struct {
	Key       string `json:"key"`
	Direction string `json:"direction"`
	Dynamic   bool   `json:"-"`
}
