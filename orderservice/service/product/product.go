package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/funthere/starset/orderservice/helper"
	"github.com/labstack/gommon/log"
)

type ProductService interface {
	GetProductByIds(ctx context.Context, ids []int64) (products map[int64]Product, err error)
}

type product struct {
	client *http.Client
	serviceURL,
	serviceUsername,
	servicePassword string
}

func NewProductService(
	client *http.Client,
	serviceURL string,
) ProductService {
	return &product{
		client:     client,
		serviceURL: serviceURL,
	}
}

func (s *product) GetProductByIds(ctx context.Context, ids []int64) (responsePayload map[int64]Product, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	request, err = http.NewRequest("GET", fmt.Sprintf("%s/product/fetch", s.serviceURL), nil)
	if err != nil {
		log.Info(err)
		return
	}

	queryParams := request.URL.Query()
	queryParams.Add("ids", helper.IntSliceToString(ids))
	request.URL.RawQuery = queryParams.Encode()

	if response, err = s.client.Do(request); err != nil {
		log.Info(err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errors.New("bad request")
		log.Info(err)
		return
	}

	err = json.NewDecoder(response.Body).Decode(&responsePayload)
	if err != nil {
		err = errors.New("err when Decode")
		log.Info(err)
		return
	}
	return
}
