package common

import (
	"bytes"
	"fmt"
	"math/rand"
	"mime/quotedprintable"
	"net/smtp"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Token struct define
type Token struct {
	UserID int64  `json:"userid" form:"userid" query:"userid"`
	Token  string `json:"token" form:"token" query:"token"`
}

// HashPassword is
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// IsMatchedPassword is
func IsMatchedPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateToken is
func CreateToken(userid int64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

// RangeIn is
func RangeIn(low, hi int) int {
	rand.Seed(time.Now().UnixNano())
	return low + rand.Intn(hi-low)
}

// VerifyToken is
func VerifyToken(token string) (bool, error) {
	return true, nil
}

// SendEmail is
func SendEmail(email string, otp int) (err error) {

	viper.SetConfigFile(`config.json`)
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var (
		from       = viper.GetString(`smtp.user`)
		password   = viper.GetString(`smtp.password`)
		recipients = []string{email}
		host       = viper.GetString(`smtp.host`)
		port       = viper.GetString(`smtp.port`)
	)

	auth := smtp.PlainAuth("", from, password, host)

	msg := viper.GetString(`smtp.msg`)

	msg = fmt.Sprintf(msg, email, otp)

	header := make(map[string]string)
	header["From"] = from
	header["To"] = email
	header["Subject"] = "Forgot password"

	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", "text/html")
	header["Content-Disposition"] = "inline"
	header["Content-Transfer-Encoding"] = "quoted-printable"

	headerMessage := ""
	for key, value := range header {
		headerMessage += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	body := "<h3>" + msg + "</h3>"
	var bodyMessage bytes.Buffer
	temp := quotedprintable.NewWriter(&bodyMessage)
	temp.Write([]byte(body))
	temp.Close()

	finalMessage := headerMessage + "\r\n" + bodyMessage.String()

	err = smtp.SendMail(host+":"+port, auth, from, recipients, []byte(finalMessage))
	if err != nil {
		return err
	}

	return nil
}
