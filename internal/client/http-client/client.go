package httpclient

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"booking.com/pkg/constants"
	"booking.com/pkg/utils"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CustomHttpClient struct {
	Client  HttpClient
	Scheama string
}

func NewHttpTlsClient(caCertPath string) (*CustomHttpClient, error) {
	caPool, err := utils.CreateCertPool(caCertPath)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			RootCAs: caPool,
		},
		},
	}
	return &CustomHttpClient{
		Client:  client,
		Scheama: constants.Https,
	}, nil
}

func (c *CustomHttpClient) Get(path string, queryParams map[string]string, response any) error {
	return c.do(http.MethodGet, path, queryParams, nil, response)
}

func (c *CustomHttpClient) Post(path string, payload, response any) error {
	return c.do(http.MethodPost, path, nil, payload, response)
}

func (c *CustomHttpClient) Put(path string, payload, response any) error {
	return c.do(http.MethodPut, path, nil, payload, response)
}

func (c *CustomHttpClient) Patch(path string, payload, response any) error {
	return c.do(http.MethodPatch, path, nil, payload, response)
}

func (c *CustomHttpClient) Delete(path string, queryParams map[string]string, response any) error {
	return c.do(http.MethodDelete, path, queryParams, nil, response)
}

func (c *CustomHttpClient) do(method, path string, queryParams map[string]string, payload, response any) error {
	request, err := utils.BuildHttpRequest(method, path, queryParams, payload)
	if err != nil {
		return err
	}
	httpRsp, err := c.Client.Do(request)
	if err != nil {
		return fmt.Errorf("failed %s %s: %w", method, request.URL.String(), err)
	}

	if err := utils.ParseHttpResponse(httpRsp, response); err != nil {
		return err
	}
	return nil
}
