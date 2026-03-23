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

func (m *MockContentTypeRepo) DeleteContentTypeBySlug(slug string, opt repository.DeleteContentTypeOptions) error {
	args := m.Called(slug, opt)
	return args.Error(0)
}

func (m *MockContentTypeRepo) DeleteContentTypeById(id uuid.UUID, opt repository.DeleteContentTypeOptions) error {
	args := m.Called(id, opt)
	return args.Error(0)
}

func (m *MockContentTypeRepo) GetContentTypeById(id uuid.UUID, opt repository.GetContentTypeOptions) (*domain.ContentType, error) {
	args := m.Called(id, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ContentType), args.Error(1)
}

func (m *MockContentTypeRepo) GetContentTypesByAccountId(account uuid.UUID, opt repository.GetContentTypeOptions) ([]*domain.ContentType, error) {
	args := m.Called(account, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ContentType), args.Error(1)
}

type MockEntryRepo struct {
	mock.Mock
}

var _ repository.EntryRepository = &MockEntryRepo{}

func (m *MockEntryRepo) GetEntryById(id uuid.UUID, opt repository.GetEntryOptions) (*domain.Entry, error) {
	args := m.Called(id, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Entry), args.Error(1)
}

func (m *MockEntryRepo) GetEntryByContentTypeAndSlug(ctId uuid.UUID, slug string, opt repository.GetEntryOptions) (*domain.Entry, error) {
	args := m.Called(ctId, slug, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Entry), args.Error(1)
}

func (m *MockEntryRepo) GetEntriesByContentType(ctId uuid.UUID, opt repository.GetEntryOptions) ([]*domain.Entry, error) {
	args := m.Called(ctId, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Entry), args.Error(1)
}

func (m *MockEntryRepo) GetEntriesByFilter(e *domain.Entry, opt repository.GetEntryOptions) ([]*domain.Entry, error) {
	args := m.Called(e, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Entry), args.Error(1)
}

func (m *MockEntryRepo) AddEntry(e *domain.Entry, opt repository.AddEntryOptions) error {
	args := m.Called(e, opt)
	return args.Error(0)
}

func (m *MockEntryRepo) UpdateEntry(e *domain.Entry, opt repository.UpdateEntryOptions) error {
	args := m.Called(e, opt)
	return args.Error(0)
}
