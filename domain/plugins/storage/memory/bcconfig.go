package memory

import "github.com/d-velop/dvelop-app-template-go/domain"

type bcConfigStore struct {
	config domain.BcConfig
}

func NewBcConfigStore() domain.BcConfigRepository {
	return &bcConfigStore{}
}

func (s *bcConfigStore) GetBcConfig() domain.BcConfig {
	return s.config
}

func (s *bcConfigStore) SaveBcConfig(config domain.BcConfig) {
	s.config = config
}
