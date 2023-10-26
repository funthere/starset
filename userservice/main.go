package main

import (
	"log"
	"sync"

	"github.com/funthere/starset/userservice/domain"
	"github.com/funthere/starset/userservice/helpers"
	"github.com/funthere/starset/userservice/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	userHandler "github.com/funthere/starset/userservice/user/delivery/http"
	userRepo "github.com/funthere/starset/userservice/user/repository"
	userUsecase "github.com/funthere/starset/userservice/user/usecase"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv CustomValidator) Validate(data interface{}) error {
	return cv.validator.Struct(data)
}

func main() {
	db, err := service.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	sqliteDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err = sqliteDB.Ping(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := sqliteDB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Validator = domain.NewValidator()

	// user init
	userRepo := userRepo.NewUserRepository(db)
	userUc := userUsecase.NewUserUsecase(userRepo)
	userHandler.NewUserHandler(e, userUc)

	// start server
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := e.Start(helpers.GetPort()); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	wg.Wait()
}
