package services

import (
	"context"
	"testing"
	m "ticketing-system/internal/core/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func (r *mockUserRepo) Create(ctx context.Context, user *m.User) error {
  	args := r.Called(ctx, user)
	return args.Error(0)
}

func (r *mockUserRepo) Update(ctx context.Context, user *m.User) error {
	args := r.Called(ctx, user)
	return args.Error(0)
}

func (r *mockUserRepo) SetActive(ctx context.Context, id uint, active bool) error {
	args := r.Called(ctx, id, active)
	return args.Error(0)
}

func (r *mockUserRepo) GetByID(ctx context.Context, id uint) (*m.User, error) {
	args := r.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*m.User), args.Error(1)
}

func (r *mockUserRepo) GetByIDForUpdate(ctx context.Context, id uint) (*m.User, error) {
	args := r.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*m.User), args.Error(1)
}

func (r *mockUserRepo) GetByUsername(ctx context.Context, username string) (*m.User, error) {
	args := r.Called(ctx, username)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*m.User), args.Error(1)
}

func (r *mockUserRepo) List(ctx context.Context) ([]m.User, error) {
	args := r.Called(ctx)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]m.User), args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name          string
		username 	  string
		password 	  string
		role 		  string
		firstName 	  string
		lastName 	  string
		phoneNumber   string
		mockRepoSetup func(mockRepo *mockUserRepo)
		expectError   bool
		expectMessage string
	}{
		{
			name:        "success",
			username:    "John",
			password:    "12345678",
			role:        "CASHIER",
			firstName:   "John",
			lastName:    "Doe",
			phoneNumber: "0812345678",
			mockRepoSetup: func(mockRepo *mockUserRepo) {
				mockRepo.On("GetByUsername", mock.Anything, "John").Return(nil, nil)
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "username too short",
			username:    "Jon",
			password:    "12345678",
			role:        "CASHIER",
			firstName:   "John",
			lastName:    "Doe",
			phoneNumber: "0812345678",
			mockRepoSetup: func(mockRepo *mockUserRepo) {},
			expectError:   true,
			expectMessage: "username must be more than 4 characters long",
		},
		{
			name:        "username already exists",
			username:    "John",
			password:    "12345678",
			role:        "CASHIER",
			firstName:   "John",
			lastName:    "Doe",
			phoneNumber: "0812345678",
			mockRepoSetup: func(mockRepo *mockUserRepo) {
				mockUser := &m.User{ID: 99, Username: "John"}
				mockRepo.On("GetByUsername", mock.Anything, "John").Return(mockUser, nil)
			},
			expectError:   true,
			expectMessage: "username already exists",
		},
		{
			name:        "password too short",
			username:    "John",
			password:    "1234",
			role:        "CASHIER",
			firstName:   "John",
			lastName:    "Doe",
			phoneNumber: "0812345678",
			mockRepoSetup: func(mockRepo *mockUserRepo) {
				mockRepo.On("GetByUsername", mock.Anything, "John").Return(nil, nil)
			},
			expectError:   true,
			expectMessage: "password must be at least 8 characters long",
		},
		{
			name:        "invalid role type",
			username:    "John",
			password:    "12345678",
			role:        "SUPER_ADMIN",
			firstName:   "John",
			lastName:    "Doe",
			phoneNumber: "0812345678",
			mockRepoSetup: func(mockRepo *mockUserRepo) {
				mockRepo.On("GetByUsername", mock.Anything, "John").Return(nil, nil)
			},
			expectError:   true,
			expectMessage: "invalid user role",
		},
		{
			name:        "invalid mobile format",
			username:    "John",
			password:    "12345678",
			role:        "CASHIER",
			firstName:   "John",
			lastName:    "Doe",
			phoneNumber: "021234567",
			mockRepoSetup: func(mockRepo *mockUserRepo) {
				mockRepo.On("GetByUsername", mock.Anything, "John").Return(nil, nil)
			},
			expectError:   true,
			expectMessage: "invalid phone number",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockUserRepo)
			tt.mockRepoSetup(mockRepo)

			service := NewUserService(mockRepo) 

			ctx := context.Background()
			
			err := service.RegisterUser(
				ctx,
				tt.username,
				tt.password,
				tt.role,
				tt.firstName,
				tt.lastName,
				tt.phoneNumber,
			)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectMessage)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}