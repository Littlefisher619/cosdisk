package account

import (
	"context"
	"fmt"

	"github.com/Littlefisher619/cosdisk/model"
)

type UserStorageMap struct {
	// [id]user
	accountmap map[string]*model.User
}

// need database or lock
// for test only
var id int64 = 1

func NewMap() *UserStorageMap {
	m := &UserStorageMap{
		accountmap: make(map[string]*model.User),
	}
	return m
}

func (u *UserStorageMap) UserCreate(ctx context.Context, Name string, Email string, Password string) (*model.User, error) {
	user := model.User{
		// todo id
		Id:       id,
		Name:     Name,
		Email:    Email,
		Password: Password,
	}
	id = id + 1
	u.accountmap[fmt.Sprint(user.Id)] = &user
	return &user, nil
}

func (u *UserStorageMap) UserGetByID(ctx context.Context, id string) (*model.User, error) {
	user := u.accountmap[id]
	if user == nil {
		return nil, model.ErrUserNotFound
	}
	return user, nil
}

func (u *UserStorageMap) UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, v := range u.accountmap {
		if v.Email == email {
			return v, nil
		}
	}
	return nil, model.ErrUserNotFound
}
