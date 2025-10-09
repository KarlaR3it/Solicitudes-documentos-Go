package solicitud

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, solicitud *Solicitud) error {
	args := m.Called(ctx, solicitud)
	return args.Error(0)
}

func (m *mockRepository) GetAll(ctx context.Context, filters GetAllReq) ([]Solicitud, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Solicitud), args.Error(1)
}

func (m *mockRepository) GetByID(ctx context.Context, id uint) (*Solicitud, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Solicitud), args.Error(1)
}

func (m *mockRepository) Update(ctx context.Context, id uint, req UpdateReq) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *mockRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
