package land

type testModel struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Active   bool   `json:"active"`
}

type testEntityContainer struct {
	test func(em EntityManager) Entity
}

const (
	testEntityName  = "tests"
	testEntityAlias = "t"
	testName        = "name"
	testLastname    = "lastname"
	testActive      = "active"
)

var (
	testEntityColumns = []string{Id, testName, testLastname, testActive}
)

func testCreatePostgresInstance() Land {
	return New(Config{
		Production:   false,
		Log:          true,
		Timezone:     false,
		DatabaseType: Postgres,
	}, nil)
}

func testEntity(em EntityManager) Entity {
	return em.CreateEntity(testEntityName).
		SetAlias(testEntityAlias).
		SetColumn(testName, Varchar, ColOpts{Limit: 255, NotNull: true}).
		SetColumn(testLastname, Varchar, ColOpts{Limit: 255, NotNull: true}).
		SetColumn(testActive, Boolean, ColOpts{NotNull: true, Default: false}).
		SetFulltext().
		SetCreatedAt().
		SetUpdatedAt()
}

func testSecondEntity(em EntityManager) Entity {
	return em.CreateEntity(testEntityName).
		SetAlias(testEntityAlias+"2").
		SetColumn(testName, Varchar, ColOpts{Limit: 255, NotNull: true}).
		SetColumn(testLastname, Varchar, ColOpts{Limit: 255, NotNull: true}).
		SetFulltext().
		SetCreatedAt().
		SetUpdatedAt()
}
