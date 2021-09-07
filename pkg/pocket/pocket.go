package pocket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	host         = "https://getpocket.com/v3"
	authorizeUrl = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"

	endpointAdd          = "/add"
	endpointRequestToken = "/oauth/request"
	endpointAuthorize    = "/oauth/authorize"

	xErrorHeader = "X-Error"

	defaultTimeout = 10 * time.Second
)

type (
	Client struct {
		client      *http.Client
		consumerKey string
	}

	requestTokenRequest struct {
		ConsumerKey string `json:"consumer_key"`
		RedirectURI string `json:"redirect_uri"`
	}

	authorizeRequest struct {
		ConsumerKey string `json:"consumer_key"`
		Code        string `json:"code"`
	}

	AuthorizeResponse struct {
		AccessToken string `json:"access_token"`
		Username    string `json:"username"`
	}

	addRequest struct {
		URL         string `json:"url"`
		Title       string `json:"title,omitempty"`
		Tags        string `json:"tags,omitempty"`
		AccessToken string `json:"access_token"`
		ConsumerKey string `json:"consumer_key"`
	}

	AddInput struct {
		URL         string
		Title       string
		Tags        []string
		AccessToken string
	}
)

// validate ...
func (i AddInput) validate() error {
	if i.URL == "" {
		return errors.New("required URL values is empty")
	}

	if i.AccessToken == "" {
		return errors.New("access  token is empty")
	}

	return nil
}

//generateRequest ...
func (i AddInput) generateRequest(consumerKey string) addRequest {
	return addRequest{
		URL:         i.URL,
		Tags:        strings.Join(i.Tags, ","),
		Title:       i.Title,
		AccessToken: i.AccessToken,
		ConsumerKey: consumerKey,
	}
}

// NewClient ...
func NewClient(consumerKey string) (*Client, error) {
	if consumerKey == "" {
		return nil, errors.New("consumer key is empty")
	}

	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		consumerKey: consumerKey,
	}, nil
}

// GetRequestToken ...
func (c *Client) GetRequestToken(ctx context.Context, redirectUrl string) (string, error) {
	input := &requestTokenRequest{
		ConsumerKey: c.consumerKey,
		RedirectURI: redirectUrl,
	}

	values, err := c.doHTTP(ctx, endpointRequestToken, input)
	if err != nil {
		return "", err
	}

	if values.Get("code") == "" {
		return "", errors.New("empty request token in API response")
	}

	return values.Get("code"), nil
}

// GetAuthorizationURL ...
func (c *Client) GetAuthorizationURL(requestToken, redirectUrl string) (string, error) {
	if requestToken == "" || redirectUrl == "" {
		return "", errors.New("empty params")
	}

	return fmt.Sprintf(authorizeUrl, requestToken, redirectUrl), nil
}

// authorization ...
func (c *Client) authorization(ctx context.Context, requestToken string) (*AuthorizeResponse, error) {
	if requestToken == "" {
		return nil, errors.New("empty request token")
	}

	input := &authorizeRequest{
		Code:        requestToken,
		ConsumerKey: c.consumerKey,
	}

	values, err := c.doHTTP(ctx, endpointAuthorize, input)
	if err != nil {
		return nil, err
	}

	accessToken, username := values.Get("acces_token"), values.Get("username")
	if accessToken == "" {
		return nil, errors.New("empty access token in API response")
	}

	return &AuthorizeResponse{
		AccessToken: accessToken,
		Username:    username,
	}, nil
}

// add ...
func (c *Client) add(ctx context.Context, input AddInput) error {
	if err := input.validate(); err != nil {
		return nil
	}

	req := input.generateRequest(c.consumerKey)
	_, err := c.doHTTP(ctx, endpointAdd, req)
	return err
}

// doHTTP ...
func (c *Client) doHTTP(ctx context.Context, endpoint string, body interface{}) (url.Values, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return url.Values{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, host+endpoint, bytes.NewReader(b))
	if err != nil {
		return url.Values{}, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF8")

	resp, err := c.client.Do(req)
	if err != nil {
		return url.Values{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Sprintf("API Error - %s", resp.Header.Get(xErrorHeader))
		return url.Values{}, errors.New(err)
	}

	respSecond, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return url.Values{}, err
	}

	values, err := url.ParseQuery(string(respSecond))
	if err != nil {
		return url.Values{}, err
	}
	return values, nil
}
