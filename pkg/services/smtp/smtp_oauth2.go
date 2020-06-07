package smtp

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/smtp"
	"strings"
)

type oauth2Auth struct {
	username, accessToken string
}

// OAuth2Auth returns an Auth that implements the SASL XOAUTH2 authentication
// as per https://developers.google.com/gmail/imap/xoauth2-protocol
func OAuth2Auth(username, accessToken string) smtp.Auth {
	return &oauth2Auth{username, accessToken}
}

func (a *oauth2Auth) Start(_ *smtp.ServerInfo) (string, []byte, error) {

	resp := []byte("user=" + a.username + "\x01auth=Bearer " + a.accessToken + "\x01\x01")

	return "XOAUTH2", resp, nil
}

func (a *oauth2Auth) Next(_ []byte, _ bool) ([]byte, error) {
	return nil, nil
}

func OAuth2GeneratorFile(file string) (string, error) {
	jsonData, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	var p struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURL  string   `json:"redirect_url"`
		AuthURL      string   `json:"auth_url"`
		TokenURL     string   `json:"token_url"`
		Hostname     string   `json:"smtp_hostname"`
		Scopes       []string `json:"scopes"`
	}

	if err := json.Unmarshal(jsonData, &p); err != nil {
		return "", err
	}

	conf := oauth2.Config{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   p.AuthURL,
			TokenURL:  p.TokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
		RedirectURL: p.RedirectURL,
		Scopes:      p.Scopes,
	}

	return generateOauth2Url(&conf, p.Hostname)
}

func OAuth2Generator() (string, error) {

	var clientId string
	fmt.Print("ClientID: ")
	_, err := fmt.Scanln(&clientId)
	if err != nil {
		return "", err
	}

	var clientSecret string
	fmt.Print("ClientSecret: ")
	_, err = fmt.Scanln(&clientSecret)
	if err != nil {
		return "", err
	}

	var authUrl string
	fmt.Print("AuthURL: ")
	_, err = fmt.Scanln(&authUrl)
	if err != nil {
		return "", err
	}

	var tokenUrl string
	fmt.Print("TokenURL: ")
	_, err = fmt.Scanln(&tokenUrl)
	if err != nil {
		return "", err
	}

	var redirectUrl string
	fmt.Print("RedirectURL: ")
	_, err = fmt.Scanln(&redirectUrl)
	if err != nil {
		return "", err
	}

	var scopes string
	fmt.Print("Scopes: ")
	_, err = fmt.Scanln(&scopes)
	if err != nil {
		return "", err
	}

	var hostname string
	fmt.Print("SMTP Hostname: ")
	_, err = fmt.Scanln(&hostname)
	if err != nil {
		return "", err
	}

	conf := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authUrl,
			TokenURL:  tokenUrl,
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
		RedirectURL: redirectUrl,
		Scopes:      strings.Split(scopes, ","),
	}

	return generateOauth2Url(&conf, hostname)
}

func OAuth2GeneratorGmail(credFile string) (string, error) {
	data, err := ioutil.ReadFile(credFile)
	if err != nil {
		return "", err
	}

	conf, err := google.ConfigFromJSON(data, "https://mail.google.com/")
	if err != nil {
		return "", err
	}

	return generateOauth2Url(conf, "smtp.gmail.com")

}

func generateOauth2Url(conf *oauth2.Config, host string) (string, error) {

	fmt.Printf("Visit the following URL to authenticate:\n%s\n\n", conf.AuthCodeURL(""))

	var verCode string
	fmt.Print("Enter verification code: ")
	_, err := fmt.Scanln(&verCode)
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	token, err := conf.Exchange(ctx, verCode)
	if err != nil {
		return "", err
	}

	var sender string
	fmt.Print("Enter sender e-mail: ")
	_, err = fmt.Scanln(&sender)
	if err != nil {
		return "", err
	}

	svcConf := &Config{
		Host:        host,
		Port:        25,
		Username:    sender,
		Password:    token.AccessToken,
		FromAddress: sender,
		FromName:    "Shoutrrr",
		ToAddresses: []string{sender},
		Auth:        authTypes.OAuth2,
		UseStartTLS: true,
		UseHTML:     true,
	}

	return svcConf.GetURL().String(), nil
}
