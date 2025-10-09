package solicitud

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestRepository_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	solicitud := &Solicitud{
		Titulo: "Test CRUD",
		Estado: "activa",
		Area:   "IT",
	}

	err := repo.Create(context.Background(), solicitud)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "titulo", "estado"}).
		AddRow(1, "Test CRUD", "activa")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := repo.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Test CRUD", result.Titulo)
}

func TestRepository_GetAll(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "titulo"}).
		AddRow(1, "Solicitud 1").
		AddRow(2, "Solicitud 2")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	results, err := repo.GetAll(context.Background(), GetAllReq{})
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestRepository_Update(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), 1, UpdateReq{
		Titulo: stringPtr("Updated"),
	})
	assert.NoError(t, err)
}

func TestRepository_Delete(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	// Configuramos la expectativa para el soft delete
	mock.ExpectBegin()
	mock.ExpectExec("(?i)UPDATE `solicitudes` SET `deleted_at`=\\? WHERE `solicitudes`\\.`id` = \\? AND `solicitudes`\\.`deleted_at` IS NULL").
		WithArgs(sqlmock.AnyArg(), uint(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSolicitud_TableName(t *testing.T) {
	s := &Solicitud{}
	assert.Equal(t, "solicitudes", s.TableName())
}

func TestRepository_GetAll_WithFilters(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	// Test con múltiples filtros
	rows := sqlmock.NewRows([]string{"id", "titulo", "estado", "area"}).
		AddRow(1, "DevOps", "activa", "IT")

	// GORM genera SQL con parámetros específicos
	mock.ExpectQuery("SELECT \\* FROM `solicitudes` WHERE titulo LIKE \\? AND estado LIKE \\? AND area LIKE \\? AND `solicitudes`\\.`deleted_at` IS NULL").
		WithArgs("%DevOps%", "%activa%", "%IT%").
		WillReturnRows(rows)

	filters := GetAllReq{
		Titulo: "DevOps",
		Estado: "activa",
		Area:   "IT",
	}

	results, err := repo.GetAll(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "DevOps", results[0].Titulo)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAll_WithPagination(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "titulo"}).
		AddRow(3, "Solicitud 3").
		AddRow(4, "Solicitud 4")

	// GORM añade automáticamente deleted_at IS NULL para soft delete
	mock.ExpectQuery("SELECT \\* FROM `solicitudes` WHERE `solicitudes`\\.`deleted_at` IS NULL LIMIT \\? OFFSET \\?").
		WithArgs(2, 2).
		WillReturnRows(rows)

	filters := GetAllReq{
		Limit: 2,
		Page:  2,
	}

	results, err := repo.GetAll(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAll_WithNumericFilters(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id", "numero_vacantes", "renta_desde"}).
		AddRow(1, 5, 1000000)

	mock.ExpectQuery("SELECT \\* FROM `solicitudes` WHERE numero_vacantes = \\? AND renta_desde = \\? AND `solicitudes`\\.`deleted_at` IS NULL").
		WithArgs(5, 1000000).
		WillReturnRows(rows)

	filters := GetAllReq{
		NumeroVacantes: 5,
		RentaDesde:     1000000,
	}

	results, err := repo.GetAll(context.Background(), filters)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAll_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	mock.ExpectQuery("SELECT").WillReturnError(assert.AnError)

	results, err := repo.GetAll(context.Background(), GetAllReq{})
	assert.Error(t, err)
	assert.Nil(t, results)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetByID_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	mock.ExpectQuery("SELECT").WillReturnError(assert.AnError)

	result, err := repo.GetByID(context.Background(), 999)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Update(context.Background(), 1, UpdateReq{
		Titulo: stringPtr("Failed Update"),
	})
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Update_AllFields(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test actualizando todos los campos posibles
	numeroVacantes := 10
	rentaDesde := 500000
	rentaHasta := 800000
	err := repo.Update(context.Background(), 1, UpdateReq{
		Titulo:              stringPtr("Full Update"),
		Estado:              stringPtr("completada"),
		Area:                stringPtr("Marketing"),
		Pais:                stringPtr("Argentina"),
		NumeroVacantes:      &numeroVacantes,
		RentaDesde:          &rentaDesde,
		RentaHasta:          &rentaHasta,
		ModalidadTrabajo:    stringPtr("remoto"),
		TipoServicio:        stringPtr("consultoria"),
		NivelExperiencia:    stringPtr("senior"),
		FechaInicioProyecto: stringPtr("2025-01-15"),
	})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Delete_Error(t *testing.T) {
	db, mock := setupTestDB(t)
	repo := NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE.*SET.*deleted_at").WillReturnError(assert.AnError)
	mock.ExpectRollback()

	err := repo.Delete(context.Background(), 999)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSolicitud_ToResponse(t *testing.T) {
	s := &Solicitud{
		ID:     1,
		Titulo: "Test Solicitud",
		Estado: "activa",
	}

	response := s.ToResponse()
	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "Test Solicitud", response.Titulo)
	assert.Equal(t, "activa", response.Estado)
}
