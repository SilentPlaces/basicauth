package users

import (
	"context"
	"database/sql"
	"github.com/SilentPlaces/basicauth.git/internal/models/user"
	"github.com/google/wire"
	"time"
)

type UserRepository interface {
	GetUserByID(id string) (*user.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(dbConnection *sql.DB) UserRepository {
	return &userRepository{db: dbConnection}
}

func (ur *userRepository) GetUserByID(id string) (*user.User, error) {
	var u = user.User{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := ur.db.QueryRowContext(ctx, "SELECT id, name, email, password, created_at FROM users WHERE id=?", id).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

var ProviderSet = wire.NewSet(NewUserRepository)
