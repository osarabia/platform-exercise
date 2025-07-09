package services

import (
	"context"
	"platform-exercise/internal/entities"
	repository "platform-exercise/internal/infra/gorm"
)

type UserService struct {
	userRepository *repository.UserRepository
	AuthService    *AuthService
}

func NewUserService(
	userRepository *repository.UserRepository,
	tokenService *AuthService,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		AuthService:    tokenService,
	}
}

func (u *UserService) NewUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	err := u.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return u.userRepository.FindUserByEmail(ctx, user.Email)
}

func (u *UserService) GetUser(ctx context.Context, ID string) (*entities.User, error) {
	return u.userRepository.FindUserByID(ctx, ID)
}

func (u *UserService) LogInUser(ctx context.Context, user *entities.User) (string, error) {
	password := user.Password
	user, err := u.userRepository.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil
	}
	if u.userRepository.ValidatePassword(user, password) {
		token, err := u.AuthService.GenerateJWT(*user.ID)
		if err != nil {
			return "", nil
		}
		return token, nil
	}

	return "", nil
}

func (u *UserService) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	user, err := u.userRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) DeleteUser(ctx context.Context, ID string) error {
	return u.userRepository.Delete(ctx, ID)
}
