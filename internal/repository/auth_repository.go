package repository

import (
	"github.com/FranzSinaga/blogcms/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *domain.User) error {
	query := `
		INSERT INTO users (email, password, name, role)
		values (:email, :password, :name, :role)
	`

	_, err := r.db.NamedExec(query, user)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.Get(user, query, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}
