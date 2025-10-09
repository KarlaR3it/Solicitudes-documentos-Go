package solicitud

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock del cliente de documentos para este archivo
type mockDocumentoClient struct {
	mock.Mock
}

func (m *mockDocumentoClient) GetBySolicitudID(solicitudID uint) ([]Documento, error) {
	args := m.Called(solicitudID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Documento), args.Error(1)
}

func TestService_GetAll(t *testing.T) {
	ctx := context.Background()
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	t.Run("debe retornar error cuando falla el repositorio", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		expectedError := errors.New("error en el repositorio")
		repo.On("GetAll", ctx, mock.AnythingOfType("solicitud.GetAllReq")).
			Return(nil, expectedError)

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.GetAll(ctx, GetAllReq{})

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("debe retornar solicitudes exitosamente", func(t *testing.T) {
		// Arrange
		testSolicitud := Solicitud{
			ID:     1,
			Titulo: "Test Solicitud",
			Estado: "pendiente",
		}

		documentos := []Documento{
			{ID: 1, NombreArchivo: "test.pdf", Extension: "pdf"},
		}

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("GetAll", ctx, mock.AnythingOfType("solicitud.GetAllReq")).
			Return([]Solicitud{testSolicitud}, nil)
		docClient.On("GetBySolicitudID", uint(1)).Return(documentos, nil)

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.GetAll(ctx, GetAllReq{})

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, "Test Solicitud", result[0].Titulo)
		assert.Len(t, result[0].Documentos, 1)
		assert.Equal(t, "test.pdf", result[0].Documentos[0].NombreArchivo)
		repo.AssertExpectations(t)
		docClient.AssertExpectations(t)
	})
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	t.Run("debe crear solicitud exitosamente", func(t *testing.T) {
		// Arrange
		req := CreateReq{
			Titulo:                   "Nueva Solicitud",
			Estado:                   "pendiente",
			Area:                     "TI",
			Pais:                     "Chile",
			Localizacion:             "Santiago",
			NumeroVacantes:           1,
			Descripcion:              "Descripción test",
			BaseEducacional:          "Universitaria",
			ConocimientosExcluyentes: "Ninguno",
			RentaDesde:               800000,
			RentaHasta:               1200000,
			ModalidadTrabajo:         "remoto",
			TipoServicio:             "freelance",
			NivelExperiencia:         "junior",
			FechaInicioProyecto:      "2024-01-01",
		}

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("Create", ctx, mock.AnythingOfType("*solicitud.Solicitud")).
			Return(nil).
			Run(func(args mock.Arguments) {
				s := args.Get(1).(*Solicitud)
				s.ID = 1 // Simular asignación de ID
			})

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.Create(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Nueva Solicitud", result.Titulo)
		repo.AssertExpectations(t)
	})

	t.Run("debe retornar error cuando falla la creación", func(t *testing.T) {
		// Arrange
		req := CreateReq{
			Titulo:                   "Solicitud Fallida",
			Estado:                   "pendiente",
			Area:                     "TI",
			Pais:                     "Chile",
			Localizacion:             "Santiago",
			NumeroVacantes:           1,
			Descripcion:              "Descripción test",
			BaseEducacional:          "Universitaria",
			ConocimientosExcluyentes: "Ninguno",
			RentaDesde:               800000,
			RentaHasta:               1200000,
			ModalidadTrabajo:         "remoto",
			TipoServicio:             "freelance",
			NivelExperiencia:         "junior",
			FechaInicioProyecto:      "2024-01-01",
		}

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		expectedError := errors.New("error de base de datos")
		repo.On("Create", ctx, mock.AnythingOfType("*solicitud.Solicitud")).Return(expectedError)

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.Create(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "error de base de datos")
		repo.AssertExpectations(t)
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	t.Run("debe retornar solicitud sin documentos exitosamente", func(t *testing.T) {
		// Arrange
		testSolicitud := &Solicitud{
			ID:     1,
			Titulo: "Test Solicitud",
			Estado: "pendiente",
		}

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("GetByID", ctx, uint(1)).Return(testSolicitud, nil)

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.GetByID(ctx, 1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Test Solicitud", result.Titulo)
		assert.Len(t, result.Documentos, 0) // GetByID no incluye documentos
		repo.AssertExpectations(t)
	})

	t.Run("debe retornar error cuando no encuentra solicitud", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		expectedError := errors.New("solicitud no encontrada")
		repo.On("GetByID", ctx, uint(999)).Return(nil, expectedError)

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.GetByID(ctx, 999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestService_GetByIDWithDocuments(t *testing.T) {
	ctx := context.Background()
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	t.Run("debe retornar solicitud con documentos exitosamente", func(t *testing.T) {
		// Arrange
		testSolicitud := &Solicitud{
			ID:     1,
			Titulo: "Test Solicitud Con Docs",
			Estado: "pendiente",
		}

		documentos := []Documento{
			{ID: 1, NombreArchivo: "doc1.pdf", Extension: "pdf"},
			{ID: 2, NombreArchivo: "doc2.jpg", Extension: "jpg"},
		}

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("GetByID", ctx, uint(1)).Return(testSolicitud, nil)
		docClient.On("GetBySolicitudID", uint(1)).Return(documentos, nil)

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.GetByIDWithDocuments(ctx, 1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Test Solicitud Con Docs", result.Titulo)
		assert.Len(t, result.Documentos, 2)
		assert.Equal(t, "doc1.pdf", result.Documentos[0].NombreArchivo)
		assert.Equal(t, "doc2.jpg", result.Documentos[1].NombreArchivo)
		repo.AssertExpectations(t)
		docClient.AssertExpectations(t)
	})

	t.Run("debe retornar solicitud sin documentos cuando cliente falla", func(t *testing.T) {
		// Arrange
		testSolicitud := &Solicitud{
			ID:     1,
			Titulo: "Test Solicitud",
			Estado: "pendiente",
		}

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("GetByID", ctx, uint(1)).Return(testSolicitud, nil)
		docClient.On("GetBySolicitudID", uint(1)).Return(nil, errors.New("error cliente documentos"))

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.GetByIDWithDocuments(ctx, 1)

		// Assert
		assert.NoError(t, err) // El servicio continúa funcionando
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Empty(t, result.Documentos) // Sin documentos por el error
		repo.AssertExpectations(t)
		docClient.AssertExpectations(t)
	})

	t.Run("debe retornar error cuando no encuentra solicitud", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		expectedError := errors.New("solicitud no encontrada")
		repo.On("GetByID", ctx, uint(999)).Return(nil, expectedError)

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.GetByIDWithDocuments(ctx, 999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		repo.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	t.Run("debe actualizar solicitud exitosamente", func(t *testing.T) {
		// Arrange
		updateReq := UpdateReq{
			Titulo: stringPtr("Solicitud Actualizada"),
			Estado: stringPtr("aprobada"),
			Area:   stringPtr("Marketing"),
		}

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("Update", ctx, uint(1), updateReq).Return(nil)

		service := NewService(repo, logger, docClient)

		// Act
		err := service.Update(ctx, 1, updateReq)

		// Assert
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("debe retornar error cuando falla la actualización", func(t *testing.T) {
		// Arrange
		updateReq := UpdateReq{
			Titulo: stringPtr("Solicitud Fallida"),
		}

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		expectedError := errors.New("error en actualización")
		repo.On("Update", ctx, uint(999), updateReq).Return(expectedError)

		service := NewService(repo, logger, docClient)

		// Act
		err := service.Update(ctx, 999, updateReq)

		// Assert
		assert.Error(t, err)
		assert.EqualError(t, err, "error en actualización")
		repo.AssertExpectations(t)
	})

	t.Run("debe actualizar con campos parciales", func(t *testing.T) {
		// Arrange - Solo actualizar título
		updateReq := UpdateReq{
			Titulo: stringPtr("Solo Título Actualizado"),
		}

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("Update", ctx, uint(1), updateReq).Return(nil)

		service := NewService(repo, logger, docClient)

		// Act
		err := service.Update(ctx, 1, updateReq)

		// Assert
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	logger := log.New(os.Stdout, "TEST: ", log.LstdFlags)

	t.Run("debe eliminar solicitud exitosamente", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("Delete", ctx, uint(1)).Return(nil)

		service := NewService(repo, logger, docClient)

		// Act
		err := service.Delete(ctx, 1)

		// Assert
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("debe retornar error cuando no encuentra solicitud para eliminar", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		expectedError := errors.New("solicitud no encontrada")
		repo.On("Delete", ctx, uint(999)).Return(expectedError)

		service := NewService(repo, logger, docClient)

		// Act
		err := service.Delete(ctx, 999)

		// Assert
		assert.Error(t, err)
		assert.EqualError(t, err, "solicitud no encontrada")
		repo.AssertExpectations(t)
	})

	t.Run("debe retornar error de base de datos", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		expectedError := errors.New("error de conexión a base de datos")
		repo.On("Delete", ctx, uint(1)).Return(expectedError)

		service := NewService(repo, logger, docClient)

		// Act
		err := service.Delete(ctx, 1)

		// Assert
		assert.Error(t, err)
		assert.EqualError(t, err, "error de conexión a base de datos")
		repo.AssertExpectations(t)
	})
}

// Función auxiliar para crear punteros a strings
func stringPtr(s string) *string {
	return &s
}
