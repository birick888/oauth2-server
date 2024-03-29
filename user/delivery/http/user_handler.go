package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/menduong/oauth2/common"
	"github.com/menduong/oauth2/domain"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// UserHandler  represent the httphandler for user
type UserHandler struct {
	UserUsecase domain.UserUsecase
}

// NewUserHandler will initialize the users/ resources endpoint
func NewUserHandler(e *echo.Echo, us domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: us,
	}
	e.POST("/register", handler.Store)
	e.POST("/login", handler.Login)
	e.GET("/user/:id", handler.GetByID)
	e.DELETE("/user/:id", handler.Delete)
	e.POST("/requestOTP", handler.RequestOTP)
	e.POST("/resetPassword", handler.ResetPassword)
}

// GetByID will get user by given id
func (u *UserHandler) GetByID(c echo.Context) error {
	log.Info("Call GetByID")
	idParam := c.Param("id")
	log.Info("idParam: ", idParam)
	if len(idParam) <= 0 {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	ctx := c.Request().Context()

	user, err := u.UserUsecase.GetByID(ctx, idParam)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// Login will check email, password
func (u *UserHandler) Login(c echo.Context) (err error) {
	log.Info("Call login service")
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	email := user.Email
	password := user.Password

	ctx := c.Request().Context()
	userStored, err := u.UserUsecase.GetByEmail(ctx, email)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: domain.ErrEmailOrPasswordNotMatch.Error()})
	}

	// compare password
	compare := common.IsMatchedPassword(password, userStored.Password)
	if !compare {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: domain.ErrEmailOrPasswordNotMatch.Error()})
	}

	// generate token
	token, err := common.CreateToken(userStored.ID)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	// init json value to response
	tokenJSON := &domain.Token{
		UserId:       userStored.ID,
		Token:        token,
		RefreshToken: "",
	}

	return c.JSON(http.StatusCreated, tokenJSON)
}

// RequestOTP will
func (u *UserHandler) RequestOTP(c echo.Context) (err error) {
	log.Info("Call RequestOTP service")
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	email := user.Email

	ctx := c.Request().Context()
	_, err = u.UserUsecase.GetByEmail(ctx, email)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: domain.ErrEmailNotExists.Error()})
	}

	// Generate random OTP number 4 length
	otp := common.RangeIn(1000, 9999)

	// Send email OTP
	err = common.SendEmail(email, otp)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	// Store OTP to redis
	err = u.UserUsecase.SetOTP(ctx, email, strconv.Itoa(otp), time.Duration(viper.GetInt("otp.expire"))*time.Minute)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	msg := "An OTP already sent to your email %s successful"
	msg = fmt.Sprintf(msg, email)

	return c.JSON(http.StatusCreated, msg)
}

// ResetPassword will
func (u *UserHandler) ResetPassword(c echo.Context) (err error) {
	log.Info("Call ResetPassword service")
	email := c.FormValue("email")
	password := c.FormValue("password")
	otp := c.FormValue("otp")

	ctx := c.Request().Context()
	userStored, err := u.UserUsecase.GetByEmail(ctx, email)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: domain.ErrEmailNotExists.Error()})
	}

	// Store OTP to redis
	otpRedis, err := u.UserUsecase.GetOTP(ctx, email)

	if otpRedis != otp {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: domain.ErrOTPWrongOrExpire.Error()})
	}

	// Update new password to DB
	userStored.Password = password
	err = u.UserUsecase.Update(ctx, &userStored)

	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, "Reset password successful")
}

// Store will store the user by given request body
func (u *UserHandler) Store(c echo.Context) (err error) {
	log.Info("Call store user")
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isValidate(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err = u.UserUsecase.Store(ctx, &user)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: err.Error()})
	}
	// dont send password
	user.Password = ""

	return c.JSON(http.StatusCreated, user)
}

// Delete will delete user by given param
func (u *UserHandler) Delete(c echo.Context) error {
	idParam := c.Param("id")
	if len(idParam) <= 0 {
		return c.JSON(http.StatusNotFound, domain.ErrNotFound.Error())
	}

	ctx := c.Request().Context()

	err := u.UserUsecase.Delete(ctx, idParam)
	if err != nil {
		return c.JSON(domain.GetStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func isValidate(m *domain.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}
