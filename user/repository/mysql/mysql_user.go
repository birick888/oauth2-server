package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/menduong/oauth2/domain"
	"github.com/menduong/oauth2/user/repository"
)

type mysqlUserRepository struct {
	Conn *sql.DB
}

// NewMysqlUserRepository will create an object that represent the article.Repository interface
func NewMysqlUserRepository(Conn *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{Conn}
}

func (m *mysqlUserRepository) fetch(ctx context.Context,
	query string,
	args ...interface{}) (result []domain.User, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.User, 0)
	for rows.Next() {
		t := domain.User{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.Email,
			&t.Password,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlUserRepository) Fetch(ctx context.Context,
	cursor string, num int64) (res []domain.User, nextCursor string, err error) {
	query := `SELECT id, username, email, password, updated_at, created_at
  						FROM user WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (m *mysqlUserRepository) GetByID(ctx context.Context, id string) (res domain.User, err error) {
	query := `SELECT id, username, email, password, updated_at, created_at FROM user WHERE ID = ?`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) <= 0 {
		return domain.User{}, domain.ErrNotFound
	}
	res = list[0]
	return res, nil
}

func (m *mysqlUserRepository) GetByEmail(ctx context.Context, email string) (res domain.User, err error) {
	query := `SELECT id, username, email, password, updated_at, created_at
  						FROM user WHERE email = ?`

	list, err := m.fetch(ctx, query, email)
	if err != nil {
		return domain.User{}, err
	}

	if len(list) <= 0 {
		return domain.User{}, domain.ErrNotFound
	}
	res = list[0]

	return res, nil
}

func (m *mysqlUserRepository) Store(ctx context.Context, user *domain.User) (id string, err error) {
	query := `INSERT INTO user(id, username, email, password, updated_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return "", err
	}

	// Generate uuid v4
	newUUID := uuid.New().String()

	_, err = stmt.ExecContext(ctx, newUUID, user.Username, user.Email, user.Password, user.UpdatedAt, user.CreatedAt)

	return newUUID, err
}

func (m *mysqlUserRepository) Delete(ctx context.Context, id string) (err error) {
	query := "DELETE FROM user WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("weird  behavior. total affected: %d", rowsAfected)
		return
	}

	return
}
func (m *mysqlUserRepository) Update(ctx context.Context, ar *domain.User) (err error) {
	query := `UPDATE user set username=?, email=?, password=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Username, ar.Email, ar.Password, ar.UpdatedAt, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  behavior. total affected: %d", affect)
		return
	}

	return
}
