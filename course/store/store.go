package store

type PostgresStore struct {
	users []string
}

func (p PostgresStore) GetUsers() []string {
	return p.users
}
