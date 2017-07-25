package spotify

import (
	"encoding/json"
	"fmt"
	client "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
)

const (
	Token        = "token"
	ClientId     = "clientid"
	ClientSecret = "clientsecret"

	CallbackUrl = "http://localhost:11075/callback"
)

var Scopes = []string{
	client.ScopePlaylistReadPrivate,
	client.ScopeUserLibraryRead,
	client.ScopeUserReadCurrentlyPlaying,
	client.ScopeUserReadPlaybackState,
	client.ScopeUserModifyPlaybackState,
	client.ScopeUserReadRecentlyPlayed,
}

func SetToken(token *oauth2.Token) error {
	tokenData, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return Set(Token, tokenData)
}

func GetClientCreds() (string, string, error) {
	clientId, err := readData(ClientId)
	if err != nil {
		return "", "", err
	}

	clientSecret, err := readData(ClientSecret)
	if err != nil {
		return "", "", err
	}

	return string(clientId), string(clientSecret), nil
}

func GetToken() (*oauth2.Token, error) {
	tokenData, err := readData(Token)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	err = json.Unmarshal(tokenData, token)
	if err != nil {
		return nil, err
	}

	config, err := AuthConfig()
	if err != nil {
		return nil, err
	}
	tokenSource := config.TokenSource(oauth2.NoContext, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if newToken.AccessToken != token.AccessToken {
		err = SetToken(newToken)
		if err != nil {
			return nil, err
		}

		token = newToken
	}

	return token, nil
}

func NewAuth() (*client.Authenticator, error) {
	auth := client.NewAuthenticator(CallbackUrl, Scopes...)
	clientId, clientSecret, err := GetClientCreds()
	if err != nil {
		return nil, err
	}

	auth.SetAuthInfo(clientId, clientSecret)
	return &auth, nil
}

func AuthConfig() (*oauth2.Config, error) {
	clientId, clientSecret, err := GetClientCreds()
	if err != nil {
		return nil, err
	}

	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  CallbackUrl,
		Scopes:       Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  client.AuthURL,
			TokenURL: client.TokenURL,
		},
	}, nil
}

func Set(dataId string, data []byte) error {
	return ioutil.WriteFile(dataPath(dataId), data, 0644)
}

func dataPath(dataId string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("alfred_workflow_data"), dataId)
}

func readData(dataId string) ([]byte, error) {
	return ioutil.ReadFile(dataPath(dataId))
}

func IsDataSet(dataId string) bool {
	// terrible, other errors could occur
	_, err := os.Stat(dataPath(dataId))
	if err == nil {
		return true
	}
	return false
}
