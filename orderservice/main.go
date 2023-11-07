package main

import (
	"log"
	"sync"

	"github.com/funthere/starset/orderservice/domain"
	helpers "github.com/funthere/starset/orderservice/helper"
	"github.com/funthere/starset/orderservice/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	orderHandler "github.com/funthere/starset/orderservice/order/delivery/http"
	orderRepo "github.com/funthere/starset/orderservice/order/repository"
	orderUsecase "github.com/funthere/starset/orderservice/order/usecase"
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

	// product init
	orderRepo := orderRepo.NewOrderRepository(db)
	ordertUc := orderUsecase.NewOrderUsecase(orderRepo)
	orderHandler.NewOrderHandler(e, ordertUc)

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