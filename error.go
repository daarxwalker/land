package land

type Error struct {
	Error   error
	Query   string
	Message string
}
