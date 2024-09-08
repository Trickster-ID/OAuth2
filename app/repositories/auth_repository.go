package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"oauth2/app/models"
	"time"
)

type IAuthRepository interface {
	GetUserByUsernameOrEmail(username, email string) (*models.Users, error)
}

type authRepository struct {
	db    *sql.DB
	mongo mongo.Database
}

func NewAuthRepository(db *sql.DB, mongo mongo.Database) IAuthRepository {
	return &authRepository{
		db:    db,
		mongo: mongo,
	}
}

func (r *authRepository) GetUserByUsernameOrEmail(username, email string) (*models.Users, error) {
	user := &models.Users{}
	whereClause := ""
	valueOfWhereCluase := ""
	if username != "" {
		whereClause = "username = ?"
		valueOfWhereCluase = username
	} else if email != "" {
		whereClause = "email = ?"
		valueOfWhereCluase = email
	} else {
		err := errors.New("username or email is empty")
		logrus.Error(err)
		return nil, err
	}
	err := r.db.QueryRow(fmt.Sprintf(`
	select 
    id, username, email, password_hash, is_admin, created_at, updated_at 
	from users 
	where %s
	`, whereClause), valueOfWhereCluase).Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return user, nil
}

func (r *authRepository) AssignToken(userId int64, token string, expired *time.Time, ctx context.Context) error {
	
	return nil
}
