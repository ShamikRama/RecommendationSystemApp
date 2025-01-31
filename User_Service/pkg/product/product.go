package product

import (
	"User_Service/internal/logger"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"go.uber.org/zap"
)

type clientlocal struct {
	logger logger.Logger
	cl     *ClientProduct
}

func NewClient(logger logger.Logger) Client {
	return &clientlocal{
		logger: logger,
		cl:     NewClientProducts(),
	}
}

func (c *clientlocal) GetProducts(ctx context.Context, pageSize int, page int) (ProductsResponse, error) {
	path, err := url.JoinPath(c.cl.BaseUrl, "products")
	if err != nil {
		c.logger.Error("Error joining path", zap.Error(err))
		return ProductsResponse{}, err
	}

	params := url.Values{}
	if pageSize > 0 && page > 0 {
		params.Add("page_size", strconv.Itoa(pageSize))
		params.Add("page", strconv.Itoa(page))
	}

	fullURL := path
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		c.logger.Error("Error creating request", zap.Error(err))
		return ProductsResponse{}, err
	}

	req.Header.Add("X-Auth-Token", "lwehvowhvowvhwovwfwefwefwefw")

	c.logger.Info("Executing request",
		zap.String("url", fullURL),
		zap.String("method", http.MethodGet))

	resp, err := c.cl.Client.Do(req)
	if err != nil {
		c.logger.Error("Error executing request", zap.Error(err))
		return ProductsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Unexpected status code", zap.Int("status", resp.StatusCode))
		return ProductsResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		c.logger.Error("Response body is nil")
		return ProductsResponse{}, errors.New("response body is nil")
	}

	var response ProductsResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		c.logger.Error("Error decoding response", zap.Error(err))
		return ProductsResponse{}, err
	}

	c.logger.Info("Response received",
		zap.Int("status", resp.StatusCode),
		zap.Int("products_count", len(response.Data)),
		zap.Int("total", response.Total))

	return response, nil
}

func (c *clientlocal) CreateCart(ctx context.Context, request CartCreateRequest) error {
	path, err := url.JoinPath(c.cl.BaseUrl, ShortClient, "/cart")
	if err != nil {
		c.logger.Info("Error joining path", zap.Error(err))
		return err
	}

	body, err := json.Marshal(request)
	if err != nil {
		c.logger.Info("Error marshaling request body", zap.Error(err))
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(body))
	if err != nil {
		c.logger.Info("Error creating request", zap.Error(err))
		return err
	}

	req.Header.Add("X-Auth-Token", c.cl.ApiKey)

	c.logger.Info("Executing request", zap.String("url", path), zap.String("method", http.MethodPost))

	resp, err := c.cl.Client.Do(req)
	if err != nil {
		c.logger.Info("Error executing request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Info("Unexpected status code", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		c.logger.Info("Response body is nil")
		return errors.New("response body is nil")
	}

	return nil
}

func (c *clientlocal) DeleteCart(ctx context.Context, request CartDeleteRequest) error {
	path, err := url.JoinPath(c.cl.BaseUrl, ShortClient, "/cart")
	if err != nil {
		c.logger.Info("Error joining path", zap.Error(err))
		return err
	}

	body, err := json.Marshal(request)
	if err != nil {
		c.logger.Info("Error creating request", zap.Error(err))
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, bytes.NewBuffer(body))
	if err != nil {
		c.logger.Info("Error creating request", zap.Error(err))
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", c.cl.ApiKey)

	resp, err := c.cl.Client.Do(req)
	if err != nil {
		c.logger.Info("Error executing request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Info("Unexpected status code", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		c.logger.Info("Response body is nil")
		return errors.New("response body is nil")
	}

	return nil
}

func (c *clientlocal) UpdateCart(ctx context.Context, request CartUpdateRequest) error {
	path, err := url.JoinPath(c.cl.BaseUrl, ShortClient, "/cart")
	if err != nil {
		c.logger.Info("Error joining path", zap.Error(err))
		return err
	}

	body, err := json.Marshal(request)
	if err != nil {
		c.logger.Info("Error creating request", zap.Error(err))
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewBuffer(body))
	if err != nil {
		c.logger.Info("Error creating request", zap.Error(err))
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", c.cl.ApiKey)

	resp, err := c.cl.Client.Do(req)
	if err != nil {
		c.logger.Info("Error executing request", zap.Error(err))
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Info("Unexpected status code", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.Body == nil {
		c.logger.Info("Response body is nil")
		return errors.New("response body is nil")
	}

	return nil
}
