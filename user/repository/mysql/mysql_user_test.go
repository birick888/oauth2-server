package mysql_test

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/menduong/oauth2/domain"
	"github.com/menduong/oauth2/user/repository"
	userMysqlRepo "github.com/menduong/oauth2/user/repository/mysql"
	"github.com/stretchr/testify/assert"
)

var userRecord domain.User

func init() {
	// init data test
	userRecord.ID = "cda6498a-235d-4f7e-ae19-661d41bc154c"
	userRecord.Username = "binhdc"
	userRecord.Email = "abc@gmail.com"
	userRecord.Password = "pabc123"
	userRecord.CreatedAt = time.Now()
	userRecord.UpdatedAt = time.Now()
}

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockUsers := []domain.User{
		{
			ID: "cda6498a-235d-4f7e-ae19-661d41bc154c", Username: "username1", Email: "email1@gmail.com",
			Password: "Password1", UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		{
			ID: "cda6498a-235d-4f7e-ae19-661d41bc154d", Username: "username2", Email: "email2@gmail.com",
			Password: "Password2", UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "updated_at", "created_at"}).
		AddRow(mockUsers[0].ID, mockUsers[0].Username, mockUsers[0].Email,
			mockUsers[0].Password, mockUsers[0].UpdatedAt, mockUsers[0].CreatedAt).
		AddRow(mockUsers[1].ID, mockUsers[1].Username, mockUsers[1].Email,
			mockUsers[1].Password, mockUsers[1].UpdatedAt, mockUsers[1].CreatedAt)

	query := "SELECT id, username, email, password, updated_at, created_at FROM user " +
		"WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	userRepo := userMysqlRepo.NewMysqlUserRepository(db)
	cursor := repository.EncodeCursor(mockUsers[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := userRepo.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "updated_at", "created_at"}).
		AddRow("cda6498a-235d-4f7e-ae19-661d41bc154c", "username1", "email1@gmail.com", "password1", time.Now(), time.Now())

	query := "SELECT id, username, email, password, updated_at, created_at FROM user WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	userRepo := userMysqlRepo.NewMysqlUserRepository(db)

	userRecord, err := userRepo.GetByID(context.TODO(), "cda6498a-235d-4f7e-ae19-661d41bc154c")
	assert.NoError(t, err)
	assert.NotNil(t, userRecord)
}

func TestStore(t *testing.T) {
	now := time.Now()
	newUserID := uuid.New().String()
	userData := &domain.User{
		ID:        newUserID,
		Username:  "username100",
		Email:     "email100@gmail.com",
		Password:  "password100",
		CreatedAt: now,
		UpdatedAt: now,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT INTO user"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(newUserID,
		userData.Username,
		userData.Email,
		userData.Password,
		userData.UpdatedAt,
		userData.CreatedAt)

	repo := userMysqlRepo.NewMysqlUserRepository(db)

	userID, _ := repo.Store(context.TODO(), userData)
	assert.NotNil(t, userID)
}

func TestGetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "updated_at", "created_at"}).
		AddRow("cda6498a-235d-4f7e-ae19-661d41bc154c", "username1", "email1@gmail.com", 1, now, now)

	query := "SELECT id, username, email, password, updated_at, created_at FROM user WHERE email = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	user := userMysqlRepo.NewMysqlUserRepository(db)

	email := "email1@gmail.com"
	record, err := user.GetByEmail(context.TODO(), email)
	assert.NoError(t, err)
	assert.NotNil(t, record)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM user WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs("cda6498a-235d-4f7e-ae19-661d41bc154c").WillReturnResult(sqlmock.NewResult(2, 1))

	user := userMysqlRepo.NewMysqlUserRepository(db)

	userID := "cda6498a-235d-4f7e-ae19-661d41bc154c"
	err = user.Delete(context.TODO(), userID)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	userData := &domain.User{
		ID:        "cda6498a-235d-4f7e-ae19-661d41bc154c",
		Username:  "user11",
		Email:     "email11",
		Password:  "pass11",
		CreatedAt: now,
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE user set username=\\?, email=\\?, password=\\?, updated_at=\\? WHERE ID = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(userData.Username,
		userData.Email, userData.Password, userData.UpdatedAt, userData.ID).WillReturnResult(sqlmock.NewResult(11, 1))

	repo := userMysqlRepo.NewMysqlUserRepository(db)

	err = repo.Update(context.TODO(), userData)
	assert.NoError(t, err)
}
