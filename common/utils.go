package common

import (
	"bytes"
	"fmt"
	"math/rand"
	"mime/quotedprintable"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// Token struct define
type Token struct {
	UserID int64  `json:"userid" form:"userid" query:"userid"`
	Token  string `json:"token" form:"token" query:"token"`
}

var (
	header     map[string]string
	from       string
	password   string
	recipients []string
	host       string
	port       string
	msg        string
	auth       smtp.Auth
)

func init() {
	LoadConfig()
	ConfigLogrus()
	ConfigSMTP()
}

func LoadConfig() {
	err := godotenv.Load(filepath.Join("./env", "test.env"))
	if err != nil {
		logrus.Fatalf("Error load env file. Err: %s", err)
		os.Exit(2)
	}

	viper.SetConfigFile(`config.json`)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Error load config file. Err: %s", err)
		os.Exit(2)
	}
}

func ConfigLogrus() {
	rotationCount := uint(viper.GetInt("log.rotate.rotationCount"))
	writer, err := rotatelogs.New(
		viper.GetString("log.path")+viper.GetString("log.logPattern"),
		rotatelogs.WithLinkName(viper.GetString("log.path")),
		rotatelogs.WithRotationTime(time.Duration(viper.GetInt("log.rotate.rotationTime"))*time.Minute),
		rotatelogs.WithRotationCount(rotationCount),
	)

	if err != nil {
		logrus.Error(err)
		os.Exit(2)
	}
	logrus.SetReportCaller(true)
	logrus.SetOutput(writer)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   viper.GetBool("log.disableCorlors"),
		TimestampFormat: viper.GetString("timeStampFormat"),
		FullTimestamp:   viper.GetBool("log.fullTimestamp"),
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// this function is required when you want to introduce your custom format.
			// In my case I wanted file and line to look like this `file="engine.go:141`
			// but f.File provides a full path along with the file name.
			// So in `formatFilePath()` function I just trimmed everything before the file name
			// and added a line number in the end
			return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
		},
	})
}

func ConfigSMTP() {
	host = viper.GetString(`smtp.host`)
	port = viper.GetString(`smtp.port`)
	from = os.Getenv("SMTP_USER")
	password = os.Getenv("SMTP_APP_PASSWORD")

	header = make(map[string]string)
	header["From"] = from
	header["Subject"] = "[oauth2-server]Forgot password"
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", "text/html")
	header["Content-Disposition"] = "inline"
	header["Content-Transfer-Encoding"] = "quoted-printable"

	msg = viper.GetString(`smtp.msg`)
	auth = smtp.PlainAuth("", from, password, host)
}

// HashPassword is
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * time.Duration(viper.GetInt("token.expire"))).Unix()
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
	msgContent := fmt.Sprintf(msg, email, otp)
	recipients = []string{email}
	headerMessage := ""
	for key, value := range header {
		headerMessage += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	body := "<h3>" + msgContent + "</h3>"
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
