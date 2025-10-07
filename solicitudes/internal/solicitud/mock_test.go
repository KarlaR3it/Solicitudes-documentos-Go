package solicitud_test

import "github.com/kramirez/solicitudes/internal/solicitud"

type mockRepository struct {
	CreateMock  func(ctx context.Context, solicitud *Solicitud) error
	GetAllMock  func(ctx context.Context, filters GetAllReq) ([]Solicitud, error)
	GetByIDMock func(ctx context.Context, id uint) (*Solicitud, error)
	UpdateMock  func(ctx context.Context, id uint, req UpdateReq) error
	DeleteMock  func(ctx context.Context, id uint) error
}

func (m *mockRepository) Create(ctx context.Context, solicitud *Solicitud) error {
	return m.CreateMock(ctx, solicitud)
}

func (m *mockRepository) GetAll(ctx context.Context, filters GetAllReq) ([]Solicitud, error) {
	return m.GetAllMock(ctx, filters)
}

func (m *mockRepository) GetByID(ctx context.Context, id uint) (*Solicitud, error) {
	return m.GetByIDMock(ctx, id)
}

func (m *mockRepository) Update(ctx context.Context, id uint, req UpdateReq) error {
	return m.UpdateMock(ctx, id, req)
}

func (m *mockRepository) Delete(ctx context.Context, id uint) error {
	return m.DeleteMock(ctx, id)
}
