package main

import (
	"log"
	"sync"

	"github.com/funthere/starset/userservice/domain"
	"github.com/funthere/starset/userservice/helpers"
	"github.com/funthere/starset/userservice/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	userHandler "github.com/funthere/starset/userservice/user/delivery/http"
	userRepoSitory "github.com/funthere/starset/userservice/user/repository"
	userUsecase "github.com/funthere/starset/userservice/user/usecase"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
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
	userRepo := userRepoSitory.NewUserRepository(db)
	userUc := userUsecase.NewUserUsecase(userRepo)
	userHandler.NewUserHandler(e, userUc)

	// start server
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := e.Start(helpers.ServerAddress()); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	wg.Wait()
}
