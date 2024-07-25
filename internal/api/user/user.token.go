package user

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
)

type UserToken struct {
	tableName struct{} `pg:"user_tokens"`
	ID        string
	UserID    string
	Token     string
}

func (uToken *UserToken) FindUserByToken(db *pg.DB) (*User, error) {
	if err := db.Model(uToken).Where("token = ?", uToken.Token).Select(); err != nil {
		return nil, err
	}

	user := &User{ID: uToken.UserID}
	if err := db.Model(user).WherePK().Select(); err != nil {
		return nil, err
	}

	return user, nil
}

func (uToken *UserToken) UpdateById(db *pg.DB) error {
	_, err := db.Model(uToken).WherePK().Update()
	if err != nil {
		return fmt.Errorf("error while updating user: %w", err)
	}

	return nil
}

func (uToken *UserToken) CreateUserToken(db *pg.DB) error {
	_, err := db.Model(uToken).Insert()
	if err != nil {
		return fmt.Errorf("couldn1t create user: %w", err)
	}

	return nil
}

func (uToken *UserToken) DeleteById(db *pg.DB) error {
	_, err := db.Model(uToken).WherePK().Delete()
	if err != nil {
		return fmt.Errorf("couldnt delete user: %w", err)
	}

	return nil
}

func (uToken *UserToken) GetUUID() {
	uToken.ID = uuid.New().String()
}
