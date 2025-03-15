package clients

import (
	"context"
	"field-service/clients/config"
	"field-service/common/util"
	config2 "field-service/config"
	"field-service/constants"
	"fmt"
	"net/http"
	"time"


	"github.com/sirupsen/logrus"

)

type UserClient struct {
	client config.IClientConfig
}

type IUserClient interface {
	GetUserByToken(context.Context) (*UserData, error)
}

func NewUserClient(client config.IClientConfig) IUserClient {
	return &UserClient{client: client}
}

func (u *UserClient) GetUserByToken(ctx context.Context) (*UserData, error) {
	unixTime := time.Now().Unix()
	logrus.Infof("field requestAt: %d", unixTime)
	logrus.Infof("field serviceName: %s", config2.Config.AppName)
	logrus.Infof("field signatureKey: %s", u.client.SignatureKey())
	generateAPIKey := fmt.Sprintf("%s:%s:%d",
			config2.Config.AppName,
			u.client.SignatureKey(),
			unixTime,
	)
	logrus.Infof("field generateAPIKey: %s", generateAPIKey)

	apiKey := util.GenerateSHA256(generateAPIKey)
	token := ctx.Value(constants.Token).(string)
	bearerToken := fmt.Sprintf("Bearer %s", token)

	logrus.Infof("field apiKey: %s", apiKey)
	var response UserResponse
	request := u.client.Client().
		Get(fmt.Sprintf("%s/api/v1/auth/user", u.client.BaseURL())).
		Set(constants.Authorization, bearerToken).
		Set(constants.XServiceName, config2.Config.AppName).
		Set(constants.XApiKey, apiKey).
		Set(constants.XRequestAt, fmt.Sprintf("%d", unixTime))

	resp, _, errs := request.EndStruct(&response)
	if len(errs) > 0 {
			return nil, errs[0]
	}

	if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("user response: %s", response.Message)
	}

	return &response.Data, nil
}
