package authentication

import (
	"github.com/FranzSinaga/blogcms/internal/domain"
	"github.com/jmoiron/sqlx"
)

// RepositoryInterface defines the methods for user repository
type RepositoryInterface interface {
	CreateUser(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

type Repository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) RepositoryInterface {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *domain.User) error {
	query := `
		INSERT INTO users (email, password, name, role)
		values (:email, :password, :name, :role)
	`

	_, err := r.db.NamedExec(query, user)
	return err
}

func (r *Repository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.Get(user, query, email)

	if err != nil {
		return nil, err
	}

	return user, nil
}
