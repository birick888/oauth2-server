package usecase

import (
	"context"
	"time"

	"github.com/menduong/oauth2/common"
	"github.com/menduong/oauth2/domain"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	userOTPRepo    domain.UserOTPRepository
	contextTimeout time.Duration
}

func NewUserUsecase(ur domain.UserRepository,
	userOTP domain.UserOTPRepository,
	timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       ur,
		userOTPRepo:    userOTP,
		contextTimeout: timeout,
	}
}

func (a *userUsecase) GetByID(c context.Context, id string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	resUser, err := a.userRepo.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	// dont send password
	resUser.Password = ""
	return resUser, nil
}

func (a *userUsecase) Update(c context.Context, u *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	password, err := common.HashPassword(u.Password)
	u.Password = password
	u.UpdatedAt = time.Now()
	return a.userRepo.Update(ctx, u)
}

func (a *userUsecase) Store(c context.Context, u *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedUser, err := a.GetByEmail(ctx, u.Email)
	if len(existedUser.ID) > 0 {
		return domain.ErrEmailAlreadyExists
	}

	password, err := common.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = password

	_, err = a.userRepo.Store(ctx, u)
	return
}

func (a *userUsecase) Delete(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	_, err = a.userRepo.GetByID(ctx, id)

	return a.userRepo.Delete(ctx, id)
}

func (a *userUsecase) GetByEmail(c context.Context, email string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.userRepo.GetByEmail(ctx, email)
	return res, nil
}

func (a *userUsecase) GetOTP(c context.Context, email string) (otp string, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	otp, err = a.userOTPRepo.GetOTP(ctx, email)
	return otp, err
}

func (a *userUsecase) SetOTP(c context.Context, email string, otp string, expireTime time.Duration) error {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.userOTPRepo.SetOTP(ctx, email, otp, expireTime)
	return err
}
