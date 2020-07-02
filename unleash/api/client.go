package api

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	loginEndpoint    = "/admin/login"
	featuresEndpoint = "/admin/features"
)

type Client struct {
	userAgent  string
	httpClient *resty.Client
	auth       AuthMechanism
}

// TODO: Validate baseUrl?
func NewClient(baseUrl string, userAgent string, authMechanism AuthMechanism) (*Client, error) {
	if authMechanism == nil {
		return nil, fmt.Errorf("an authentication mechanism must be provided")
	}

	jar, _ := cookiejar.New(nil)

	c := resty.New().
		SetTimeout(10 * time.Second).
		SetHostURL(baseUrl).
		SetCookieJar(jar).
		SetRetryCount(2)

	// Do login if unauthenticated
	c.AddRetryCondition(func(r *resty.Response, err error) bool {
		// TODO: Ignore auth route
		if r.StatusCode() == http.StatusUnauthorized {
			_, err := c.R().
				SetHeaders(authMechanism.headers()).
				SetBody(authMechanism.body()).
				Post(loginEndpoint)

			return err == nil
		}

		return false
	})

	return &Client{
		userAgent:  userAgent,
		httpClient: c,
		auth:       authMechanism,
	}, nil
}

func (c *Client) ListFeatureFlags() ([]Feature, error) {
	resp, err := c.httpClient.R().SetResult(FeatureResponse{}).Get(featuresEndpoint)

	if err := determineError(resp, err); err != nil {
		return nil, err
	}

	return resp.Result().(*FeatureResponse).Features, nil
}

func (c *Client) FeatureFlagByName(featureName string) (*Feature, error) {
	features, err := c.ListFeatureFlags()
	if err != nil {
		return nil, err
	}

	for _, f := range features {
		if strings.ToLower(f.Name) == strings.ToLower(featureName) {
			return &f, nil
		}
	}

	return nil, nil
}

func (c *Client) CreateFeatureFlag(feature Feature) error {
	resp, err := c.httpClient.R().SetBody(feature).Post(featuresEndpoint)

	return determineError(resp, err)
}

func (c *Client) UpdateFeatureFlag(name string, feature Feature) error {
	resp, err := c.httpClient.R().SetBody(feature).Put(featuresEndpoint + "/" + name)

	return determineError(resp, err)
}

func (c *Client) DeleteFeatureFlag(name string) error {
	resp, err := c.httpClient.R().Delete(featuresEndpoint + "/" + name)

	return determineError(resp, err)
}

func determineError(resp *resty.Response, err error) error {
	if err != nil {
		return fmt.Errorf("technical error while communicating with unleash API: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("unleash API returned an error (http response code: %d): %s", resp.StatusCode(), resp.Body())
	}

	return nil
}
