package memory

import (
	"sync"

	"github.com/mehix/vf-sec-ctrls/categ"
)

type repository struct {
	data map[string][]categ.Entry
	m    sync.RWMutex
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string][]categ.Entry),
	}
}

func (r *repository) ListForUser(user string) ([]categ.Entry, error) {
	r.m.RLock()
	entries := r.data[user]
	r.m.RUnlock()

	return entries, nil
}

func (r *repository) Save(user string, entries []categ.Entry) error {
	r.m.Lock()
	defer r.m.Unlock()
	r.data[user] = entries
	return nil
}
