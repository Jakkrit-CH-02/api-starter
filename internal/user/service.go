package user

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"api-starter/domain/entity"
	"api-starter/domain/model"

	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength = 8

var (
	ErrNotFound   = errors.New("user not found")
	ErrEmailTaken = errors.New("email is already in use")
)

// ValidationError describes a request that failed business validation.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

type (
	// Service holds the business logic of the user module.
	Service interface {
		Create(ctx context.Context, req model.UserRequest) (*model.UserResponse, error)
		GetByID(ctx context.Context, id uint) (*model.UserResponse, error)
		List(ctx context.Context) ([]model.UserResponse, error)
	}
	service struct {
		repository Repository
	}
)

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, req model.UserRequest) (*model.UserResponse, error) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if err := validate(req); err != nil {
		return nil, err
	}

	existing, err := s.repository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hash),
	}
	if req.RoleID != 0 {
		user.UserRoleID = &req.RoleID
	}

	if err := s.repository.Create(ctx, &user); err != nil {
		return nil, err
	}

	res := model.NewUserResponse(user)
	return &res, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*model.UserResponse, error) {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	res := model.NewUserResponse(*user)
	return &res, nil
}

func (s *service) List(ctx context.Context) ([]model.UserResponse, error) {
	users, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]model.UserResponse, 0, len(users))
	for _, user := range users {
		res = append(res, model.NewUserResponse(user))
	}
	return res, nil
}

func validate(req model.UserRequest) error {
	switch {
	case req.FirstName == "":
		return &ValidationError{"firstName is required"}
	case req.LastName == "":
		return &ValidationError{"lastName is required"}
	case !validEmail(req.Email):
		return &ValidationError{"email is invalid"}
	case len(req.Password) < minPasswordLength:
		return &ValidationError{"password must be at least 8 characters"}
	case len(req.Password) > 72:
		// bcrypt only accepts up to 72 bytes.
		return &ValidationError{"password must be at most 72 bytes"}
	}
	return nil
}

func validEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email
}
