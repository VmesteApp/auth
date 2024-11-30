package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/VmesteApp/auth-service/internal/entity"
	"github.com/VmesteApp/auth-service/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo        UserRepo
	api         VkWebApi
	tokenSecret string
	tokenTTL    time.Duration
	privateKey  string
}

// New - make user usecase.
func New(repo UserRepo, webapi VkWebApi, tokenSecret string, tokenTTL time.Duration, privateKey string) *UserUseCase {
	return &UserUseCase{
		repo:        repo,
		api:         webapi,
		tokenSecret: tokenSecret,
		tokenTTL:    tokenTTL,
		privateKey:  privateKey,
	}
}

func (u *UserUseCase) CreateAccount(ctx context.Context, email, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("can't generate password hash: %w", err)
	}

	err = u.repo.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return err
		}

		return fmt.Errorf("can't save user: %w", err)
	}

	return nil
}

func (u *UserUseCase) Login(ctx context.Context, email, password string) (*entity.User, string, error) {
	user, err := u.repo.User(ctx, email)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return nil, "", entity.ErrUserNotFound
		}

		return nil, "", fmt.Errorf("can't get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return nil, "", entity.ErrInvalidCredentials
	}

	token, err := u.doToken(user.ID, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("can't make token: %w", err)
	}

	return user, token, nil
}

func (u *UserUseCase) VkLoginByAccessToken(ctx context.Context, userAccessToken string) (*entity.User, string, error) {
	tokenInfo, err := u.api.ValidateUserAccessToken(userAccessToken)
	if err != nil {
		return nil, "", err
	}

	return u.doVkLogin(ctx, tokenInfo.UserId)
}

func (u *UserUseCase) VkLogin(ctx context.Context, launchParams string) (*entity.User, string, error) {
	parsedUrl, err := url.Parse(launchParams)
	if err != nil {
		return nil, "", entity.ErrBadVkLaunchParams
	}
	queryParams := parsedUrl.Query()

	queryMap := make(map[string]string)
	for k, v := range queryParams {
		queryMap[k] = v[0]
	}

	if !u.verifyLaunchParams(queryMap) {
		return nil, "", entity.ErrBadVkLaunchParams
	}

	vkUserIDParsed, err := strconv.Atoi(queryParams.Get("vk_user_id"))

	if err != nil {
		return nil, "", entity.ErrBadVkLaunchParams
	}
	return u.doVkLogin(ctx, vkUserIDParsed)
}

func (u *UserUseCase) doVkLogin(ctx context.Context, userID int) (*entity.User, string, error) {
	user, err := u.repo.SocialUser(ctx, "vk", strconv.Itoa(userID))
	if errors.Is(err, entity.ErrUserNotFound) {
		user, err := u.repo.SaveSocialUser(ctx, "vk", strconv.Itoa(userID))
		if err != nil {
			return nil, "", fmt.Errorf("failed save social login: %w", err)
		}

		token, err := u.doToken(user.ID, user.Role)
		if err != nil {
			return nil, "", fmt.Errorf("can't make token: %w", err)
		}

		return user, token, nil
	}
	if err != nil {
		return nil, "", fmt.Errorf("failed get user by social login: %w", err)
	}

	token, err := u.doToken(user.ID, user.Role)
	if err != nil {
		return nil, "", fmt.Errorf("failed make token: %w", err)
	}

	return user, token, nil
}

func (u *UserUseCase) doToken(userId uint64, role entity.Role) (string, error) {
	payload := map[string]any{
		"uid":  userId,
		"role": role,
	}

	token, err := jwt.NewToken(payload, u.tokenSecret, u.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("can't generate token: %w", err)
	}

	return token, nil
}

func (u *UserUseCase) verifyLaunchParams(query map[string]string) bool {
	vkSubset := make(map[string]string)

	for k, v := range query {
		if strings.HasPrefix(k, "vk_") {
			vkSubset[k] = v
		}
	}

	keys := make([]string, 0, len(vkSubset))
	for k := range vkSubset {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var encodedParams []string
	for _, k := range keys {
		encodedParams = append(encodedParams, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(vkSubset[k])))
	}
	queryString := strings.Join(encodedParams, "&")

	h := hmac.New(sha256.New, []byte(u.privateKey))
	h.Write([]byte(queryString))
	hashCode := h.Sum(nil)

	base64Hash := base64.StdEncoding.EncodeToString(hashCode)
	decodedHashCode := strings.TrimRight(base64Hash, "=")
	decodedHashCode = strings.ReplaceAll(decodedHashCode, "+", "-")
	decodedHashCode = strings.ReplaceAll(decodedHashCode, "/", "_")

	return query["sign"] == decodedHashCode
}
