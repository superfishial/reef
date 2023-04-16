package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/superfishial/reef/server/config"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	Email             string   `json:"email"`
	EmailVerified     bool     `json:"email_verified"`
	Name              string   `json:"name"`
	GivenName         string   `json:"given_name"`
	FamilyName        string   `json:"family_name"`
	PreferredUsername string   `json:"preferred_username"`
	Nickname          string   `json:"nickname"`
	Groups            []string `json:"groups"`
	Sub               string   `json:"sub"`
}

func getOAuthClient(conf config.Config) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  conf.Server.RootURL + "/v1/auth/callback",
		ClientID:     conf.OAuth.ClientID,
		ClientSecret: conf.OAuth.ClientSecret,
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.OAuth.AuthEndpoint,
			TokenURL: conf.OAuth.TokenEndpoint,
		},
	}
}

// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func GenerateStateToken() (string, error) {
	token, err := GenerateRandomString(64)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString([]byte(token)), nil
}

func GetUserInfo(accessToken string) (UserInfo, error) {
	resp, err := http.Get("https://auth.super.fish/application/o/userinfo/?access_token=" + accessToken)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed while fetching user info: %w", err)
	}
	if resp.StatusCode != 200 {
		return UserInfo{}, fmt.Errorf("failed while fetching user info with status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed while reading body of user info response: %w", err)
	}

	userInfo, err := ParseUserInfo(contents)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed while parsing body of user info response as UserInfo struct: %w", err)
	}
	return userInfo, nil
}

func ParseUserInfo(data []byte) (UserInfo, error) {
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return UserInfo{}, fmt.Errorf("failed while unmarshalling user info as JSON: %w", err)
	}
	return userInfo, nil
}

func SignToken(jwtSecret string, userInfo UserInfo) (string, error) {
	groupsString, err := json.Marshal(userInfo.Groups)
	if err != nil {
		return "", fmt.Errorf("failed to stringify groups into JSON: %w", err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     userInfo.Email,
		"name":      userInfo.Name,
		"nickname":  userInfo.Nickname,
		"sub":       userInfo.Sub,
		"givenName": userInfo.GivenName,
		"groups":    groupsString,
	})
	tokenString, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return tokenString, nil
}
