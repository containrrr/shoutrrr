package xouath2

import (
	"encoding/json"
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/services/smtp"
	"github.com/containrrr/shoutrrr/pkg/types"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"strings"
)

// Generator is the XOAuth2 Generator implementation
type Generator struct{}

// Generate generates a service URL from a set of user questions/answers
func (g *Generator) Generate(_ types.Service, props map[string]string, args []string) (types.ServiceConfig, error) {

	if provider, found := props["provider"]; found {
		if provider == "gmail" {
			return oauth2GeneratorGmail(args[0])
		}
	}

	if len(args) > 0 {
		return oauth2GeneratorFile(args[0])
	}

	return oauth2Generator()
}

func oauth2GeneratorFile(file string) (*smtp.Config, error) {
	jsonData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
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
		return nil, err
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

	return generateOauth2Config(&conf, p.Hostname)
}

func oauth2Generator() (*smtp.Config, error) {

	var clientID string
	fmt.Print("ClientID: ")
	_, err := fmt.Scanln(&clientID)
	if err != nil {
		return nil, err
	}

	var clientSecret string
	fmt.Print("ClientSecret: ")
	_, err = fmt.Scanln(&clientSecret)
	if err != nil {
		return nil, err
	}

	var authURL string
	fmt.Print("AuthURL: ")
	_, err = fmt.Scanln(&authURL)
	if err != nil {
		return nil, err
	}

	var tokenURL string
	fmt.Print("TokenURL: ")
	_, err = fmt.Scanln(&tokenURL)
	if err != nil {
		return nil, err
	}

	var redirectURL string
	fmt.Print("RedirectURL: ")
	_, err = fmt.Scanln(&redirectURL)
	if err != nil {
		return nil, err
	}

	var scopes string
	fmt.Print("Scopes: ")
	_, err = fmt.Scanln(&scopes)
	if err != nil {
		return nil, err
	}

	var hostname string
	fmt.Print("SMTP Hostname: ")
	_, err = fmt.Scanln(&hostname)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL,
			TokenURL:  tokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
		RedirectURL: redirectURL,
		Scopes:      strings.Split(scopes, ","),
	}

	return generateOauth2Config(&conf, hostname)
}

func oauth2GeneratorGmail(credFile string) (*smtp.Config, error) {
	data, err := ioutil.ReadFile(credFile)
	if err != nil {
		return nil, err
	}

	conf, err := google.ConfigFromJSON(data, "https://mail.google.com/")
	if err != nil {
		return nil, err
	}

	return generateOauth2Config(conf, "smtp.gmail.com")

}

func generateOauth2Config(conf *oauth2.Config, host string) (*smtp.Config, error) {

	fmt.Printf("Visit the following URL to authenticate:\n%s\n\n", conf.AuthCodeURL(""))

	var verCode string
	fmt.Print("Enter verification code: ")
	_, err := fmt.Scanln(&verCode)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	token, err := conf.Exchange(ctx, verCode)
	if err != nil {
		return nil, err
	}

	var sender string
	fmt.Print("Enter sender e-mail: ")
	_, err = fmt.Scanln(&sender)
	if err != nil {
		return nil, err
	}

	svcConf := &smtp.Config{
		Host:        host,
		Port:        25,
		Username:    sender,
		Password:    token.AccessToken,
		FromAddress: sender,
		FromName:    "Shoutrrr",
		ToAddresses: []string{sender},
		Auth:        smtp.AuthTypes.OAuth2,
		UseStartTLS: true,
		UseHTML:     true,
	}

	return svcConf, nil
}
