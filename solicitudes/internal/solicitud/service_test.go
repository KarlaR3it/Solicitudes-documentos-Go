package solicitud_test

import (
	"context"
	"testing"

	"github.com/kramirez/solicitudes/internal/solicitud"
)

func TestService_GetAll(t *testing.T) {

	t.Run("should return an error", func(t *testing.T) {
		repo := &mockRepository{
			GetAllMock: func(ctx context.Context, filters Filters, offset, limit int) ([]solicitud.Solicitud, error) {
				return []solicitud.Solicitud{}, nil
			},
		}
	})
}
