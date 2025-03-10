package clients

import (
	"context"
	configClients "field-service/clients/config"
	"field-service/common/util"
	config "field-service/config"
	"field-service/constants"
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

type UserClient struct {
	client configClients.IClientConfig
}

type IUserClient interface {
	GetUserByToken(context.Context) (*UserData, error)
}

func NewUserClient(client configClients.IClientConfig) IUserClient {
	return &UserClient{client: client}
}

func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
	unixTime := time.Now().Unix()
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
		config.Config.AppName,
		u.client.SignatureKey(),
		unixTime,
	)
	apiKey := util.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	var response UserResponse
	request := gorequest.New().
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set(constants.Authorization, bearerToken).
		Set(constants.XServiceName, config.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
		Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL()))

	// request := u.client.Client().Clone().
	// 	Set(constants.Authorization, bearerToken).
	// 	Set(constants.XServiceName, config2.Config.AppName).
	// 	Set(constants.XApiKey, apiKey).
	// 	Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime)).
	// 	Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL()))

	resp, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user response: %s", response.Message)
	}

	return &response.Data, nil
}
