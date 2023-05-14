package categ

import (
	old "github.com/mehix/vf-sec-ctrls/categ"
	"github.com/mehix/vf-sec-ctrls/pkg/domain/categ"
)

type Service struct {
	repo categ.Repository
}

func NewService(r categ.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) ForUser(user string) (*old.Hierarchy, error) {
	entries, err := s.repo.ListForUser(user)
	if err != nil {
		return nil, err
	}

	h := old.NewFromList(entries)

	return h, err
}

func (s *Service) Save(user string, h *old.Hierarchy) error {
	if h == nil {
		return nil
	}
	return s.repo.Save(user, h.Entries())
}
