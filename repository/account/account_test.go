package account

import (
	"context"
	"testing"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/repository/connections"
	"github.com/stretchr/testify/require"
)

func TestAccountMAP(t *testing.T) {
	ass := require.New(t)
	userDB := NewMap()

	user1, err := userDB.UserCreate(context.Background(), "admin", "123@abc.com", "123456")
	ass.NoError(err)
	ass.Equal(user1.Id, int64(1))

	user2, err := userDB.UserCreate(context.Background(), "fish", "f@sea.com", "88888888")
	ass.NoError(err)
	ass.Equal(user2.Id, int64(2))

	userA, err := userDB.UserGetByID(context.Background(), "1")
	ass.NoError(err)
	ass.Equal(userA.Id, int64(1))
	ass.Equal(userA.Name, "admin")
	ass.Equal(userA.Email, "123@abc.com")
	ass.Equal(userA.Password, "123456")

	userB, err := userDB.UserGetByID(context.Background(), "2")
	ass.NoError(err)
	ass.Equal(userB.Id, int64(2))
	ass.Equal(userB.Name, "fish")
	ass.Equal(userB.Email, "f@sea.com")
	ass.Equal(userB.Password, "88888888")

	userC, err := userDB.UserGetByEmail(context.Background(), "123@abc.com")
	ass.NoError(err)
	ass.Equal(userC.Id, int64(1))
	ass.Equal(userC.Name, "admin")
}

func TestAccountPG(t *testing.T) {
	ass := require.New(t)

	models := []interface{}{
		(*model.User)(nil),
	}
	db, err := connections.NewPostgresDB(models, true)
	ass.NoError(err)
	userDB := NewPG(db)
	defer db.Db.Close()

	user1, err := userDB.UserCreate(context.Background(), "admin", "123@abc.com", "123456")
	ass.NoError(err)
	ass.Equal(user1.Id, int64(1))

	user2, err := userDB.UserCreate(context.Background(), "fish", "f@sea.com", "88888888")
	ass.NoError(err)
	ass.Equal(user2.Id, int64(2))

	userA, err := userDB.UserGetByID(context.Background(), "1")
	ass.NoError(err)
	ass.Equal(userA.Id, int64(1))
	ass.Equal(userA.Name, "admin")
	ass.Equal(userA.Email, "123@abc.com")
	ass.Equal(userA.Password, "123456")

	userB, err := userDB.UserGetByID(context.Background(), "2")
	ass.NoError(err)
	ass.Equal(userB.Id, int64(2))
	ass.Equal(userB.Name, "fish")
	ass.Equal(userB.Email, "f@sea.com")
	ass.Equal(userB.Password, "88888888")

	userC, err := userDB.UserGetByEmail(context.Background(), "123@abc.com")
	ass.NoError(err)
	ass.Equal(userC.Id, int64(1))
	ass.Equal(userC.Name, "admin")
}
