package gorm

import (
	"context"
	"errors"
	"platform-exercise/internal/entities"
	"platform-exercise/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) Create(ctx context.Context, user *entities.User) error {
	userModel := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Birthday: user.Birthday,
	}
	tx := u.db.WithContext(ctx).Create(&userModel)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (u *UserRepository) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := u.db.WithContext(ctx).Where("email = ? and deleted_at is null", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) FindUserByID(ctx context.Context, ID string) (*entities.User, error) {
	var user entities.User
	if err := u.db.WithContext(ctx).Where("id = ? and deleted_at is null", ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) ValidatePassword(user *entities.User, password string) bool {
	userModel := models.User{
		Email:    user.Email,
		Password: user.Password,
	}

	return userModel.CheckPassword(password)
}

func (u *UserRepository) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	if err := u.db.WithContext(ctx).Model(&models.User{}).Where("id = ? and deleted_at is null", user.ID).Updates(map[string]any{
		"name":     user.Name,
		"birthday": user.Birthday,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return u.FindUserByID(ctx, *user.ID)
}

func (u *UserRepository) Delete(ctx context.Context, ID string) error {
	if err := u.db.WithContext(ctx).Where("id = ? and deleted_at is null", ID).Delete(&models.User{}).Error; err != nil {
		return err
	}

	return nil
}
