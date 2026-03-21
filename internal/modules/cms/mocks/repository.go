package mocks

import (
	"github.com/google/uuid"
	"github.com/manaschubby/gocms/internal/modules/cms/domain"
	"github.com/manaschubby/gocms/internal/modules/cms/repository"
	"github.com/stretchr/testify/mock"
)

// Mocking the repositories
type MockAccountRepo struct{ mock.Mock }

var _ repository.AccountRepository = &MockAccountRepo{}

func (m *MockAccountRepo) GetAccountByUUID(id uuid.UUID, opt repository.GetAccountOptions) (*domain.Account, error) {
	args := m.Called(id, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepo) GetAccounts(opt repository.GetAccountsOptions) ([]*domain.Account, error) {
	args := m.Called(opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Account), args.Error(1)
}

func (m *MockAccountRepo) CreateAccount(account *domain.Account, opt repository.CreateAccountOptions) error {
	args := m.Called(account, opt)
	return args.Error(0)
}

type MockContentTypeRepo struct{ mock.Mock }

var _ repository.ContentTypeRepository = &MockContentTypeRepo{}

func (m *MockContentTypeRepo) GetContentTypeBySlug(slug string, opt repository.GetContentTypeOptions) (*domain.ContentType, error) {
	args := m.Called(slug, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ContentType), args.Error(1)
}
func (m *MockContentTypeRepo) CreateNewContentType(ct *domain.ContentType, opt repository.CreateNewContentTypeOptions) error {
	return m.Called(ct, opt).Error(0)
}
