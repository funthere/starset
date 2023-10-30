package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/funthere/starset/productservice/domain"
	helpers "github.com/funthere/starset/productservice/helper"
	"github.com/funthere/starset/productservice/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	productHandler "github.com/funthere/starset/productservice/product/delivery/http"
	productRepo "github.com/funthere/starset/productservice/product/repository"
	productUsecase "github.com/funthere/starset/productservice/product/usecase"
	"github.com/funthere/starset/productservice/service/user"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("err")
	}
}

func main() {
	db, err := service.InitDatabase()
	if err != nil {
		log.Fatalln(err)
	}

	sqliteDB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}

	if err := sqliteDB.Ping(); err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err := sqliteDB.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Validator = domain.NewValidator()

	// user service init
	envUserSrvURL, ok := os.LookupEnv("USER_SERVICE_URL")
	if !ok {
		log.Fatalln("Error lookup ENV USER_SERVICE_URL")
	}
	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: &http.Transport{},
	}
	userSrvURL, err := url.Parse(envUserSrvURL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	userSvc := user.NewUserService(httpClient, userSrvURL.String())

	// product init
	productRepo := productRepo.NewProductRepository(db)
	productUc := productUsecase.NewProductUsecase(productRepo, userSvc)
	productHandler.NewProductHandler(e, productUc)

	// start server
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := e.Start(helpers.ServerAddress()); err != nil {
			log.Fatalln(err)
		}
	}()

	wg.Wait()

}
