package account

import (
	"context"
	"strconv"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/repository/connections"
)

type UserStoragePG struct {
	conn *connections.DBPgConn
}

func NewPG(DBPgConn *connections.DBPgConn) *UserStoragePG {
	return &UserStoragePG{
		conn: DBPgConn,
	}
}

func (u *UserStoragePG) UserCreate(ctx context.Context, Name string, Email string, Password string) (*model.User, error) {
	user := model.User{
		// todo id
		Name:     Name,
		Email:    Email,
		Password: Password,
	}
	u.conn.Db.Model(&user).Insert()
	return &user, nil
}

func (u *UserStoragePG) UserGetByID(ctx context.Context, id string) (*model.User, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	user := &model.User{Id: int64(i)}
	err = u.conn.Db.Model(user).WherePK().Select()
	if err != nil {
		return nil, model.ErrUserNotFound
	}
	if user == nil {
		return nil, model.ErrUserNotFound
	}
	return user, nil
}

func (u *UserStoragePG) UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := u.conn.Db.Model(user).Where("email = ?", email).Select()
	if err != nil {
		return nil, model.ErrUserNotFound
	}
	return user, nil
}
