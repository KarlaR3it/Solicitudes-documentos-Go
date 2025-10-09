package solicitud

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEndpoint_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange: mock repo and doc client, real service
	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.POST("/solicitudes", ep.Create)

	// Expect repo.Create to be called and set ID
	repo.On("Create", mock.Anything, mock.AnythingOfType("*solicitud.Solicitud")).
		Return(nil).
		Run(func(args mock.Arguments) {
			s := args.Get(1).(*Solicitud)
			s.ID = 1
		})

	// Valid request payload
	payload := CreateReq{
		Titulo:                   "DevOps Engineer",
		Estado:                   "pendiente",
		Area:                     "Infraestructura",
		Pais:                     "Chile",
		Localizacion:             "Concepción",
		NumeroVacantes:           1,
		Descripcion:              "Ingeniero DevOps",
		BaseEducacional:          "Ingeniería",
		ConocimientosExcluyentes: "AWS",
		RentaDesde:               1500000,
		RentaHasta:               2200000,
		ModalidadTrabajo:         "presencial",
		TipoServicio:             "infraestructura",
		NivelExperiencia:         "senior",
		FechaInicioProyecto:      "2025-12-01",
		UsuarioID:                uintPtr(3),
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/solicitudes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var resp Solicitud
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, payload.Titulo, resp.Titulo)
	assert.Equal(t, payload.Estado, resp.Estado)

	repo.AssertExpectations(t)
}

func TestEndpoint_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.POST("/solicitudes", ep.Create)

	// invalid JSON body
	req := httptest.NewRequest(http.MethodPost, "/solicitudes", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	// Repo should not be called
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestEndpoint_Create_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.POST("/solicitudes", ep.Create)

	// Configure repo to return error on Create
	repo.On("Create", mock.Anything, mock.AnythingOfType("*solicitud.Solicitud")).
		Return(assert.AnError)

	payload := CreateReq{
		Titulo:                   "Role",
		Estado:                   "pendiente",
		Area:                     "IT",
		Pais:                     "CL",
		Localizacion:             "SCL",
		NumeroVacantes:           1,
		Descripcion:              "desc",
		BaseEducacional:          "base",
		ConocimientosExcluyentes: "aws",
		RentaDesde:               1,
		RentaHasta:               2,
		ModalidadTrabajo:         "remote",
		TipoServicio:             "dev",
		NivelExperiencia:         "jr",
		FechaInicioProyecto:      "2025-12-01",
		UsuarioID:                uintPtr(3),
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/solicitudes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}

func TestEndpoint_GetByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.GET("/solicitudes/:id", ep.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/solicitudes/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
}

func TestEndpoint_GetByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.GET("/solicitudes/:id", ep.GetByID)

	// Service.GetByID calls repo.GetByID and propagates error; endpoint maps to 404
	repo.On("GetByID", mock.Anything, uint(99)).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/solicitudes/99", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	repo.AssertExpectations(t)
}

func TestEndpoint_GetByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.GET("/solicitudes/:id", ep.GetByID)

	s := &Solicitud{ID: 7, Titulo: "Algo", Estado: "pendiente"}
	repo.On("GetByID", mock.Anything, uint(7)).Return(s, nil)

	req := httptest.NewRequest(http.MethodGet, "/solicitudes/7", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp SolicitudResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, uint(7), resp.ID)
	assert.Equal(t, "Algo", resp.Titulo)
	repo.AssertExpectations(t)
}

func TestEndpoint_Update_ForbiddenFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.PATCH("/solicitudes/:id", ep.Update)

	// Body tries to update forbidden field usuario_id
	body := bytes.NewBufferString(`{"usuario_id": 10}`)
	req := httptest.NewRequest(http.MethodPatch, "/solicitudes/1", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}

func TestEndpoint_Update_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.PATCH("/solicitudes/:id", ep.Update)

	// Service.Update calls repo.GetByID first and returns error mapped to 500
	repo.On("GetByID", mock.Anything, uint(42)).Return(nil, assert.AnError)

	body := bytes.NewBufferString(`{"titulo":"nuevo"}`)
	req := httptest.NewRequest(http.MethodPatch, "/solicitudes/42", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}

func TestEndpoint_Update_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.PATCH("/solicitudes/:id", ep.Update)

	// Existence check then update ok
	repo.On("GetByID", mock.Anything, uint(5)).Return(&Solicitud{ID: 5}, nil)
	repo.On("Update", mock.Anything, uint(5), mock.AnythingOfType("solicitud.UpdateReq")).Return(nil)

	body := bytes.NewBufferString(`{"titulo":"nuevo"}`)
	req := httptest.NewRequest(http.MethodPatch, "/solicitudes/5", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestEndpoint_Delete_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.DELETE("/solicitudes/:id", ep.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/solicitudes/xyz", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestEndpoint_Delete_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.DELETE("/solicitudes/:id", ep.Delete)

	// Service.Delete checks existence then calls repo.Delete; simulate existence ok then delete error
	repo.On("GetByID", mock.Anything, uint(10)).Return(&Solicitud{ID: 10}, nil)
	repo.On("Delete", mock.Anything, uint(10)).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/solicitudes/10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	repo.AssertExpectations(t)
}

func TestEndpoint_Delete_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.DELETE("/solicitudes/:id", ep.Delete)

	repo.On("GetByID", mock.Anything, uint(3)).Return(&Solicitud{ID: 3}, nil)
	repo.On("Delete", mock.Anything, uint(3)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/solicitudes/3", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	repo.AssertExpectations(t)
}

func TestEndpoint_Update_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.PATCH("/solicitudes/:id", ep.Update)

	req := httptest.NewRequest(http.MethodPatch, "/solicitudes/1", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}

func TestEndpoint_Update_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.PATCH("/solicitudes/:id", ep.Update)

	req := httptest.NewRequest(http.MethodPatch, "/solicitudes/not-a-number", bytes.NewBufferString(`{"titulo":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
}

func TestEndpoint_GetAll_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.GET("/solicitudes", ep.GetAll)

	// repo returns one Solicitud, doc client returns no documents
	repo.On("GetAll", mock.Anything, mock.AnythingOfType("solicitud.GetAllReq")).
		Return([]Solicitud{{ID: 100, Titulo: "T", Estado: "pendiente"}}, nil)
	docClient.On("GetBySolicitudID", uint(100)).Return([]Documento{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/solicitudes", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []SolicitudResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, uint(100), resp[0].ID)
	repo.AssertExpectations(t)
	docClient.AssertExpectations(t)
}

func TestEndpoint_GetByIDWithDocuments_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := new(mockRepository)
	docClient := new(mockDocumentoClient)
	logger := log.New(io.Discard, "", 0)
	svc := NewService(repo, logger, docClient)
	ep := NewEndpoint(svc)

	r := gin.New()
	r.GET("/solicitudes/:id/con-documentos", ep.GetByIDWithDocuments)

	repo.On("GetByID", mock.Anything, uint(55)).Return(&Solicitud{ID: 55, Titulo: "ConDocs", Estado: "pendiente"}, nil)
	docClient.On("GetBySolicitudID", uint(55)).Return([]Documento{{ID: 1, NombreArchivo: "a.pdf", Extension: "pdf"}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/solicitudes/55/con-documentos", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp SolicitudResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, uint(55), resp.ID)
	assert.Len(t, resp.Documentos, 1)
	assert.Equal(t, "a.pdf", resp.Documentos[0].NombreArchivo)
	repo.AssertExpectations(t)
	docClient.AssertExpectations(t)
}
