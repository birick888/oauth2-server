package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	_userHttpDelivery "github.com/menduong/oauth2/user/delivery/http"
	_userHttpDeliveryMiddleware "github.com/menduong/oauth2/user/delivery/http/middleware"
	_userRepo "github.com/menduong/oauth2/user/repository/mysql"
	_redisRepo "github.com/menduong/oauth2/user/repository/redis"
	_userUcase "github.com/menduong/oauth2/user/usecase"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
		log.Println(viper.GetString(`smtp.content`))
	}
}

func main() {

	// dbHost := viper.GetString(`database.host`)
	// dbPort := viper.GetString(`database.port`)
	// dbUser := viper.GetString(`database.user`)
	// dbPass := viper.GetString(`database.pass`)
	// dbName := viper.GetString(`database.name`)

	err := godotenv.Load(filepath.Join("./env", "test.env"))
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	// Connect mysql
	mySQL_DB_NAME := os.Getenv("MYSQL_DB_NAME")
	mySQL_USER := os.Getenv("MYSQL_USER")
	mySQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	mySQL_HOST := os.Getenv("MYSQL_HOST")
	mySQL_PORT := os.Getenv("MYSQL_PORT")

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mySQL_USER, mySQL_PASSWORD, mySQL_HOST, mySQL_PORT, mySQL_DB_NAME)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Connect redis
	// redisHost := viper.GetString(`redis.host`)
	// redisPort := viper.GetString(`redis.port`)
	// redisDB := viper.GetInt(`redis.db`)
	// redisPassword := viper.GetString(`redis.password`)
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisDBSlot, err := strconv.Atoi(os.Getenv("REDIS_DB_SLOT"))
	if err != nil {
		os.Exit(2)
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisConn := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       redisDBSlot, // use default DB
	})

	e := echo.New()
	middL := _userHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	ur := _userRepo.NewMysqlUserRepository(dbConn)

	userRedis := _redisRepo.NewRedisUserRepository(redisConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	uu := _userUcase.NewUserUsecase(ur, userRedis, timeoutContext)
	_userHttpDelivery.NewUserHandler(e, uu)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
