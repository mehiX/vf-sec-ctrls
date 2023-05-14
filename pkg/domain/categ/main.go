package categ

import "github.com/mehix/vf-sec-ctrls/categ"

type Repository interface {
	ListForUser(user string) ([]categ.Entry, error)
	Save(user string, entries []categ.Entry) error
}
