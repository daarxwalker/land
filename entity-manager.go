package land

type EntityManager interface {
	CreateEntity(name string) Entity
}

type entityManager struct {
	land     *land
	entities []*entity
}

func createEntityManager(land *land) *entityManager {
	return &entityManager{
		land:     land,
		entities: make([]*entity, 0),
	}
}

func (m *entityManager) CreateEntity(name string) Entity {
	e := createEntity(m, name)
	m.entities = append(m.entities, e)
	return e
}
