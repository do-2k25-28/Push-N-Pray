package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	auth       Authenticator
}

type Option func(*Client)

func NewClient(baseURL string, opts ...Option) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("api: base URL is required")
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("api: parse base URL: %w", err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("api: base URL must include scheme and host")
	}

	if !strings.HasSuffix(parsed.Path, "/") {
		parsed.Path += "/"
	}

	client := &Client{
		baseURL:    parsed,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(client)
	}

	if client.httpClient == nil {
		return nil, fmt.Errorf("api: http client is required")
	}

	return client, nil
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithBearerToken(token string) Option {
	return func(c *Client) {
		c.auth = BearerToken(token)
	}
}

func WithBasicAuth(email, token string) Option {
	return func(c *Client) {
		c.auth = BasicAuth{Email: email, Token: token}
	}
}

func WithAuth(auth Authenticator) Option {
	return func(c *Client) {
		c.auth = auth
	}
}

func (c *Client) Health(ctx context.Context) error {
	req, err := c.newRequest(ctx, http.MethodGet, "health", nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

func (c *Client) Register(ctx context.Context, payload RegisterRequest) (*AuthResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "auth/register", payload)
	if err != nil {
		return nil, err
	}

	var resp AuthResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Login(ctx context.Context, payload LoginRequest) (*AuthResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "auth/login", payload)
	if err != nil {
		return nil, err
	}

	var resp AuthResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Token(ctx context.Context, payload TokenRequest) (*AuthResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "auth/token", payload)
	if err != nil {
		return nil, err
	}

	var resp AuthResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) ListPATs(ctx context.Context) (*ListPATsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "tokens", nil)
	if err != nil {
		return nil, err
	}

	var resp ListPATsResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreatePAT(ctx context.Context, payload CreatePATRequest) (*CreatePATResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "tokens", payload)
	if err != nil {
		return nil, err
	}

	var resp CreatePATResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) DeletePAT(ctx context.Context, tokenID string) error {
	path := fmt.Sprintf("tokens/%s", url.PathEscape(tokenID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

func (c *Client) CreateProject(ctx context.Context, payload CreateProjectRequest) (*CreateProjectResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "projects", payload)
	if err != nil {
		return nil, err
	}

	var resp CreateProjectResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) DeleteProject(ctx context.Context, projectID string) error {
	path := fmt.Sprintf("projects/%s", url.PathEscape(projectID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

func (c *Client) DeployProject(ctx context.Context, projectID string, payload DeployProjectRequest) (*DeployProjectResponse, error) {
	path := fmt.Sprintf("projects/%s/deploy", url.PathEscape(projectID))
	req, err := c.newRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}

	var resp DeployProjectResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) DeploymentInfo(ctx context.Context, projectID, deploymentID string) (*DeploymentInfoResponse, error) {
	path := fmt.Sprintf("projects/%s/deployments/%s", url.PathEscape(projectID), url.PathEscape(deploymentID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var resp DeploymentInfoResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) SetProjectEnv(ctx context.Context, projectID string, payload SetProjectEnvRequest) error {
	path := fmt.Sprintf("projects/%s/env", url.PathEscape(projectID))
	req, err := c.newRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	var buf io.Reader

	if body != nil {
		payload := &bytes.Buffer{}
		enc := json.NewEncoder(payload)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, fmt.Errorf("api: encode request body: %w", err)
		}
		buf = payload
	}

	target := c.baseURL.ResolveReference(&url.URL{Path: strings.TrimPrefix(path, "/")})
	req, err := http.NewRequestWithContext(ctx, method, target.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("api: create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.auth != nil {
		c.auth.Apply(req)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("api: request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return readHTTPError(resp)
	}

	if out == nil || resp.StatusCode == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("api: decode response: %w", err)
	}

	return nil
}
