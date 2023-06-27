package land

type Error struct {
	error   error
	query   string
	message string
}
