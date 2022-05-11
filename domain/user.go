package domain

import (
	"context"
	"time"
)

// User ...
type User struct {
	ID        int64     `json:"id" form:"id" query:"id"`
	Username  string    `json:"username" form:"username" query:"username"`
	Email     string    `json:"email" form:"email" query:"email"`
	Password  string    `json:"password" form:"password" query:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserUsecase represent the user's usecases
type UserUsecase interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, us *User) error
	Store(context.Context, *User) error
	Delete(ctx context.Context, id int64) error
	GetOTP(ctx context.Context, email string) (string, error)
	SetOTP(ctx context.Context, email string, otp string, expireTime time.Duration) error
}

// UserRepository ...
type UserRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, us *User) error
	Store(context.Context, *User) error
	Delete(ctx context.Context, id int64) error
}

// UserOTPRepository ...
type UserOTPRepository interface {
	GetOTP(ctx context.Context, email string) (string, error)
	SetOTP(ctx context.Context, email string, otp string, expireTime time.Duration) error
}
