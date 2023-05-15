package categ

import (
	"github.com/mehix/vf-sec-ctrls/pkg/domain/categ"
)

type Service struct {
	repo categ.Repository
}

func NewService(r categ.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ForUser(user string) (*Hierarchy, error) {
	entries, err := s.repo.ListForUser(user)
	if err != nil {
		return nil, err
	}

	h := NewFromList(entries)

	return h, err
}

func (s *Service) Save(user string, h *Hierarchy) error {
	if h == nil {
		return nil
	}
	return s.repo.Save(user, h.Entries())
}
