package land

type Migration interface {
	Up(func(l Land)) Migration
	Down(func(l Land)) Migration
}

type migration struct {
	id   string
	up   func(l Land)
	down func(l Land)
}

func (m *migration) Up(fn func(l Land)) Migration {
	m.up = fn
	return m
}

func (m *migration) Down(fn func(l Land)) Migration {
	m.down = fn
	return m
}
