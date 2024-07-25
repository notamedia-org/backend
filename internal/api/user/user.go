package user

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
)

type User struct {
	tableName struct{} `pg:"user"`
	ID        string   `pg:"id"`
	FirstName string   `pg:"first_name"`
	LastName  string   `pg:"last_name"`
}

func (u *User) findById(db *pg.DB) (*User, error) {
	if err := db.Model(&u).WherePK().Select(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) UpdateById(db *pg.DB) error {
	result, err := db.Model(u).WherePK().Update()
	if err != nil {
		return fmt.Errorf("error while updating user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found with id: %s", u.ID)
	}

	return nil
}

func (u *User) CreateUser(db *pg.DB) error {
	_, err := db.Model(u).Insert()
	if err != nil {
		return fmt.Errorf("couldnt create user: %w", err)
	}

	return nil
}

func (u *User) DeleteById(db *pg.DB) error {
	result, err := db.Model(u).WherePK().Delete()
	if err != nil {
		return fmt.Errorf("couldnt delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("not found user with id: %s", u.ID)
	}

	return nil
}

func (u *User) GetUUID() {
	u.ID = uuid.New().String()
}
