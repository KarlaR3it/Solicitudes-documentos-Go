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

func (m *mockDocumentoClient) DeleteBySolicitudID(solicitudID uint) error {
	args := m.Called(solicitudID)
	return args.Error(0)
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

	validRequest := CreateReq{
		Titulo:                   "DevOps Engineer",
		Estado:                   "pendiente",
		Area:                     "Infraestructura",
		Pais:                     "Chile",
		Localizacion:             "Concepción",
		NumeroVacantes:           1,
		Descripcion:              "Ingeniero DevOps para automatizar procesos de CI/CD y gestionar infraestructura en la nube",
		BaseEducacional:          "Ingeniería en Informática o experiencia equivalente demostrable",
		ConocimientosExcluyentes: "AWS, Docker, Kubernetes, Jenkins, Terraform, Linux",
		RentaDesde:               1500000,
		RentaHasta:               2200000,
		ModalidadTrabajo:         "presencial",
		TipoServicio:             "infraestructura",
		NivelExperiencia:         "senior",
		FechaInicioProyecto:      "2025-12-01",
		UsuarioID:                uintPtr(3),
	}

	t.Run("debe crear solicitud exitosamente con todos los campos", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("Create", ctx, mock.AnythingOfType("*solicitud.Solicitud")).
			Return(nil).
			Run(func(args mock.Arguments) {
				s := args.Get(1).(*Solicitud)
				s.ID = 1
			})

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.Create(ctx, validRequest)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, validRequest.Titulo, result.Titulo)
		assert.Equal(t, validRequest.Estado, result.Estado)
		assert.Equal(t, validRequest.Area, result.Area)
		assert.Equal(t, validRequest.RentaDesde, result.RentaDesde)
		repo.AssertExpectations(t)
	})

	t.Run("debe retornar error cuando falla la creación en el repositorio", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		expectedError := errors.New("error de base de datos")
		repo.On("Create", ctx, mock.AnythingOfType("*solicitud.Solicitud")).Return(expectedError)

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.Create(ctx, validRequest)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "error de base de datos")
		repo.AssertExpectations(t)
	})

	t.Run("debe validar campos requeridos", func(t *testing.T) {
		tests := []struct {
			name       string
			req        CreateReq
			errMsg     string
			setupMocks func(*mockRepository, *mockDocumentoClient)
		}{
			{
				name: "falta título",
				req: func() CreateReq {
					r := validRequest
					r.Titulo = ""
					return r
				}(),
				errMsg: "el título es requerido",
				setupMocks: func(r *mockRepository, d *mockDocumentoClient) {
					// No se esperan llamadas al repositorio
				},
			},
			{
				name: "falta estado",
				req: func() CreateReq {
					r := validRequest
					r.Estado = ""
					return r
				}(),
				errMsg: "", // No debería fallar ya que tiene valor por defecto
				setupMocks: func(r *mockRepository, d *mockDocumentoClient) {
					r.On("Create", ctx, mock.AnythingOfType("*solicitud.Solicitud")).
						Return(nil).
						Run(func(args mock.Arguments) {
							s := args.Get(1).(*Solicitud)
							s.ID = 1
						})
				},
			},
			{
				name: "fecha inválida",
				req: func() CreateReq {
					r := validRequest
					r.FechaInicioProyecto = "fecha-invalida"
					return r
				}(),
				errMsg: "formato de fecha inválido, debe ser YYYY-MM-DD",
				setupMocks: func(r *mockRepository, d *mockDocumentoClient) {
					// No se esperan llamadas al repositorio
				},
			},
			{
				name: "rango de renta inválido",
				req: func() CreateReq {
					r := validRequest
					r.RentaDesde = 2000000
					r.RentaHasta = 1000000
					return r
				}(),
				errMsg: "el rango de renta es inválido",
				setupMocks: func(r *mockRepository, d *mockDocumentoClient) {
					// No se esperan llamadas al repositorio
				},
			},
			{
				name: "falta área",
				req: func() CreateReq {
					r := validRequest
					r.Area = ""
					return r
				}(),
				errMsg: "el área es requerida",
				setupMocks: func(r *mockRepository, d *mockDocumentoClient) {
					// No se esperan llamadas al repositorio
				},
			},
			{
				name: "falta país",
				req: func() CreateReq {
					r := validRequest
					r.Pais = ""
					return r
				}(),
				errMsg: "el país es requerido",
				setupMocks: func(r *mockRepository, d *mockDocumentoClient) {
					// No se esperan llamadas al repositorio
				},
			},
			{
				name: "falta localización",
				req: func() CreateReq {
					r := validRequest
					r.Localizacion = ""
					return r
				}(),
				errMsg: "la localización es requerida",
				setupMocks: func(r *mockRepository, d *mockDocumentoClient) {
					// No se esperan llamadas al repositorio
				},
			},
			{
				name: "falta ID de usuario",
				req: func() CreateReq {
					r := validRequest
					r.UsuarioID = nil
					return r
				}(),
				errMsg: "el ID de usuario es requerido",
				setupMocks: func(r *mockRepository, d *mockDocumentoClient) {
					// No se esperan llamadas al repositorio
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				repo := new(mockRepository)
				docClient := new(mockDocumentoClient)
				service := NewService(repo, logger, docClient)

				// Configurar mocks si es necesario
				if tt.setupMocks != nil {
					tt.setupMocks(repo, docClient)
				}

				// Act
				result, err := service.Create(ctx, tt.req)
				
				// Assert
				if tt.errMsg != "" {
					assert.Error(t, err)
					assert.Nil(t, result)
					assert.Contains(t, err.Error(), tt.errMsg)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, result)
				}

				// Verificar que no se llamó al repositorio a menos que se espere
				if tt.setupMocks == nil {
					repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
				} else {
					repo.AssertExpectations(t)
				}
			})
		}
	})

	t.Run("debe manejar valores por defecto correctamente", func(t *testing.T) {
		req := validRequest
		req.Estado = "" // Estado debería tener un valor por defecto

		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		repo.On("Create", ctx, mock.AnythingOfType("*solicitud.Solicitud")).
			Return(nil).
			Run(func(args mock.Arguments) {
				s := args.Get(1).(*Solicitud)
				s.ID = 1
				// Verificar que el estado se haya establecido correctamente
				assert.Equal(t, "pendiente", s.Estado, "El estado debería tener el valor por defecto 'pendiente'")
			})

		service := NewService(repo, logger, docClient)

		// Act
		result, err := service.Create(ctx, req)

		// Assert
		assert.NoError(t, err, "No debería haber error al crear la solicitud")
		assert.NotNil(t, result, "El resultado no debería ser nulo")
		assert.Equal(t, "pendiente", result.Estado, "El estado debería ser 'pendiente'")
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

		// Mock para verificar que existe la solicitud
		existingSolicitud := &Solicitud{ID: 1, Titulo: "Original"}
		repo.On("GetByID", ctx, uint(1)).Return(existingSolicitud, nil)
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

		// Mock para verificar que existe la solicitud
		existingSolicitud := &Solicitud{ID: 999, Titulo: "Original"}
		repo.On("GetByID", ctx, uint(999)).Return(existingSolicitud, nil)
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

		// Mock para verificar que existe la solicitud
		existingSolicitud := &Solicitud{ID: 1, Titulo: "Original"}
		repo.On("GetByID", ctx, uint(1)).Return(existingSolicitud, nil)
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

		// Mock para verificar que existe la solicitud
		existingSolicitud := &Solicitud{ID: 1, Titulo: "Original"}
		repo.On("GetByID", ctx, uint(1)).Return(existingSolicitud, nil)
		docClient.On("DeleteBySolicitudID", uint(1)).Return(nil)
		repo.On("Delete", ctx, uint(1)).Return(nil)

		service := NewService(repo, logger, docClient)

		// Act
		err := service.Delete(ctx, 1)

		// Assert
		assert.NoError(t, err)
		repo.AssertExpectations(t)
		docClient.AssertExpectations(t)
	})

	t.Run("debe retornar error cuando no encuentra solicitud para eliminar", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		expectedError := errors.New("solicitud no encontrada")
		repo.On("GetByID", ctx, uint(999)).Return(nil, expectedError)

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

		// Mock para verificar que existe la solicitud
		existingSolicitud := &Solicitud{ID: 1, Titulo: "Original"}
		repo.On("GetByID", ctx, uint(1)).Return(existingSolicitud, nil)
		docClient.On("DeleteBySolicitudID", uint(1)).Return(nil)

		expectedError := errors.New("error de base de datos")
		repo.On("Delete", ctx, uint(1)).Return(expectedError)

		service := NewService(repo, logger, docClient)

		// Act
		err := service.Delete(ctx, 1)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		repo.AssertExpectations(t)
		docClient.AssertExpectations(t)
	})

	t.Run("debe continuar eliminando solicitud aunque falle eliminación de documentos", func(t *testing.T) {
		// Arrange
		repo := new(mockRepository)
		docClient := new(mockDocumentoClient)

		// Mock para verificar que existe la solicitud
		existingSolicitud := &Solicitud{ID: 1, Titulo: "Original"}
		repo.On("GetByID", ctx, uint(1)).Return(existingSolicitud, nil)
		
		// Simular que falla la eliminación de documentos
		docClient.On("DeleteBySolicitudID", uint(1)).Return(errors.New("error al eliminar documentos"))
		
		// Pero la solicitud se elimina exitosamente
		repo.On("Delete", ctx, uint(1)).Return(nil)

		service := NewService(repo, logger, docClient)

		// Act
		err := service.Delete(ctx, 1)

		// Assert
		assert.NoError(t, err) // No debe retornar error aunque fallen los documentos
		repo.AssertExpectations(t)
		docClient.AssertExpectations(t)
	})
}

// Función auxiliar para crear punteros a strings
func stringPtr(s string) *string {
	return &s
}

// Función auxiliar para crear punteros a uint
func uintPtr(u uint) *uint {
	return &u
}
