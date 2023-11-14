package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/funthere/starset/orderservice/domain"
	helpers "github.com/funthere/starset/orderservice/helper"
	dbService "github.com/funthere/starset/orderservice/service/db"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	orderHandler "github.com/funthere/starset/orderservice/order/delivery/http"
	orderRepository "github.com/funthere/starset/orderservice/order/repository"
	orderUsecase "github.com/funthere/starset/orderservice/order/usecase"
	"github.com/funthere/starset/orderservice/service/product"
	"github.com/funthere/starset/orderservice/service/user"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	db, err := dbService.InitPostgreDatabase()
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
	e.Use(middleware.Logger())
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

	// product service init
	envProductSrvURL, ok := os.LookupEnv("PRODUCT_SERVICE_URL")
	if !ok {
		log.Fatalln("Error lookup ENV PRODUCT_SERVICE_URL")
	}
	productSrvURL, err := url.Parse(envProductSrvURL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	productSvc := product.NewProductService(httpClient, productSrvURL.String())

	// order init
	orderRepo := orderRepository.NewOrderRepository(db)
	orderUc := orderUsecase.NewOrderUsecase(orderRepo, productSvc, userSvc)
	orderHandler.NewOrderHandler(e, orderUc)

	// Start server
	e.Logger.Fatal(e.Start(helpers.ServerAddress()))

}
