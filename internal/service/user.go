package service

import (
	"github.com/HsimWong/ecommerce/internal/dao"
	"github.com/HsimWong/ecommerce/internal/utils"
)

type UserService struct {
	userDao dao.UserDAO
	// jwtSecret string
}

func NewUserService( /* 依赖注入 */ ) *UserService {
	return &UserService{
		userDao: dao.NewUserDAO(),
		// jwtSecret: config.AppConfig.JWT.Secret,
	}
}

func (us *UserService) Register(username, email, password string) error {
	// Validate existence of username and email
	existence, err := us.userDao.CheckUsernameExistence(username)
	if err != nil {
		return err
	}
	if existence {
		return utils.ErrRegisterExistedUsername
	}

	existence, err = us.userDao.CheckEmailExistence(email)
	if err != nil {
		return err
	}
	if existence {
		return utils.ErrRegisterExistedEmail
	}

	return nil
}
