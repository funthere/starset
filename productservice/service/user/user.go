package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/funthere/starset/productservice/helper"
	"github.com/labstack/gommon/log"
)

type UserService interface {
	GetUserByIds(ctx context.Context, ids []uint32) (users map[uint32]UserModel, err error)
}

type user struct {
	client *http.Client
	serviceURL,
	serviceUsername,
	servicePassword string
}

func NewUserService(
	client *http.Client,
	serviceURL string,
) UserService {
	return &user{
		client:     client,
		serviceURL: serviceURL,
	}
}

func (s *user) GetUserByIds(ctx context.Context, ids []uint32) (responsePayload map[uint32]UserModel, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	request, err = http.NewRequest("GET", fmt.Sprintf("%s/users", s.serviceURL), nil)
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
