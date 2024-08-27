package webapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/VmesteApp/auth-service/internal/entity"
)

type VkWebApi struct {
	AppId         int
	ServiceSecret string
}

func New(appId int, serviceSecret string) *VkWebApi {
	return &VkWebApi{
		AppId:         appId,
		ServiceSecret: serviceSecret,
	}
}

func (vk *VkWebApi) ValidateUserAccessToken(userAccessToken string) (*entity.VkTokenInfo, error) {
	u, _ := url.Parse("https://api.vk.com/method/secure.checkToken")
	q := u.Query()
	q.Add("v", "5.101")
	q.Add("access_token", vk.ServiceSecret)
	q.Add("token", userAccessToken)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("can't request to validate token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read body: %w", err)
	}

	var data map[string]entity.VkTokenInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("can't unmarshal body: %w", err)
	}

	tokenInfoWrapper := data["response"]

	return &tokenInfoWrapper, nil
}
