package dao

import (
	"github.com/HsimWong/ecommerce/internal/model"
)

type UserDAO interface {
	GetByID(id int) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id int) error
	CheckEmailExistence(email string) (bool, error)
	CheckUsernameExistence(username string) (bool, error)
}

func NewUserDAO() UserDAO {
	return &UserDAOImpl{}
}

type UserDAOImpl struct {
}

// ud *UserDAOImpl
func (ud *UserDAOImpl) GetByID(id int) (*model.User, error)                { panic("Not Implemented") }
func (ud *UserDAOImpl) GetByUsername(username string) (*model.User, error) { panic("Not Implemented") }
func (ud *UserDAOImpl) Create(user *model.User) error                      { panic("Not Implemented") }
func (ud *UserDAOImpl) Update(user *model.User) error                      { panic("Not Implemented") }
func (ud *UserDAOImpl) Delete(id int) error                                { panic("Not Implemented") }
func (ud *UserDAOImpl) CheckEmailExistence(email string) (bool, error)     { panic("Not Implemented") }
func (ud *UserDAOImpl) CheckUsernameExistence(username string) (bool, error) {
	panic("Not Implemented")
}
