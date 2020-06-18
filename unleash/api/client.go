package api

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	loginEndpoint    = "/admin/login"
	featuresEndpoint = "/admin/features"
)

type client struct {
	userAgent  string
	httpClient *resty.Client
	auth       AuthMechanism
}

// TODO: Validate baseUrl?
func NewClient(baseUrl string, userAgent string, authMechanism AuthMechanism) (*client, error) {
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

	return &client{
		userAgent:  userAgent,
		httpClient: c,
		auth:       authMechanism,
	}, nil
}

func (c *client) ListFeatureFlags() ([]Feature, error) {
	resp, err := c.httpClient.R().SetResult(FeatureResponse{}).Get(featuresEndpoint)

	if err := determineError(resp, err); err != nil {
		return nil, err
	}

	return resp.Result().(*FeatureResponse).Features, nil
}

func determineError(resp *resty.Response, err error) error {
	if err != nil {
		return fmt.Errorf("technical error while listing feature flags: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("unleash API returned an error (http response code: %d): %s", resp.StatusCode(), resp.Body())
	}

	return nil
}
