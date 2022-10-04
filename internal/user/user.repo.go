package user

import (
	"context"
	"database/sql"
	"fmt"

	entsql "entgo.io/ent/dialect/sql"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent"
	"git.teko.vn/loyalty-system/loyalty-file-processing/internal/ent/ent/user"
)

type (
	Repo interface {
		CreateUser(ctx context.Context, user *User) (*User, error)
		GetUser(ctx context.Context, username string) (*User, error)
	}

	repoImpl struct {
		client *ent.Client
	}
)

const dbEngine = "mysql"

var _ Repo = &repoImpl{} // only for mention that repoImpl implement Repo

func NewRepo(db *sql.DB) *repoImpl {
	drv := entsql.OpenDB(dbEngine, db)
	client := ent.NewClient(ent.Driver(drv))
	return &repoImpl{client: client}
}

func (r *repoImpl) CreateUser(ctx context.Context, user *User) (*User, error) {
	u, err := mapUserCreate(r.client, *user).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}
	return &User{*u}, nil
}

func (r *repoImpl) GetUser(ctx context.Context, phoneNumber string) (*User, error) {
	u, err := r.client.User.
		Query().
		Where(user.PhoneNumber(phoneNumber)).
		// `Only` fails if no user found,
		// or more than 1 user returned.
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying singular user: %w", err)
	}

	return &User{*u}, nil
}

// private method ----------

// mapFileAwardPoint ... must update this function when schema is modified
func mapUserCreate(client *ent.Client, user User) *ent.UserCreate {
	return client.User.
		Create().
		SetName(user.Name).
		SetActive(user.Active).
		SetEmail(user.Email).
		SetPhoneNumber(user.PhoneNumber).
		SetPasswordHash(user.PasswordHash)
}
