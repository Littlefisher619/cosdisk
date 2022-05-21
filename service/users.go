package service

import (
	"context"
	"fmt"

	"github.com/Littlefisher619/cosdisk/model"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(s string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(hashed)
}

func comparePassword(hashed string, normal string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
}

func (c *CosDisk) UserRegister(
	ctx context.Context, name string, email string, password string,
) (user *model.User, err error) {
	c.Logger.Info("UserRegister " + name + " " + email)
	var createdUser *model.User

	err = c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		// Check Email
		_, err = c.accountRepository.UserGetByEmail(ctx, email)
		if err == nil {
			// exists
			return ErrUserAlreadyRegister
		}

		password = hashPassword(password)
		createdUser, err = c.accountRepository.UserCreate(ctx, name, email, password)
		if err != nil {
			return err
		}

		err = txn.AddDir(fmt.Sprint(createdUser.Id), "/")
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (c *CosDisk) UserLogin(
	ctx context.Context, email string, password string,
) (user *model.User, err error) {
	c.Logger.Info("UserLogin " + email)

	gotUser, err := c.accountRepository.UserGetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := comparePassword(gotUser.Password, password); err != nil {
		return nil, ErrPasswordIncorrect
	}

	return gotUser, nil
}
