package user

import (
	"context"
	"database/sql"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/mitchellh/mapstructure"

	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/user"
)

type Repo interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUser(ctx context.Context, username string) (*User, error)
}

type RepoImpl struct {
	client *ent.Client
}

func NewRepo(db *sql.DB) *RepoImpl {
	drv := entsql.OpenDB("mysql", db)
	client := ent.NewClient(ent.Driver(drv))
	return &RepoImpl{client: client}
}

func (r *RepoImpl) CreateUser(ctx context.Context, user *User) (*User, error) {
	u, err := r.client.User.
		Create().
		SetName(user.Name).
		SetActive(user.Active).
		SetEmail(user.Email).
		SetPhoneNumber(user.PhoneNumber).
		SetPasswordHash(user.PasswordHash).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	user.ID = u.ID
	return user, nil
}

func (r *RepoImpl) GetUser(ctx context.Context, phonenumber string) (*User, error) {
	u, err := r.client.User.
		Query().
		Where(user.PhoneNumber(phonenumber)).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying singular user: %w", err)
	}

	user := &User{}
	if err := mapstructure.Decode(u, user); err != nil {
		return nil, fmt.Errorf("failed decode user info from DB")
	}
	return user, nil
}
