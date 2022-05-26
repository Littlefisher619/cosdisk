package account

import (
	"context"
	"strconv"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/pkg/dbdriver"
)

type UserStorageXORM struct {
	conn *dbdriver.DBEngine
}

func NewXORM(DBPgConn *dbdriver.DBEngine) *UserStorageXORM {
	return &UserStorageXORM{
		conn: DBPgConn,
	}
}

func (u *UserStorageXORM) UserCreate(ctx context.Context, Name string, Email string, Password string) (*model.User, error) {
	user := model.User{
		// todo id
		Name:     Name,
		Email:    Email,
		Password: Password,
	}
	u.conn.Db.Insert(&user)
	return &user, nil
}

func (u *UserStorageXORM) UserGetByID(ctx context.Context, id string) (*model.User, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	user := &model.User{Id: int64(i)}
	has, err := u.conn.Db.Get(user)
	if err != nil || !has {
		return nil, model.ErrUserNotFound
	}
	if user == nil {
		return nil, model.ErrUserNotFound
	}
	return user, nil
}

func (u *UserStorageXORM) UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	has, err := u.conn.Db.Where("email = ?", email).Get(user)
	if err != nil || !has {
		return nil, model.ErrUserNotFound
	}
	return user, nil
}
