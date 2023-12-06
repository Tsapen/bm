package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/google/go-querystring/query"

	bm "github.com/Tsapen/bm/internal/bm"
	"github.com/Tsapen/bm/pkg/api"
)

func booksPath(id int64) string {
	if id > 0 {
		return path.Join("/api/v1/books", strconv.FormatInt(id, 10))
	}

	return "/api/v1/books"
}

func collectionsPath(id int64) string {
	if id > 0 {
		return path.Join("/api/v1/collections", strconv.FormatInt(id, 10))
	}

	return "/api/v1/collections"
}

func booksCollectionPath(id int64) string {
	return path.Join("/api/v1/collections", strconv.FormatInt(id, 10), "books")
}

func (c *Client) GetBook(ctx context.Context, req *api.GetBookReq) (*api.GetBookResp, error) {
	resp := new(api.GetBookResp)
	err := c.doRequestWithURLParams(ctx, booksPath(req.ID), nil, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetBooks(ctx context.Context, req *api.GetBooksReq) (*api.GetBooksResp, error) {
	resp := new(api.GetBooksResp)
	err := c.doRequestWithURLParams(ctx, booksPath(0), req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) CreateBook(ctx context.Context, req *api.CreateBookReq) (*api.CreateBookResp, error) {
	resp := new(api.CreateBookResp)
	err := c.doRequestWithJSON(ctx, booksPath(0), http.MethodPost, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) UpdateBook(ctx context.Context, req *api.UpdateBookReq) (bool, error) {
	err := c.doRequestWithJSON(ctx, booksPath(req.ID), http.MethodPut, req, nil)
	if err != nil {
		return false, fmt.Errorf("do request: %w", err)
	}

	return true, nil
}

func (c *Client) DeleteBooks(ctx context.Context, req *api.DeleteBooksReq) (bool, error) {
	err := c.doRequestWithJSON(ctx, booksPath(0), http.MethodDelete, req, nil)
	if err != nil {
		return false, fmt.Errorf("do request: %w", err)
	}

	return true, nil
}

func (c *Client) GetCollection(ctx context.Context, req *api.GetCollectionReq) (*api.GetCollectionResp, error) {
	resp := new(api.GetCollectionResp)
	err := c.doRequestWithURLParams(ctx, collectionsPath(req.ID), nil, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetCollections(ctx context.Context, req *api.GetCollectionsReq) (*api.GetCollectionsResp, error) {
	resp := new(api.GetCollectionsResp)
	err := c.doRequestWithURLParams(ctx, collectionsPath(0), req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) CreateCollection(ctx context.Context, req *api.CreateCollectionReq) (*api.CreateCollectionResp, error) {
	resp := new(api.CreateCollectionResp)
	err := c.doRequestWithJSON(ctx, collectionsPath(0), http.MethodPost, req, resp)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

func (c *Client) UpdateCollection(ctx context.Context, req *api.UpdateCollectionReq) (bool, error) {
	err := c.doRequestWithJSON(ctx, collectionsPath(req.ID), http.MethodPut, req, nil)
	if err != nil {
		return false, fmt.Errorf("do request: %w", err)
	}

	return true, nil
}

func (c *Client) DeleteCollection(ctx context.Context, req *api.DeleteCollectionReq) (bool, error) {
	err := c.doRequestWithJSON(ctx, collectionsPath(req.ID), http.MethodDelete, nil, nil)
	if err != nil {
		return false, fmt.Errorf("do request: %w", err)
	}

	return true, nil
}

func (c *Client) CreateBooksCollection(ctx context.Context, req *api.CreateBooksCollectionReq) (bool, error) {
	err := c.doRequestWithJSON(ctx, booksCollectionPath(req.CID), http.MethodPost, req, nil)
	if err != nil {
		return false, fmt.Errorf("do request: %w", err)
	}

	return true, nil
}

func (c *Client) DeleteBooksCollection(ctx context.Context, req *api.DeleteBooksCollectionReq) (bool, error) {
	err := c.doRequestWithJSON(ctx, booksCollectionPath(req.CID), http.MethodDelete, req, nil)
	if err != nil {
		return false, fmt.Errorf("do request: %w", err)
	}

	return true, nil
}

func (c *Client) doRequestWithJSON(ctx context.Context, urlPath, method string, reqData, respData any) (err error) {
	body, err := json.Marshal(reqData)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	u, err := url.Parse(c.cfg.Address)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	u.Path = path.Join(u.Path, urlPath)

	req, err := http.NewRequest(method, u.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("construct request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer func() {
		err = bm.HandleErrPair(resp.Body.Close(), err)
	}()

	if resp.StatusCode >= http.StatusInternalServerError {
		return fmt.Errorf("server error: %d", resp.StatusCode)
	}

	if http.StatusOK == resp.StatusCode {
		if respData == nil {
			return nil
		}

		if err = json.NewDecoder(resp.Body).Decode(respData); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}

		return nil
	}

	errResp := map[string]any{}
	if err = json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return fmt.Errorf("decode error response: %w", err)
	}

	return fmt.Errorf("%+v", errResp)
}

func (c *Client) doRequestWithURLParams(ctx context.Context, urlPath string, reqData, respData any) (err error) {
	u, err := url.Parse(c.cfg.Address)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	u.Path = path.Join(u.Path, urlPath)

	if reqData != nil {
		vals, err := query.Values(reqData)
		if err != nil {
			return fmt.Errorf("construct request: %w", err)
		}

		u.RawQuery = vals.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("construct request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer func() {
		err = bm.HandleErrPair(resp.Body.Close(), err)
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get error http status: %d", resp.StatusCode)
	}

	if respData != nil || http.StatusOK < resp.StatusCode && resp.StatusCode < http.StatusInternalServerError {
		if err = json.NewDecoder(resp.Body).Decode(respData); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}
