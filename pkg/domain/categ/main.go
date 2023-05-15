package categ

type Repository interface {
	ListForUser(user string) ([]Entry, error)
	Save(user string, entries []Entry) error
}
