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

type validateUserAccessTokenResponse struct {
	Response entity.VkTokenInfo `json:"response"`
	Error    *struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"error,omitempty"`
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

	var data validateUserAccessTokenResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("can't unmarshal body: %w", err)
	}

	if data.Error != nil {
		switch data.Error.ErrorMsg {
		case "Access denied: Incorrect token invalid_token":
			return nil, entity.ErrBadVkToken
		case "Access denied: Incorrect token session_expired":
			return nil, entity.ErrVkTokenExpired
		default:
			return nil, fmt.Errorf("unknown error: %s", data.Error.ErrorMsg)
		}
	}

	return &data.Response, nil
}
