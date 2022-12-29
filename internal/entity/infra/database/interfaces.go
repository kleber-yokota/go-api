package database

import "github.com/kleber-yokota/go-api/internal/entity"



type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(emaild string) (*entity.User,error)
}

type ProdudctInterface interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, sort string) ([]entity.Product, error)
	FindByID(id string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}