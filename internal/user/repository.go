package user

import (
	"context"
	"errors"

	"api-starter/domain/entity"

	"gorm.io/gorm"
)

type (
	// Repository is the data access layer of the user module.
	// The Find* methods return (nil, nil) when no row matches.
	Repository interface {
		Create(ctx context.Context, user *entity.User) error
		FindByID(ctx context.Context, id uint) (*entity.User, error)
		FindByEmail(ctx context.Context, email string) (*entity.User, error)
		FindAll(ctx context.Context) ([]entity.User, error)
	}
	repository struct {
		db *gorm.DB
	}
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return oneOrNil(&user, err)
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return oneOrNil(&user, err)
}

func (r *repository) FindAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).Order("id").Find(&users).Error
	return users, err
}

func oneOrNil(user *entity.User, err error) (*entity.User, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
