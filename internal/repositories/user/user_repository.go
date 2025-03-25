package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/SilentPlaces/basicauth.git/internal/models/models"
	consulService "github.com/SilentPlaces/basicauth.git/internal/services/consul"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper/hash"
	"github.com/google/uuid"
	"github.com/google/wire"
	"time"
)

type UserRepository interface {
	GetUserByID(id string) (*models.User, error)
	GetUserByMail(mail string) (*models.User, error)
	InsertUser(user models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUserByID(id string) error
}

type userRepository struct {
	db      *sql.DB
	timeout time.Duration
}

func NewUserRepository(dbConnection *sql.DB, consul consulService.ConsulService) UserRepository {
	cfg := consul.GetRegistrationConfig()

	return &userRepository{db: dbConnection, timeout: cfg.MailVerificationTimeInSeconds}
}

// Helper function to create a context with timeout
func (ur *userRepository) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// Common function to query a user by a given condition (ID or Mail)
func (ur *userRepository) getUserByCondition(query string, args ...interface{}) (*models.User, error) {
	var u models.User
	ctx, cancel := ur.newContext()
	defer cancel()

	err := ur.db.QueryRowContext(ctx, query, args...).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.IsVerified, &u.VerifiedAt, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	return &u, nil
}

func (ur *userRepository) GetUserByID(id string) (*models.User, error) {
	query := "SELECT id, name, email, password, is_verified, verified_at, created_at FROM users WHERE id=?"
	return ur.getUserByCondition(query, id)
}

func (ur *userRepository) GetUserByMail(mail string) (*models.User, error) {
	query := "SELECT id, name, email, password, is_verified, verified_at, created_at FROM users WHERE email=?"
	return ur.getUserByCondition(query, mail)
}

func (ur *userRepository) InsertUser(user models.User) (*models.User, error) {
	// Generate a new UUID as user id (uuid version 4)
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID: %w", err)
	}
	user.ID = uid.String()

	// Hash the user's password
	user.Password = helpers.HashToSHA1(user.Password)
	//insert into db
	ctx, cancel := ur.newContext()
	defer cancel()
	result, err := ur.db.ExecContext(ctx,
		"INSERT INTO users (id, name, email, password) VALUES (?,?,?,?)",
		user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to determine rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, errors.New("user already exists")
	}
	return &user, nil
}

func (ur *userRepository) DeleteUserByID(id string) error {
	ctx, cancel := ur.newContext()
	defer cancel()

	_, err := ur.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (ur *userRepository) UpdateUser(user *models.User) (*models.User, error) {
	ctx, cancel := ur.newContext()
	defer cancel()
	query := "UPDATE users SET name=?, email=?, password=?, is_verified=?, verified_at=? WHERE id=?"
	result, err := ur.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.IsVerified, user.VerifiedAt, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no user found with id: %d", user.ID)
	}
	updatedUser, err := ur.GetUserByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated user: %w", err)
	}
	return updatedUser, nil
}

var UserRepositoryProviderSet = wire.NewSet(NewUserRepository)
