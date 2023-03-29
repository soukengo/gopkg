package http

import (
	"context"
	"github.com/go-resty/resty/v2"
	"io"
	"strings"
)

const (
	JSON = "application/json"
)

type Client struct {
	cli  *resty.Client
	opts *options
}

func NewClient(opts ...Option) *Client {
	opt := new(options)
	opt.apply(opts...)
	return &Client{opts: opt, cli: newRestyClient(opt)}
}

func newRestyClient(opts *options) *resty.Client {
	cli := resty.New().SetTimeout(opts.timeout)
	if len(opts.baseUrl) > 0 {
		cli.SetBaseURL(opts.baseUrl)
	}
	if opts.basicAuth != nil {
		cli.SetBasicAuth(opts.basicAuth.username, opts.basicAuth.password)
	}
	if len(opts.proxyURL) > 0 {
		cli.SetProxy(opts.proxyURL)
	}
	return cli
}

// Post send POST request
func (c *Client) Post(ctx context.Context, addr string, body io.Reader, headers map[string]string) (out *Response, err error) {
	resp, err := c.newRequest(ctx).SetHeaders(headers).SetBody(body).Post(addr)
	if err != nil {
		return
	}
	out = newResponse(resp.Body(), resp.RawResponse)
	return
}

// PostForm send FORM POST request
func (c *Client) PostForm(ctx context.Context, addr string, params map[string]string, headers map[string]string) (out *Response, err error) {
	resp, err := c.newRequest(ctx).SetHeaders(headers).SetFormData(params).Post(addr)
	if err != nil {
		return
	}
	out = newResponse(resp.Body(), resp.RawResponse)
	return
}

// PostJSON send JSON POST request
func (c *Client) PostJSON(ctx context.Context, addr string, json string, headers map[string]string) (out *Response, err error) {
	if len(headers) == 0 {
		headers = make(map[string]string)
	}
	headers["content-type"] = JSON
	resp, err := c.newRequest(ctx).SetHeaders(headers).SetBody(strings.NewReader(json)).Post(addr)
	if err != nil {
		return
	}
	out = newResponse(resp.Body(), resp.RawResponse)
	return
}

func (c *Client) newRequest(ctx context.Context) *resty.Request {
	return c.cli.NewRequest().SetContext(ctx)
}
